package source

import (
	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

var (
	maxCountPerChunk = 3600 * 24
	// epoch            = time.Date(2015, time.July, 1, 0, 0, 0, 0, time.UTC)
	epoch = time.Date(2015, time.September, 27, 0, 0, 0, 0, time.UTC)
	// epoch = time.Date(2007, time.January, 0, 0, 0, 0, 0, time.UTC)
	// epoch = time.Date(2037, time.January, 0, 0, 0, 0, 0, time.UTC)
)

// ReadAll creates and return a channel of Entry
func ReadAll(db *sqlx.DB) <-chan Entry {
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

func oneChunk(db *sqlx.DB, startTimeExcl time.Time, maxCountPerChunk int, src chan<- Entry) (time.Time, int) {
	log.Println("one chunk")
	start := time.Now()
	sql := "SELECT stamp,watt FROM watt where stamp>? ORDER BY stamp ASC LIMIT ?"

	rows, err := db.Query(sql, startTimeExcl, maxCountPerChunk)
	Checkerr(err)

	chunkRowCount := 0
	var lastStamp time.Time
	for rows.Next() {
		var stamp mysql.NullTime
		var watt int

		err = rows.Scan(&stamp, &watt)
		Checkerr(err)

		chunkRowCount++
		if stamp.Valid {
			lastStamp = stamp.Time
			src <- Entry{Stamp: stamp.Time, Watt: watt}
		}
	}
	TimeTrack(start, "source.ReadAll.checkpoint", chunkRowCount)
	return lastStamp, chunkRowCount
}
