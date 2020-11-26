package synth

import (
	"time"

	"github.com/daneroo/go-ted1k/timer"
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
	Epoch     time.Time
	TotalRows int
}

// Read() creates and returns a channel of types.Entry
func (r *Reader) Read() <-chan types.Entry {
	src := make(chan types.Entry)

	go func(r *Reader) {
		start := time.Now()

		totalCount := 0
		stamp := r.Epoch
		for {

			entry := types.Entry{Stamp: stamp, Watt: int(stamp.Unix() % 5000)}
			src <- entry

			totalCount++
			stamp = stamp.Add(time.Second)

			// break if there are no more rows.
			if totalCount >= r.TotalRows {
				break
			}
		}
		// close the channel
		close(src)
		timer.Track(start, "synth.Read", totalCount)
	}(r)

	return src
}
