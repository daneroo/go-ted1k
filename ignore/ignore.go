package ignore

import (
	"time"

	"github.com/daneroo/go-ted1k/timer"
	"github.com/daneroo/go-ted1k/types"
)

// Write creates an Entry channel
func Write(src <-chan types.Entry) {
	start := time.Now()
	count := 0
	for range src {
		count++
	}
	timer.Track(start, "ignore.Write", count)
}
