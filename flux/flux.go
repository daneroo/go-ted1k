package flux

import (
	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
	"github.com/influxdb/influxdb/client"
	"math/rand"
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
	writeBatchSize = 3600 * 24 //10000
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

func WriteAll(src <-chan Entry) {
	start := time.Now()
	count := 0

	con, err := connect()
	Checkerr(err)
	// defer close?

	var entries = make([]Entry, 0, writeBatchSize)

	for entry := range src {
		entries = append(entries, entry)
		count++
		if len(entries) == cap(entries) {
			log.Printf("Write::checkpoint at %d records %v", count, entry.Stamp)
			log.Printf(" entries len:%d cap %d", len(entries), cap(entries))
			writeEntries(con, entries)
			entries = make([]Entry, 0, writeBatchSize)
		}
	}
	log.Printf("Write::checkpoint at %d records", count)
	log.Printf(" entries len:%d cap %d", len(entries), cap(entries))
	writeEntries(con, entries)
	TimeTrack(start, "flux.WriteAll", count)
}

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
	log.Printf("Connected to %v, %s, ping:%v %v", u, ver, dur)

	return con, nil
}

func Try() {
	log.Println("Trying influxdb")

	con, err := connect()
	Checkerr(err)
	// writePoints(con)

	entries := makeEntries()
	writeEntries(con, entries)
}

func makeEntries() []Entry {
	const (
		wattRange  = 1000
		sampleSize = 10
	)
	var entries = make([]Entry, sampleSize)
	rand.Seed(42)
	for i := 0; i < sampleSize; i++ {
		entries[i] = Entry{
			Watt:  rand.Intn(wattRange),
			Stamp: time.Now().Round(time.Millisecond),
			// Stamp: client.SetPrecision(time.Now(), "ms"),
		}
		fmt.Printf("entry[%d]: %v\n", i, entries[i])
		time.Sleep(100 * time.Millisecond)
	}
	log.Println("Done making Points")
	return entries
}

func writeEntries(con *client.Client, entries []Entry) {
	log.Printf("Writing %d entries\n", len(entries))
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
	response, err := con.Write(bps)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Write response: %v\n", response)
}
