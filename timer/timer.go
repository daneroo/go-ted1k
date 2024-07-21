package timer

// Just some auxilliary functions

import (
	"fmt"
	"log"
	"math"
	"time"
)

// Track calculates elapsed time as well as rate
// e.g.: progress.Monitor.inner took 242ms, rate ~ 357024.8/s count: 86400
func Track(start time.Time, name string, count int) {
	elapsed := time.Since(start)
	log.Println(format(elapsed, name, count))
}

func format(elapsed time.Duration, name string, count int) string {
	// Round elapsed to Milliseconds for display (but after rate calc)
	elapsed = elapsed.Round(time.Millisecond)
	rate := Rate(elapsed, count)
	return fmt.Sprintf("%s took %s, %s count: %d", name, elapsed, rate, count)
}

// Rate formats count/elapsed as a string
func Rate(elapsed time.Duration, count int) string {
	rate := float64(count) / elapsed.Seconds()
	if math.IsNaN(rate) { // onlu happns if count and elapsed are both 0
		rate = 0
	}
	units := "/s"
	if math.IsInf(rate, 0) {
		// leave as is
	} else if rate > 1e6 {
		rate /= 1e6
		units = "M/s"
	} else if rate > 1e3 {
		rate /= 1e3
		units = "k/s"
	}

	return fmt.Sprintf("rate ~ %.1f%s", rate, units)
}
