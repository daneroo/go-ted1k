package jsonl

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
)

const (
	// BySecond = "2006-01-02T15:04:05"
	// These are time.Format layouts, used to detect file rollover...
	ByMinute = "2006-01-02T15:04:05"
	ByHour   = "2006-01-02T15:04:05"
	ByDay    = "2006-01-02T15:04:05"
)

type Writer struct {
	FlushBoundary string
}

func DefaultWriter() *Writer {
	w := &Writer{
		FlushBoundary: ByMinute,
	}
	return w
}

// Consume the Entry (receive only) channel
// preforming batched writes (of size writeBatchSize)
// Also performs progress logging (and timing)
func (w *Writer) Write(src <-chan Entry) {
	start := time.Now()
	count := 0

	f, err := os.Create("./data/data.jsonl")
	Checkerr(err)
	defer f.Close()

	enc := json.NewEncoder(f)
	bufw := bufio.NewWriter(f) // default size 4k
	for entry := range src {
		count++

		// bytes, err := json.Marshal(entry)
		// Checkerr(err)
		// written, err := bufw.Write(bytes)
		// bufw.WriteByte('\n')
		// log.Printf("line: %s (%d, %d)", bytes, written, count)
		// Checkerr(err)

		err := enc.Encode(&entry)
		Checkerr(err)

	}
	bufw.Flush()
	TimeTrack(start, "jsonl.Write", count)
}
