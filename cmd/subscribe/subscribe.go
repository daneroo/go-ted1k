package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
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

	// Encoded subscription
	sub, err := c.Subscribe(topic, enqueueMessage)
	if err != nil {
		log.Fatalf("Unable to subscribe to topic: %s\n", topic)
	}

	go processEntryQueue(context.Background()) // Start goroutine to process queued messages

	// time.Sleep(10 * time.Hour)
	// Sleep forever
	<-make(chan int)

	// Unsubscribe
	sub.Unsubscribe()

	// Drain
	sub.Drain()

	// Drain connection (Preferred for responders)
	// Close() not needed if this is called.
	nc.Drain()

	// Close connection
	nc.Close()
}

type message struct {
	Stamp time.Time `json:"stamp"`
	Host  string    `json:"host"`
	Text  string    `json:"text"` // or "volt,omitempty"
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

func processEntryQueue(ctx context.Context) {
	backoffTime := 1 * time.Second
	maxBackoffTime := 10 * time.Second

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
			backoffTime = minDuration(2*backoffTime, maxBackoffTime)
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
	entry.Stamp = entry.Stamp.Round(time.Second)
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
