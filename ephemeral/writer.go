package ephemeral

import (
	"github.com/daneroo/go-ted1k/types"
)

// Writer is ...
type Writer struct {
}

// NewWriter is a constructor for the Writer struct
func NewWriter() *Writer {
	return &Writer{}
}

// Write consumes an Entry channel - returns (count,error)
func (w *Writer) Write(src <-chan []types.Entry) (int, error) {
	count := 0
	for slice := range src {
		for range slice { // index,entry
			count++
		}
	}
	return count, nil
}
