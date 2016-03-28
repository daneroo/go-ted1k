package progress

import (
	"fmt"
	"time"

	. "github.com/daneroo/go-ted1k/types"
	. "github.com/daneroo/go-ted1k/util"
)

const (
	BatchByDay = 3600 * 24
)

type Monitor struct {
	Batch int
}

func (p *Monitor) Monitor(src <-chan Entry) <-chan Entry {
	dst := make(chan Entry)

	go func(p *Monitor) {
		start := time.Now()
		innerStart := start // so we ca track the inner loop rate
		count := 0
		for entry := range src {
			count++
			// this is the passthrough
			dst <- entry

			if (count % p.Batch) == 0 {
				TimeTrack(innerStart, "progress.Monitor.inner", p.Batch)
				// reset the inner timer
				innerStart = time.Now()

				day := entry.Stamp.Format("2006-01-02")
				TimeTrack(start, fmt.Sprintf("progress.Monitor.global (%s)", day), count)
				// TimeTrack(start, "progress.Monitor.global", count)
			}
		}
		// close the channel
		close(dst)
		TimeTrack(start, "progress.Monitor", count)
	}(p)

	return dst
}
