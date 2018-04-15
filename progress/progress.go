package progress

import (
	"fmt"
	"log"
	"time"

	"github.com/daneroo/go-ted1k/timer"
	"github.com/daneroo/go-ted1k/types"
)

const (
	// BatchByDay is an approximate count of samples per day
	BatchByDay = 3600 * 24
)

// Monitor represents a tracking channel monitor
type Monitor struct {
	Batch int
}

// Monitor creates a passthrough channel of Entry, which is monitored
func (p *Monitor) Monitor(src <-chan types.Entry) <-chan types.Entry {
	dst := make(chan types.Entry)

	go func(p *Monitor) {
		start := time.Now()
		innerStart := start // so we ca track the inner loop rate
		count := 0
		for entry := range src {
			count++
			// this is the passthrough
			dst <- entry

			if (count % p.Batch) == 0 {
				innerRate := timer.Rate(time.Since(innerStart), p.Batch)
				// reset the inner timer
				innerStart = time.Now()

				day := entry.Stamp.Format("2006-01-02")
				msg := fmt.Sprintf("progress.Monitor.global (%s) inner %s,", day, innerRate)
				timer.Track(start, msg, count)
			}
		}
		// close the channel
		close(dst)
		timer.Track(start, "progress.Monitor", count)
	}(p)

	return dst
}

// Gaps prints any succesive gaps in consecutive entries > 1 hour
func (p *Monitor) Gaps(src <-chan types.Entry) <-chan types.Entry {
	dst := make(chan types.Entry)

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
