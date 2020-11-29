package ephemeral

import (
	"fmt"
	"time"

	"github.com/daneroo/go-ted1k/timer"
	"github.com/daneroo/go-ted1k/types"
)

const monitorBatch = 1e7

// Monitor creates a passthrough channel of Entry, which is monitored
// TODO(daneroo): add configuration and state (receiver), with Name and break behaviour
func Monitor(name string, src <-chan []types.Entry) <-chan []types.Entry {
	dst := make(chan []types.Entry)

	go func() {
		start := time.Now()
		innerStart := start // so we ca track the inner loop rate
		count := 0
		for slice := range src {
			for _, entry := range slice { // index,entry
				count++

				if (count % monitorBatch) == 0 {
					innerRate := timer.Rate(time.Since(innerStart), monitorBatch)
					// reset the inner timer
					innerStart = time.Now()

					day := entry.Stamp.Format("2006-01-02")
					msg := fmt.Sprintf("%s (%s) inner %s,", name, day, innerRate)
					timer.Track(start, msg, count)
				}
			}
			// this is the passthrough
			dst <- slice

		}

		// close the channel
		close(dst)
		timer.Track(start, name, count)
	}()

	return dst
}
