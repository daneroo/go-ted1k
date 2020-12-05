package ephemeral

import (
	"time"

	"github.com/daneroo/go-ted1k/types"
)

var (
	// ThisYear is a reference date
	ThisYear = time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	// LastYear is a reference date
	LastYear = time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
)

// Reader contains the initialized params
type Reader struct {
	Epoch time.Time
	// Slice Batching, 1000 is better if useful, 100 performant, 10 is still ok.
	Batch     int
	TotalRows int
}

// NewReader is a constructor for the Reader struct
func NewReader() *Reader {
	return &Reader{
		Epoch: ThisYear,
		Batch: 1000,
		// math.PI * 1e7 ~ 1 year in seconds!
		// TotalRows: 3.1415926e7,
		TotalRows: 1e7,
	}
}

// Read() creates and returns a channel of []types.Entry
func (r *Reader) Read() <-chan []types.Entry {
	src := make(chan []types.Entry)

	go func(r *Reader) {
		slice := make([]types.Entry, 0, r.Batch)

		totalCount := 0
		stamp := r.Epoch
		for {

			entry := types.Entry{Stamp: stamp, Watt: int(stamp.Unix() % 5000)}
			slice = append(slice, entry)
			if len(slice) == cap(slice) {
				src <- slice
				slice = make([]types.Entry, 0, r.Batch)
			}

			totalCount++
			stamp = stamp.Add(time.Second)

			// break if there are no more rows.
			if totalCount >= r.TotalRows {
				break
			}
		}
		// flush the slice
		src <- slice
		// close the channel
		close(src)
	}(r)

	return src
}
