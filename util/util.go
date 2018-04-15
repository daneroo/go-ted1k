package util

// Just some auxilliary functions

import (
	"log"
	"time"
)

// TimeTrack calculates elapsed time as well as rate
func TimeTrack(start time.Time, name string, count int) {
	elapsed := time.Since(start)

	//TODO(daneroo) try som adaptive rounding....
	// nanos := time.Since(start).Nanoseconds()
	// secs := nanos / 1e9
	// elapsed := time.Duration(secs) * time.Second

	if count > 0 {
		rate := float64(count) / elapsed.Seconds()
		log.Printf("%s took %s, rate ~ %.1f/s count: %d", name, elapsed, rate, count)
	} else {
		log.Printf("%s took %s", name, elapsed)
	}
}

func Checkerr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
