package progress

import (
	"log"
	"time"

	"github.com/daneroo/go-ted1k/types"
)

// Gaps prints any succesive gaps in consecutive entries (> 1 hour)
func Gaps(src <-chan []types.Entry) <-chan []types.Entry {
	dst := make(chan []types.Entry)

	go func() {
		count := 0
		gaps := 0
		var sumGaps time.Duration
		var firstStamp time.Time
		var lastStamp time.Time
		for slice := range src {
			// this is the passthrough
			dst <- slice
			for _, entry := range slice { // index,entry
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
		}
		// close the channel
		close(dst)
		log.Printf("Progress.Gaps: %d gaps totaling %s (%d entries)", gaps, sumGaps, int(sumGaps.Seconds()))
		log.Printf("Progress.Gaps: %d total entries in [%s, %s] %s", count, firstStamp.Format(time.RFC3339), lastStamp.Format(time.RFC3339), lastStamp.Sub(firstStamp))
	}()

	return dst
}
