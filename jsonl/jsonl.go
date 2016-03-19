package jsonl

import (
	"fmt"
	"log"
	"time"

	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
	"github.com/daneroo/timewalker"
)

type Writer struct {
	FlushBoundary timewalker.Duration
	enc           FBJE
	intvl         timewalker.Interval
}

func DefaultWriter() *Writer {
	w := &Writer{
		FlushBoundary: timewalker.Day,
	}
	return w
}

// Consume the Entry (receive only) channel
// preforming batched writes (of size writeBatchSize)
// Also performs progress logging (and timing)
func (w *Writer) Write(src <-chan Entry) {
	start := time.Now()
	count := 0

	for entry := range src {
		count++

		w.openFor(entry)
		err := w.enc.Encode(&entry)
		Checkerr(err)

	}
	w.close()
	TimeTrack(start, "jsonl.Write", count)
}

func (w *Writer) close() {
	w.enc.Close()
}

// Does 4 things; open File/buffer/encoder/Interval
//
func (w *Writer) openFor(entry Entry) {
	// could test Start==End (not initialized)
	if !w.intvl.Start.IsZero() {
		// log.Printf("-I: %s : %s %s", w.FlushBoundary, entry.Stamp, w.intvl)
	} else {
		s := w.FlushBoundary.Floor(entry.Stamp)
		e := w.FlushBoundary.AddTo(s)
		w.intvl = timewalker.Interval{Start: s, End: e}
		log.Printf("+I: %s : %s %s", w.FlushBoundary, entry.Stamp, w.intvl)
	}

	if !entry.Stamp.Before(w.intvl.End) {
		if w.enc.isOpen {
			log.Printf("Should close: %s", w.intvl)
			s := w.FlushBoundary.Floor(entry.Stamp)
			e := w.FlushBoundary.AddTo(s)
			w.intvl = timewalker.Interval{Start: s, End: e}
			w.enc.Close()
		}
	}

	if !w.enc.isOpen {
		log.Printf("Should open: %s", w.intvl)
		err := w.enc.Open(fmt.Sprintf("./data/%s-%s.jsonl", w.FlushBoundary, w.intvl.Start.Format(time.RFC3339)))
		Checkerr(err)
	}

}
