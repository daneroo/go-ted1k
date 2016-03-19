package flux

// TODO(daneroo) should use the v2 client, it has close
import (
	"time"

	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
	"github.com/influxdb/influxdb/client"
)

import (
	"fmt"
	"log"
	"net/url"
	// "os"
)

type Writer struct {
	Host           string
	Port           int
	DB             string
	Measurement    string
	WriteBatchSize int
	con            *client.Client
}

func DefaultWriter() *Writer {
	w := &Writer{
		Host:           "docker",
		Port:           8086,
		DB:             "ted",
		Measurement:    "watt",
		WriteBatchSize: 3600 * 24,
	}
	return w
}

// Consume the Entry (receive only) channel
// preforming batched writes (of size writeBatchSize)
// Also performs progress logging (and timing)
func (w *Writer) Write(src <-chan Entry) {
	start := time.Now()
	count := 0

	// should I close if not nil?
	err := w.connect()
	Checkerr(err)
	// defer close? when we move to v2 client

	var entries = make([]Entry, 0, w.WriteBatchSize)
	for entry := range src {
		entries = append(entries, entry)
		count++
		if len(entries) == cap(entries) {
			entries = w.flush(entries)
		}
	}
	_ = w.flush(entries)
	TimeTrack(start, "flux.Write", count)
}

// Write out the entries to con, and reallocate a new empty slice
func (w *Writer) flush(entries []Entry) []Entry {
	w.writeEntries(entries)
	return make([]Entry, 0, w.WriteBatchSize)
}

// Perform the batch write
func (w *Writer) writeEntries(entries []Entry) {
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
		Database:        w.DB,
		RetentionPolicy: "default",
	}
	_, err := w.con.Write(bps)
	Checkerr(err)
}

// Create the client connection
func (w *Writer) connect() error {
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", w.Host, w.Port))
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL: *u,
		// Username: os.Getenv("INFLUX_USER"),
		// Password: os.Getenv("INFLUX_PWD"),
	}

	w.con, err = client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	dur, ver, err := w.con.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to %v, InfluxDB:%s, ping:%v", u, ver, dur)

	return nil
}
