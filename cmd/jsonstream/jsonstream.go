package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/daneroo/go-ted1k/ephemeral"
	"github.com/daneroo/go-ted1k/jsonl"
	"github.com/daneroo/go-ted1k/logsetup"
	"github.com/daneroo/go-ted1k/progress"
	"github.com/daneroo/go-ted1k/timer"
	"github.com/daneroo/go-ted1k/types"
	"github.com/intel-go/fastjson"
)

func main() {
	logsetup.SetupFormat()
	log.Printf("Starting json stream test\n")
	// prepare data -> 900k/s
	doTest("ephemeral -> jsonl", ephemeral.NewReader(), jsonl.NewWriter())
	// reference implementation - jsonl.NewReader - 500k/s
	doTest("jsonl -> ephemeral", jsonl.NewReader(), ephemeral.NewWriter())

	// now read a single jsonl file
	filename := "data/jsonl/month/2020-01-01T00:00:00Z.jsonl"
	parseRef(filename)

	parseFastJSON(filename)
}

const (
	bufferedReaderSize = 32 * 1024 // default is 4k, 32k ~5% improvement
	batchSize          = 24 * 3600
)

// with slices ~ 520k/s - no slices ~ 548k/s
//  buf: 0k ~500k/s
//  buf: 4k ~520k/s
//  buf:32k ~548
func parseRef(filename string) int {
	start := time.Now()

	// Open the file
	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// dec := json.NewDecoder(reader)

	// bufferedReader := bufio.NewReader(reader) // 4k is the default
	// dec := json.NewDecoder(bufferedReader)

	bufferedReader := bufio.NewReaderSize(reader, bufferedReaderSize)
	dec := json.NewDecoder(bufferedReader)

	count := 0
	// slice := make([]types.Entry, 0, batchSize)
	var entry types.Entry // the entry we decode into
	for dec.More() {

		// decode an array value (Message)
		if err := dec.Decode(&entry); err != nil {
			log.Fatal(err)
		}

		count++
	}
	timer.Track(start, fmt.Sprintf("encoding/json %45s", filename), count)

	return count
}

// stolen from pump
func doTest(name string, r types.EntryReader, w types.EntryWriter) (int, error) {
	log.Printf("-=- %s\n", name)
	return w.Write(progress.Monitor(name, r.Read()))
}

func parseFastJSON(filename string) int {
	start := time.Now()

	// Open the file
	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// dec := json.NewDecoder(reader)

	// bufferedReader := bufio.NewReader(reader) // 4k is the default
	// dec := json.NewDecoder(bufferedReader)

	bufferedReader := bufio.NewReaderSize(reader, bufferedReaderSize)
	// dec := json.NewDecoder(bufferedReader)
	dec := fastjson.NewDecoder(bufferedReader)

	count := 0
	// slice := make([]types.Entry, 0, batchSize)
	var entry types.Entry // the entry we decode into

	for dec.More() {

		// decode an array value (Message)
		if err := dec.Decode(&entry); err != nil {
			// log.Fatal(err)
			log.Println(err)
			break
		}
		// log.Printf("entry: %s %d\n", entry.Stamp, entry.Watt)

		count++
	}
	timer.Track(start, fmt.Sprintf("fastjson %45s", filename), count)

	return count
}
