package mysql

import (
	"fmt"
	"time"

	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	// maximum number of rows read in one iteration
	AboutADay = 3600 * 24
)

var (
	ThisYear  = time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	SixMonths = time.Date(2015, time.July, 1, 0, 0, 0, 0, time.UTC)
	LastYear  = time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC)
	AllTime   = time.Date(1970, time.January, 0, 0, 0, 0, 0, time.UTC)
	FarFuture = time.Date(2037, time.January, 0, 0, 0, 0, 0, time.UTC)
)

type Reader struct {
	DB        *sqlx.DB
	TableName string
	Epoch     time.Time
	MaxRows   int
}

// Read() creates and returns a channel of Entry
func (r *Reader) Read() <-chan Entry {
	src := make(chan Entry)

	go func(r *Reader) {
		start := time.Now()

		totalCount := 0
		startTime := r.Epoch
		for {
			lastStamp, rowCount := r.readRows(startTime, src)

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
	}(r)

	return src
}

// readRows reads a set of entries and sends them into the src channel.
// Returned rows starts at stamp > startTime (does not include the startTime bound).
// A maximum of maxRows rows are read.
// Return the maximum time stamp read, as well as the number of rows.
func (r *Reader) readRows(startTime time.Time, src chan<- Entry) (time.Time, int) {
	start := time.Now()
	sql := fmt.Sprintf("SELECT stamp,watt FROM %s where stamp>? ORDER BY stamp ASC LIMIT ?", r.TableName)

	rows, err := r.DB.Query(sql, startTime, r.MaxRows)
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
