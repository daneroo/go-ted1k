package jsonl

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/daneroo/go-ted1k/types"
	"github.com/daneroo/go-ted1k/util"
	"github.com/daneroo/timewalker"
)

// Reader is ...
type Reader struct {
	Grain    timewalker.Duration
	BasePath string
	// Slice Batching: the capacity of a single slice []types.Entry
	Batch int
	// this state is shared to preserve split of Read(),readOneFile()
	src   chan []types.Entry
	slice []types.Entry
}

const (
	// channelCapacity    = 100       // this made a huge difference, 10 is not enough 1000 makesno difference
	channelCapacity    = 2         // this is now a channel of slices
	bufferedReaderSize = 32 * 1024 // default is 4k, 32k ~5% improvement
)

// NewReader is a constructor for the Reader struct
func NewReader() *Reader {
	return &Reader{
		Grain:    timewalker.Month,
		BasePath: defaultBasePath,
		Batch:    1000,
	}
}

// Read() creates and returns a channel of []types.Entry
func (r *Reader) Read() <-chan []types.Entry {
	// TODO(daneroo) tweak this capacity, probably related to the efficiency of the encoder
	r.src = make(chan []types.Entry, channelCapacity)

	go func(r *Reader) {
		// start := time.Now()
		r.slice = make([]types.Entry, 0, r.Batch)

		// get the files
		filenames, err := filesIn(r.BasePath, r.Grain)
		util.Checkerr(err)

		totalCount := 0

		for _, filename := range filenames {
			count := r.readOneFile(filename)
			totalCount += count
		}

		// flush the slice
		r.src <- r.slice

		// close the channel
		close(r.src)
		r.src = nil
		// timer.Track(start, "jsonl.Read", totalCount)
	}(r)

	return r.src
}

// TODO(daneroo): error handling
// TODO(daneroo): close readers and decoder
// add a bufio.NewReader
func (r *Reader) readOneFile(filename string) int {
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

		// appends the entry to the slice
		r.slice = append(r.slice, entry)
		// send the slice to te channel
		if len(r.slice) == cap(r.slice) {
			r.src <- r.slice
			r.slice = make([]types.Entry, 0, r.Batch)
		}

	}
	return count
}
