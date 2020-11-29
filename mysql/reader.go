package mysql

import (
	"fmt"
	"time"

	"github.com/daneroo/go-ted1k/timer"
	"github.com/daneroo/go-ted1k/types"
	"github.com/daneroo/go-ted1k/util"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	// AboutADay is used to size the maximum number of rows read in one iteration
	AboutADay = 3600 * 24
)

var (
	// ThisYear is ...
	ThisYear = time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	// Recent is ...
	Recent = time.Date(2015, time.September, 25, 0, 0, 0, 0, time.UTC)
	// SixMonths is ...
	SixMonths = time.Date(2015, time.July, 1, 0, 0, 0, 0, time.UTC)
	// LastYear  is ...
	LastYear = time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC)
	// AllTime is ...
	AllTime = time.Date(1970, time.January, 0, 0, 0, 0, 0, time.UTC)
	// FarFuture is ...
	FarFuture = time.Date(2037, time.January, 0, 0, 0, 0, 0, time.UTC)
)

// Reader is ...
type Reader struct {
	DB        *sqlx.DB
	TableName string
	Epoch     time.Time
	MaxRows   int
}

// Read() creates and returns a channel of types.Entry
func (r *Reader) Read() <-chan types.Entry {
	src := make(chan types.Entry)

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
		timer.Track(start, "mysql.Read", totalCount)
	}(r)

	return src
}

// readRows reads a set of entries and sends them into the src channel.
// Returned rows starts at stamp > startTime (does not include the startTime bound).
// A maximum of maxRows rows are read.
// Return the maximum time stamp read, as well as the number of rows.
func (r *Reader) readRows(startTime time.Time, src chan<- types.Entry) (time.Time, int) {
	sql := fmt.Sprintf("SELECT stamp,watt FROM %s where stamp>? ORDER BY stamp ASC LIMIT ?", r.TableName)

	rows, err := r.DB.Query(sql, startTime, r.MaxRows)
	defer rows.Close()
	util.Checkerr(err)

	count := 0
	var lastStamp time.Time
	for rows.Next() {
		var stamp mysql.NullTime
		var watt int

		err = rows.Scan(&stamp, &watt)
		util.Checkerr(err)

		// count even null stamp rows (which should never happen)
		count++
		if stamp.Valid {
			lastStamp = stamp.Time
			src <- types.Entry{Stamp: stamp.Time, Watt: watt}
		}
	}
	return lastStamp, count
}
