package jsonl

import (
	"log"
	"time"

	"github.com/daneroo/go-ted1k/timer"
	"github.com/daneroo/go-ted1k/types"
	"github.com/daneroo/go-ted1k/util"
	"github.com/daneroo/timewalker"
)

// Writer is ...
type Writer struct {
	Grain    timewalker.Duration
	BasePath string
	enc      FBJE
	intvl    timewalker.Interval
}

// NewWriter is a constructor for the Writer struct
func NewWriter() *Writer {
	return &Writer{
		Grain:    timewalker.Month,
		BasePath: defaultBasePath,
	}
}

// Consume the types.Entry (receive only) channel
// preforming batched writes (of size writeBatchSize)
// Also performs progress logging (and timing)
func (w *Writer) Write(src <-chan []types.Entry) {
	start := time.Now()
	count := 0

	for slice := range src {
		for _, entry := range slice {
			count++

			w.openFor(entry)
			err := w.enc.Encode(&entry)
			util.Checkerr(err)

		}
	}
	w.close()
	timer.Track(start, "jsonl.Write", count)
}

func (w *Writer) close() {
	w.enc.Close()
}

// Does 4 things; open File, buffer, encoder, Interval
func (w *Writer) openFor(entry types.Entry) {
	// could test Start==End (not initialized)
	if !w.intvl.Start.IsZero() {
		// log.Printf("-I: %s : %s %s", w.Grain, entry.Stamp, w.intvl)
	} else {
		s := w.Grain.Floor(entry.Stamp)
		e := w.Grain.AddTo(s)
		w.intvl = timewalker.Interval{Start: s, End: e}
		log.Printf("+I: %s : %s %s", w.Grain, entry.Stamp, w.intvl)
	}

	if !entry.Stamp.Before(w.intvl.End) {
		if w.enc.isOpen {
			// log.Printf("Should close: %s", w.intvl)

			// new interval: for loop
			s := w.Grain.Floor(entry.Stamp)
			e := w.Grain.AddTo(s)
			w.intvl = timewalker.Interval{Start: s, End: e}

			w.enc.Close()
		}
	}

	if !w.enc.isOpen {
		// log.Printf("Should open: %s", w.intvl)

		// this make directories as well...
		file, err := pathFor(w.BasePath, w.Grain, w.intvl)
		util.Checkerr(err)

		err = w.enc.Open(file)
		util.Checkerr(err)
	}

}
