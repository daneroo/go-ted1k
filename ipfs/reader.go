package ipfs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/daneroo/go-ted1k/types"
	"github.com/daneroo/go-ted1k/util"
	"github.com/daneroo/timewalker"
	shell "github.com/ipfs/go-ipfs-api"
)

// Reader is ...
type Reader struct {
	sh    *shell.Shell
	Grain timewalker.Duration
	Cid   string
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
func NewReader(sh *shell.Shell, cid string) *Reader {
	return &Reader{
		sh:    sh,
		Grain: timewalker.Month,
		Batch: 1000,
		Cid:   cid,
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
		filenames, err := r.filesIn()
		util.Checkerr(err)

		totalCount := 0

		for _, filename := range filenames {
			// log.Printf("Reading %s\n", filename)
			count := r.readOneFile(filename)
			totalCount += count
		}

		// flush the slice
		r.src <- r.slice

		// close the channel
		close(r.src)
		r.src = nil
		// timer.Track(start, "ipfs.Read", totalCount)
	}(r)

	return r.src
}

// fileIn returns a slice of full paths to the files in the top level r.Cid
// assume flat directory for now
func (r *Reader) filesIn() ([]string, error) {
	path := fmt.Sprintf("%s/%s", r.Cid, dirFor(r.Grain))
	objects, err := r.sh.FileList(path)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, link := range objects.Links {
		// log.Printf("Link %d: %+v\n", idx, link)
		if link.Type == "File" { // File or Directory
			filename := fmt.Sprintf("%s/%s", path, link.Name)
			filenames = append(filenames, filename)
		} else {
			// when we want to recurse ...
		}
	}
	// Important: We must guarantee file order
	sort.Strings(filenames)

	return filenames, nil
}

// TODO(daneroo): error handling
// filename is the fully qualified pah (cid included)
func (r *Reader) readOneFile(path string) int {

	reader, err := r.sh.Cat(path)
	util.Checkerr(err)
	defer reader.Close()

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
