package ephemeral

import (
	"time"

	"github.com/daneroo/go-ted1k/timer"
	"github.com/daneroo/go-ted1k/types"
)

// Writer is ...
type Writer struct {
}

// NewWriter is a constructor for the Writer struct
func NewWriter() *Writer {
	return &Writer{}
}

// Write creates an Entry channel
func (w *Writer) Write(src <-chan []types.Entry) {
	start := time.Now()
	count := 0
	for slice := range src {
		for range slice { // index,entry
			count++
		}
	}
	for range src {
	}
	timer.Track(start, "ephemeral.Write", count)
}
