package flux

// TODO(daneroo) should use the v2 client, it has close
import (
	"log"
	"time"

	. "github.com/daneroo/go-ted1k/types"
	. "github.com/daneroo/go-ted1k/util"
	// client "github.com/influxdata/influxdb/client/v2"
	client "github.com/influxdata/influxdb/client/v2"
)

import "fmt"

// "os"

type Writer struct {
	Host           string
	Port           int
	DB             string
	Measurement    string
	WriteBatchSize int
	con            client.Client
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
// and closes the connection
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
	w.close()
	TimeTrack(start, "flux.Write", count)
}

// Write out the entries to con, and reallocate a new empty slice
func (w *Writer) flush(entries []Entry) []Entry {
	w.writeEntries(entries)
	return make([]Entry, 0, w.WriteBatchSize)
}

// Perform the batch write
func (w *Writer) writeEntries(entries []Entry) {
	bps, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:        w.DB,
		RetentionPolicy: "default",
		Precision:       "s",
		// WriteConsistency: string,
	})
	Checkerr(err)

	for _, entry := range entries {
		name := "watt" // Measurement
		tags := map[string]string{ /*"ted1k",...*/ }
		fields := map[string]interface{}{
			"value": entry.Watt,
		}
		pt, err := client.NewPoint(name, tags, fields, entry.Stamp)
		Checkerr(err)
		bps.AddPoint(pt)

		// fmt.Printf("point: %v\n", pt)
	}

	// TODO(daneroo): retry, if error is timeout?
	err = w.con.Write(bps)
	Checkerr(err)
}

// Create the client connection
// TODO(daneroo): We need to close this thing!
func (w *Writer) connect() error {

	url := fmt.Sprintf("http://%s:%d", w.Host, w.Port)
	var err error
	w.con, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: url,
		// Username: os.Getenv("INFLUX_USER"),
		// Password: os.Getenv("INFLUX_PWD"),
	})
	Checkerr(err)

	dur, ver, err := w.con.Ping(time.Minute)
	Checkerr(err)
	log.Printf("Connected to %s, InfluxDB:%s, ping:%v", url, ver, dur)

	return nil
}
func (w *Writer) close() error {
	return w.con.Close()
}
