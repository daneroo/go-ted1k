package main

// Just some auxilliary functions

import (
	"log"
	"time"
)

func timeTrack(start time.Time, name string, count int) {
	elapsed := time.Since(start)
	if count > 0 {
		rate := float64(count) / elapsed.Seconds()
		log.Printf("%s took %s, rate ~ %.1f/s", name, elapsed, rate)
	} else {
		log.Printf("%s took %s", name, elapsed)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
