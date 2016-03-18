package ignore

import (
	"log"
	"time"

	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
)

const (
	BatchByDay = 3600 * 24
)

func New(batch int) (*Ignorer, error) {
	i := &Ignorer{batch}
	if batch <= 0 {
		i.batch = BatchByDay
	}
	log.Printf("Ignorer.batch: %d", i.batch)
	return i, nil
}

type Ignorer struct {
	batch int
}

func (i Ignorer) Write(src <-chan Entry) {
	start := time.Now()
	count := 0
	for entry := range src {
		count++
		if (count % i.batch) == 0 {
			log.Printf("ignore.Write::checkpoint at %d records %v", count, entry.Stamp)
		}
	}
	TimeTrack(start, "ignore.Write", count)
}
