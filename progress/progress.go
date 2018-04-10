package progress

import (
	"fmt"
	"log"
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

// Gaps prints any succesive gaps in consecutive entries > 1 hour
func (p *Monitor) Gaps(src <-chan Entry) <-chan Entry {
	dst := make(chan Entry)

	go func(p *Monitor) {
		count := 0
		gaps := 0
		var sumGaps time.Duration
		var firstStamp time.Time
		var lastStamp time.Time
		for entry := range src {
			// this is the passthrough
			dst <- entry
			count++
			if lastStamp.IsZero() {
				firstStamp = entry.Stamp
			} else { // avoid first compare
				gap := entry.Stamp.Sub(lastStamp)
				if gap > time.Hour {
					gaps++
					sumGaps += gap
					log.Printf("progress.Gaps: %s %s : %s", lastStamp.Format(time.RFC3339), entry.Stamp.Format(time.RFC3339), gap)
				}
			}
			lastStamp = entry.Stamp
		}
		// close the channel
		close(dst)
		log.Printf("Progress.Gaps: %d gaps totaling %s (%d entries)", gaps, sumGaps, int(sumGaps.Seconds()))
		log.Printf("Progress.Gaps: %d total entries in [%s, %s] %s", count, firstStamp.Format(time.RFC3339), lastStamp.Format(time.RFC3339), lastStamp.Sub(firstStamp))
	}(p)

	return dst
}
