package progress

import (
	"time"

	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
)

const (
	BatchByDay = 3600 * 24
)

type Monitor struct {
	Batch int
}

func (p *Monitor) Pipe(src <-chan Entry) <-chan Entry {
	dst := make(chan Entry)

	go func(p *Monitor) {
		start := time.Now()
		count := 0
		for entry := range src {
			count++
			// this is the passthrough
			dst <- entry

			if (count % p.Batch) == 0 {
				TimeTrack(start, "progress.Pipe.checkpoint", count)
			}
		}
		// close the channel
		close(dst)
		TimeTrack(start, "progress.Pipe", count)
	}(p)

	return dst
}
