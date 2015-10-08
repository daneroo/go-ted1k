package source

import (
	"database/sql"
	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var (
	maxCountPerChunk = 3600 * 24
	epoch            = time.Date(2015, time.July, 1, 0, 0, 0, 0, time.UTC)
	// epoch            = time.Date(2015, time.September, 27, 0, 0, 0, 0, time.UTC)
	// epoch = time.Date(2007, time.January, 0, 0, 0, 0, 0, time.UTC)
	// epoch = time.Date(2037, time.January, 0, 0, 0, 0, 0, time.UTC)
)

// ReadAll creates and return a channel of Entry
func ReadAll(db *sql.DB) <-chan Entry {
	src := make(chan Entry)

	go func() {
		start := time.Now()
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
		TimeTrack(start, "source.ReadAll", rowCount)
	}()

	return src
}

func oneChunk(db *sql.DB, startTimeExcl time.Time, maxCountPerChunk int, src chan<- Entry) (time.Time, int) {
	start := time.Now()
	// defer TimeTrack(time.Now(), "oneChunk", maxCountPerChunk)
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
	TimeTrack(start, "source.ReadAll.checkpoint", chunkRowCount)
	return lastStamp, chunkRowCount
}
