package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/daneroo/go-ted1k/timer"
	"github.com/daneroo/go-ted1k/types"
	"github.com/daneroo/go-ted1k/util"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4"
)

const (
	// AboutADay is used to size the maximum number of rows read in one iteration
	AboutADay = 3600 * 24
)

var (
	ThisYear  = time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	Recent    = time.Date(2015, time.September, 25, 0, 0, 0, 0, time.UTC)
	SixMonths = time.Date(2015, time.July, 1, 0, 0, 0, 0, time.UTC)
	LastYear  = time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC)
	AllTime   = time.Date(1970, time.January, 0, 0, 0, 0, 0, time.UTC)
	FarFuture = time.Date(2037, time.January, 0, 0, 0, 0, 0, time.UTC)
)

type Reader struct {
	Conn      *pgx.Conn
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
		timer.Track(start, "postgres.Read", totalCount)
	}(r)

	return src
}

// readRows reads a set of entries and sends them into the src channel.
// Returned rows starts at stamp > startTime (does not include the startTime bound).
// A maximum of maxRows rows are read.
// Return the maximum time stamp read, as well as the number of rows.
// TODO(daneroo) error handling
func (r *Reader) readRows(startTime time.Time, src chan<- types.Entry) (time.Time, int) {
	sql := fmt.Sprintf("SELECT stamp,watt FROM %s where stamp>$1 ORDER BY stamp ASC LIMIT $2", r.TableName)

	rows, err := r.Conn.Query(context.Background(), sql, startTime, r.MaxRows)
	// TODO(daneroo) error handling
	// if err != nil {
	// 	return err
	// }
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
	// TODO(daneroo) error handling
	// Any errors encountered by rows.Next or rows.Scan will be returned here
	// if rows.Err() != nil {
	// 	return rows.Err()
	// }
	return lastStamp, count
}
