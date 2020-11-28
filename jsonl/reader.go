package jsonl

import (
	"bufio"
	"encoding/json"
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

const (
	channelCapacity    = 100       // this made a huge difference, 10 is not enough 1000 makesno difference
	bufferedReaderSize = 32 * 1024 // default is 4k, 32k ~5% improvement
)

// Read() creates and returns a channel of types.Entry
func (r *Reader) Read() <-chan types.Entry {
	// TODO(daneroo) tweak this capacity, probably related to the efficiency of the encoder
	src := make(chan types.Entry, channelCapacity)

	go func(r *Reader) {
		start := time.Now()

		// get the files
		filenames, err := filesIn(r.Grain)
		util.Checkerr(err)

		totalCount := 0

		for _, filename := range filenames {
			// log.Printf("-jsonl.Read: %s : %s", r.Grain, filename)
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
// TODO(daneroo): close readers and decoder
// add a bufio.NewReader
func readOneFile(filename string, src chan<- types.Entry) int {
	// Open the file
	reader, err := os.Open(filename)
	util.Checkerr(err)

	// dec := json.NewDecoder(reader)
	// bufferedReader := bufio.NewReader(reader) // 4k is the default
	bufferedReader := bufio.NewReaderSize(reader, bufferedReaderSize)
	dec := json.NewDecoder(bufferedReader)

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
