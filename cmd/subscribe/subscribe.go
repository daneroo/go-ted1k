package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/daneroo/go-ted1k/logsetup"
	"github.com/daneroo/go-ted1k/postgres"
	"github.com/daneroo/go-ted1k/types"
	"github.com/jackc/pgx/v4"
	nats "github.com/nats-io/nats.go"
)

const (
	// myCredentials = "ted:secret@tcp(0.0.0.0:3306)/ted"
	pgCredentials      = "postgres://postgres:secret@0.0.0.0:5432/ted"
	pgTablename        = "watt"
	natsURL            = "nats://nats.dl.imetrical.com:4222"
	natsConnectionName = "subscribe.ted1k"
	topic              = "im.qcic.heartbeat"
	host               = natsConnectionName
)

// should NOT be global!!
var conn *pgx.Conn

func main() {
	logsetup.SetupFormat()
	log.Printf("Starting TED1K subscribe\n") // TODO(daneroo): add version,buildDate

	conn = postgres.Setup(context.Background(), []string{pgTablename}, pgCredentials)
	defer conn.Close(context.Background())

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

	// Regular subscription (not encoded)
	// sub, err := nc.Subscribe(topic, func(m *nats.Msg) {
	// 	fmt.Printf("Received a message: %s\n", string(m.Data))
	// })
	// if err != nil {
	// 	log.Fatalf("Unable to subscribe to topic: %s\n", topic)
	// }

	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	// Encoded subscription
	sub, err := c.Subscribe(topic, receiveMessage)

	if err != nil {
		log.Fatalf("Unable to subscribe to topic: %s\n", topic)
	}

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

func receiveMessage(m *message) {
	// fmt.Printf("Received a message! %+v\n", m)
	if m.Host == "capture.ted1k" {
		entry := types.Entry{Stamp: m.Stamp, Watt: 0}
		_, err := fmt.Sscanf(m.Text, "watts: %d", &entry.Watt)
		if err != nil {
			log.Println(err)
		}
		// no longer called in a go routine
		insertEntry(conn, entry)
	}
}

var count = 0

// If I want to invoke from go routine, should use a connection pool
func insertEntry(conn *pgx.Conn, entry types.Entry) {
	count++
	// log.Printf("Writing an entry: %v\n", entry)
	sql := fmt.Sprintf("INSERT INTO %s (stamp, watt) VALUES ($1,$2)  ON CONFLICT (stamp) DO NOTHING", pgTablename)
	entry.Stamp = entry.Stamp.Round(time.Second)
	vals := []interface{}{entry.Stamp, entry.Watt}

	commandTag, err := conn.Exec(context.Background(), sql, vals...)
	if err != nil {
		log.Printf("Unable to insert %d entry: %v %v\n", count, entry, err)
		return
	}
	log.Printf("Wrote an entry: %d %v, affected rows:%d\n", count, entry, commandTag.RowsAffected())

}
