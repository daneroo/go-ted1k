package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/daneroo/go-ted1k/logsetup"
	"github.com/daneroo/go-ted1k/postgres"
	"github.com/daneroo/go-ted1k/types"
	"github.com/jackc/pgx/v4/pgxpool"
	nats "github.com/nats-io/nats.go"
)

const (
	// myCredentials = "ted:secret@tcp(0.0.0.0:3306)/ted"
	pgCredentialsDefault = "postgres://postgres:secret@0.0.0.0:5432/ted"
	pgTablename          = "watt"
	natsURL              = "nats://nats.ts.imetrical.com:4222"
	natsConnectionName   = "subscribe.ted1k"
	topic                = "im.qcic.heartbeat"
	host                 = natsConnectionName
)

var (
	dbpool *pgxpool.Pool
	// count  = 0
	entryQueue []types.Entry
	// entryQueueCond is a condition variable used to synchronize access to the entryQueue variable.
	// It allows the processEntryQueue function to efficiently wait for new entries to be added to the queue
	// This variable is used in the enqueueMessage() and processEntryQueue() functions to add and remove entries
	// from the queue and to signal the processEntryQueue() function when new entries are available.
	entryQueueCond *sync.Cond = sync.NewCond(&sync.Mutex{})
)

// Subscribe to a nats topic, extract an entry and insert into the database
// tl;dr:
// When an entry is received, it is added to the entryQueue and the processEntryQueue() function is signalled
// to process the entry.
// If while processing the entryQueue:
// - The queue is empty, the processEntryQueue() function will wait for a signal
// - If an error occurs, during insertEntry, the entry is put back in the queue and the processing is paused (with capped exponential backoff)
// When an interrupt signal is received, the program will exit after the entryQueue is empty
// This is done to ensure that all entries are processed before the program exits
// However, if the entryQueue is not empty after 20 seconds, the program will exit anyway
// This can happen if the database is not available, at the same time as the program is being shutdown
func main() {
	logsetup.SetupFormat()
	log.Printf("Starting TED1K subscribe\n") // TODO(daneroo): add version,buildDate

	pgCredentials := os.Getenv("PGCONN")
	if pgCredentials == "" {
		pgCredentials = pgCredentialsDefault
	}

	dbpoolTmp, err := postgres.SetupPool(context.Background(), []string{pgTablename}, pgCredentials)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	dbpool = dbpoolTmp // assign to global dbpool variable
	defer dbpool.Close()

	// Connect to a nats server
	nc, err := nats.Connect(natsURL,
		// RetryOnFailedConnect is on master, ut not released
		// nats.RetryOnFailedConnect(true)
		nats.MaxReconnects(-1), // 60 is the default
		nats.Name(natsConnectionName),
	)
	if err != nil {
		log.Fatalf("Unable to connect to Nats: %s\n", natsURL)
	}

	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	// Encoded subscription - we don't need to keep a handle on the subscription since we only close the entire connection
	_, err = c.Subscribe(topic, enqueueMessage)
	if err != nil {
		log.Fatalf("Unable to subscribe to topic: %s\n", topic)
	}

	go processEntryQueue(context.Background()) // Start goroutine to process queued messages

	// Wait for an interrupt signal to gracefully shutdown the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt // Block until an interrupt signal is received

	log.Println("Shutdown signal received, exiting...")

	// Nats clean up
	// Drain connection (Preferred for responders)
	// This should all be confirmed. I think this is the correct way to do it
	// "Drain will put a connection into a drain state. All subscriptions will immediately be put into a drain state.
	//  Upon completion, the publishers will be drained and can not publish any additional messages. Upon draining of the publishers,
	//  the connection will be closed. Use the ClosedCB() option to know when the connection has moved from draining to closed
	// My Assumptions:
	//  sub.Unsubscribe() not needed if nc.Drain() called
	//  sub.Drain() not needed if this is nc.Drain() called
	//  nc.Close() not needed if this is nc.Drain() called
	nc.Drain()
	waitUntilEntryQueueEmpty()

}

type message struct {
	Stamp time.Time `json:"stamp"`
	Host  string    `json:"host"`
	Text  string    `json:"text"` // or "volt,omitempty"
}

func waitUntilEntryQueueEmpty() {
	now := time.Now()
	maxWait := 20 * time.Second
	for {
		// Wait for the queue to become empty
		entryQueueCond.L.Lock()
		if len(entryQueue) == 0 {
			log.Println("EntryQueue is empty we can safely exit")

			entryQueueCond.L.Unlock()
			return
		}
		log.Printf("EntryQueue has %d entries waiting (max=%v)", len(entryQueue), maxWait)
		entryQueueCond.L.Unlock()

		time.Sleep(1 * time.Second)
		if time.Since(now) > maxWait {
			log.Printf("Queue is not empty after %s, exiting anyway", maxWait)
			return
		}
	}
}

func enqueueMessage(m *message) {
	if m.Host == "capture.ted1k" {
		entry := types.Entry{Stamp: m.Stamp, Watt: 0}
		_, err := fmt.Sscanf(m.Text, "watts: %d", &entry.Watt)
		if err != nil {
			log.Println(err)
		}
		entryQueueCond.L.Lock()
		entryQueue = append(entryQueue, entry)
		entryQueueCond.Signal()
		entryQueueCond.L.Unlock()
	}
}

// processEntryQueue is a goroutine that processes entries from the entryQueue
// insertEntry inserts the given entry into the database
// If an error occurs, the entry is re-enqueued for retry (prepended instead of appended to preserve the intended order)
// but the processing is paused (with capped exponential backoff)
func processEntryQueue(ctx context.Context) {
	minBackoffTime := 1 * time.Second
	maxBackoffTime := 10 * time.Second
	backoffTime := minBackoffTime

	for {
		entryQueueCond.L.Lock()
		for len(entryQueue) == 0 {
			// Wait for a new entry to be added to the queue
			entryQueueCond.Wait()
		}
		entry := entryQueue[0]
		entryQueue = entryQueue[1:]
		entryQueueCond.L.Unlock()

		// Attempt to insert the entry into the database
		err := insertEntry(entry)
		if err != nil {
			log.Printf("processEntryQueue Unable to insert entry: %v\n", entry)
			// Re-enqueue the entry for retry (by prepending instead of appending to preserve the intended order)
			entryQueueCond.L.Lock()
			entryQueue = append([]types.Entry{entry}, entryQueue...)
			entryQueueCond.Signal()
			entryQueueCond.L.Unlock()
			// Wait for an exponentially (capped) increasing amount of time before retrying
			log.Printf("EntryQueue pausing (for %v) with %d entries in queue\n", backoffTime, len(entryQueue))
			time.Sleep(backoffTime)
			// double the backoff time, but cap it at maxBackoffTime
			backoffTime = minDuration(2*backoffTime, maxBackoffTime)
		} else {
			backoffTime = minBackoffTime
		}
	}
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func insertEntry(entry types.Entry) error {
	// log.Printf("Writing an entry: %v\n", entry)
	sql := fmt.Sprintf("INSERT INTO %s (stamp, watt) VALUES ($1,$2)  ON CONFLICT (stamp) DO NOTHING", pgTablename)
	// Stamp should now be Truncated at source, but we ere Rounding, now we correctly Truncate
	entry.Stamp = entry.Stamp.Truncate(time.Second)
	vals := []interface{}{entry.Stamp, entry.Watt}

	ctx := context.Background()
	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		// log.Printf("Unable to acquire database connection: %v", err)
		return err
	}
	defer conn.Release()

	commandTag, err := conn.Exec(context.Background(), sql, vals...)
	if err != nil {
		// log.Printf("Unable to insert entry: %v - %v\n", entry, err)
		return err
	}
	log.Printf("Wrote an entry: %v, affected rows:%d\n", entry, commandTag.RowsAffected())
	return nil
}
