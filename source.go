package main

import (
	"database/sql"
	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
)

// readAll creates and return a channel of Entry
func readAll() <-chan Entry {
	src := make(chan Entry)

	go func() {
		rowCount := 0
		startTimeExcl := epoch
		for {
			chunkRowCount := 0
			startTimeExcl, chunkRowCount = oneChunk(db, startTimeExcl, maxCountPerChunk, src)
			rowCount += chunkRowCount

			if chunkRowCount == 0 {
				break
			}
		}
		close(src)
		log.Printf("Fetched a total of %d rows", rowCount)
	}()

	return src
}

func oneChunk(db *sql.DB, startTimeExcl time.Time, maxCountPerChunk int, src chan<- Entry) (time.Time, int) {
	defer TimeTrack(time.Now(), "oneChunk", maxCountPerChunk)
	sql := "SELECT stamp,watt FROM watt where stamp>? ORDER BY stamp ASC LIMIT ?"
	// sql := "SELECT stamp,watt FROM watt where stamp<? ORDER BY stamp DESC LIMIT ?"
	rows, err := db.Query(sql, startTimeExcl, maxCountPerChunk)
	Checkerr(err)
	defer rows.Close()

	avgWatt := 0
	chunkRowCount := 0
	var lastStamp time.Time
	for rows.Next() {
		// var stamp string
		var stamp mysql.NullTime
		var watt int
		err = rows.Scan(&stamp, &watt)
		if err != nil {
			log.Println(err)
		}
		avgWatt += watt
		chunkRowCount++
		if stamp.Valid {
			lastStamp = stamp.Time
			src <- Entry{Stamp: stamp.Time, Watt: watt}
		}
		// log.Printf(" %v: %v", stamp, watt)
	}
	err = rows.Err() // get any error encountered during iteration
	Checkerr(err)

	if chunkRowCount != 0 {
		avgWatt /= chunkRowCount
	}
	log.Printf("average between (%v - %v]: %v (%v)", startTimeExcl, lastStamp, avgWatt, chunkRowCount)
	return lastStamp, chunkRowCount
}
