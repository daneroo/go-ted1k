package jsonl

import (
	"encoding/json"
	"log"
	"os"
	"time"

	. "github.com/daneroo/go-ted1k/types"
	. "github.com/daneroo/go-ted1k/util"
	"github.com/daneroo/timewalker"
)

type Reader struct {
	Grain timewalker.Duration
}

// Read() creates and returns a channel of Entry
func (r *Reader) Read() <-chan Entry {
	src := make(chan Entry)

	go func(r *Reader) {
		start := time.Now()

		// get the files
		filenames, err := filesIn(r.Grain)
		Checkerr(err)

		totalCount := 0

		for _, filename := range filenames {
			log.Printf("-jsonl.Read: %s : %s", r.Grain, filename)
			count := readOneFile(filename, src)
			totalCount += count
		}

		// close the channel
		close(src)
		TimeTrack(start, "json.Read", totalCount)
	}(r)

	return src
}

// TODO(daneroo): error handling
func readOneFile(filename string, src chan<- Entry) int {
	// Open the file
	reader, err := os.Open(filename)
	Checkerr(err)

	dec := json.NewDecoder(reader)

	count := 0
	var entry Entry // the entry we decode into
	for dec.More() {

		// decode an array value (Message)
		err := dec.Decode(&entry)
		Checkerr(err)

		count++

		// send the entry into the channel
		src <- entry
	}
	return count
}
