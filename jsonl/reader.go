package jsonl

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/daneroo/go-ted1k/timer"
	"github.com/daneroo/go-ted1k/types"
	"github.com/daneroo/go-ted1k/util"
	"github.com/daneroo/timewalker"
)

// Reader is ...
type Reader struct {
	Grain timewalker.Duration
}

// Read() creates and returns a channel of types.Entry
func (r *Reader) Read() <-chan types.Entry {
	src := make(chan types.Entry)

	go func(r *Reader) {
		start := time.Now()

		// get the files
		filenames, err := filesIn(r.Grain)
		util.Checkerr(err)

		totalCount := 0

		for _, filename := range filenames {
			log.Printf("-jsonl.Read: %s : %s", r.Grain, filename)
			count := readOneFile(filename, src)
			totalCount += count
		}

		// close the channel
		close(src)
		timer.Track(start, "json.Read", totalCount)
	}(r)

	return src
}

// TODO(daneroo): error handling
func readOneFile(filename string, src chan<- types.Entry) int {
	// Open the file
	reader, err := os.Open(filename)
	util.Checkerr(err)

	dec := json.NewDecoder(reader)

	count := 0
	var entry types.Entry // the entry we decode into
	for dec.More() {

		// decode an array value (Message)
		err := dec.Decode(&entry)
		util.Checkerr(err)

		count++

		// send the entry into the channel
		src <- entry
	}
	return count
}
