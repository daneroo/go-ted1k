package source

import (
	"time"

	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	// maximum number of rows read in one iteration
	maxRows = 3600 * 24
	epoch   = time.Date(2015, time.July, 1, 0, 0, 0, 0, time.UTC)
	// epoch   = time.Date(2015, time.September, 27, 0, 0, 0, 0, time.UTC)
	// epoch = time.Date(2007, time.January, 0, 0, 0, 0, 0, time.UTC)
	// epoch = time.Date(2037, time.January, 0, 0, 0, 0, 0, time.UTC)
)

// ReadAll creates and return a channel of Entry
func ReadAll(db *sqlx.DB) <-chan Entry {
	src := make(chan Entry)

	go func() {
		start := time.Now()
		totalCount := 0
		startTime := epoch
		for {
			lastStamp, rowCount := readRows(db, startTime, maxRows, src)

			totalCount += rowCount
			startTime = lastStamp

			// break if there are no more rows.
			if rowCount == 0 {
				break
			}
		}
		// close the channel
		close(src)
		TimeTrack(start, "source.ReadAll", totalCount)
	}()

	return src
}

// readRows reads a set of entries and sends them into the src channel.
// Returned rows starts at stamp > startTime (does not include the startTime bound).
// A maximum of maxRows rows are read.
// Return the maximum time stamp read, as well as the number of rows.
func readRows(db *sqlx.DB, startTime time.Time, maxRows int, src chan<- Entry) (time.Time, int) {
	start := time.Now()
	sql := "SELECT stamp,watt FROM watt where stamp>? ORDER BY stamp ASC LIMIT ?"

	rows, err := db.Query(sql, startTime, maxRows)
	defer rows.Close()
	Checkerr(err)

	count := 0
	var lastStamp time.Time
	for rows.Next() {
		var stamp mysql.NullTime
		var watt int

		err = rows.Scan(&stamp, &watt)
		Checkerr(err)

		// count even null stamp rows (which should never happen)
		count++
		if stamp.Valid {
			lastStamp = stamp.Time
			src <- Entry{Stamp: stamp.Time, Watt: watt}
		}
	}
	TimeTrack(start, "source.ReadAll.checkpoint", count)
	return lastStamp, count
}
