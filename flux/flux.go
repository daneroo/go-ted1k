package flux

import (
	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
	"github.com/influxdb/influxdb/client"
	"time"
)

import (
	"fmt"
	"log"
	"net/url"
	// "os"
)

const (
	MyHost = "docker"
	MyPort = 8086
	MyDB   = "ted"
	// MyMeasurement = "watt"
	writeBatchSize = 3600 * 24
)

func IgnoreAll(src <-chan Entry) {
	start := time.Now()
	count := 0
	for entry := range src {
		count++
		if (count % writeBatchSize) == 0 {
			log.Printf("Ignore::checkpoint at %d records %v", count, entry.Stamp)
		}
	}
	TimeTrack(start, "flux.IgnoreAll", count)
}

// Consume the Entry (receive only) channel
// preforming batched writes (of size writeBatchSize)
// Also performs progress logging (and timing)
func WriteAll(src <-chan Entry) {
	start := time.Now()
	startBatch := time.Now() // timer for internal (chunk) loop iterations
	count := 0

	con, err := connect()
	Checkerr(err)
	// defer close?

	var entries = make([]Entry, 0, writeBatchSize)
	for entry := range src {
		entries = append(entries, entry)
		count++
		if len(entries) == cap(entries) {
			entries = flush(con, entries)
			TimeTrack(startBatch, "flux.WriteAll.checkpoint+", cap(entries))
			startBatch = time.Now()
		}
	}
	_ = flush(con, entries)
	TimeTrack(start, "flux.WriteAll", count)
}

// Write out the entries to con, and reallocate a new empty slice
func flush(con *client.Client, entries []Entry) []Entry {
	writeEntries(con, entries)
	return make([]Entry, 0, writeBatchSize)
}

// Perform the batch write
func writeEntries(con *client.Client, entries []Entry) {
	var (
		pts = make([]client.Point, len(entries))
	)

	for i, entry := range entries {
		pts[i] = client.Point{
			Measurement: "watt",
			Fields: map[string]interface{}{
				"value": entry.Watt,
			},
			Time: entry.Stamp,
		}
		// fmt.Printf("point[%d]: %v\n", i, pts[i])
	}

	bps := client.BatchPoints{
		Points:          pts,
		Database:        MyDB,
		RetentionPolicy: "default",
	}
	_, err := con.Write(bps)
	Checkerr(err)
}

// Create the client connection
func connect() (*client.Client, error) {
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", MyHost, MyPort))
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL: *u,
		// Username: os.Getenv("INFLUX_USER"),
		// Password: os.Getenv("INFLUX_PWD"),
	}

	con, err := client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	dur, ver, err := con.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to %v, InfluxDB:%s, ping:%v", u, ver, dur)

	return con, nil
}
