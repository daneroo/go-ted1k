package types

import (
	"time"
)

// Entry is the basic datatype for an energy measurement
type Entry struct {
	Stamp time.Time `json:"stamp"`
	Watt  int       `json:"watt"`
}

// EntryReader is a way to produce an input channel of Entry slices
type EntryReader interface {
	Read() <-chan []Entry
}

// EntryWriter is a way to produce an output channel of Entry slices
type EntryWriter interface {
	Write(src <-chan []Entry) (int, error)
}
