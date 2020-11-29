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

// Reader needs a NewContructor()...
type Reader struct {
	Conn      *pgx.Conn
	TableName string
	Epoch     time.Time
	// MaxRows is the max number rows readb from database per query (LIMIT)
	MaxRows int
	// Batch is the capacity of a single slice []types.Entry
	Batch int
	// this state is shared to preserve split of Read(),readOneFile()
	src   chan []types.Entry
	slice []types.Entry
}

const (
	channelCapacity = 2 // this is now a channel of slices
)

// NewReader is a constructor for the Reader struct
func NewReader(conn *pgx.Conn, tableName string) *Reader {
	return &Reader{
		Conn:      conn,
		TableName: tableName,
		Epoch:     AllTime,
		MaxRows:   AboutADay,
	}
}

// Read() creates and returns a channel of types.Entry
func (r *Reader) Read() <-chan []types.Entry {
	r.src = make(chan []types.Entry, channelCapacity)

	go func(r *Reader) {
		start := time.Now()
		r.slice = make([]types.Entry, 0, r.Batch)

		totalCount := 0
		startTime := r.Epoch
		for {
			lastStamp, rowCount := r.readRows(startTime)

			totalCount += rowCount
			startTime = lastStamp

			// break if there are no more rows.
			if rowCount == 0 {
				break
			}
		}
		// flush the slice
		r.src <- r.slice

		// close the channel
		close(r.src)
		r.src = nil
		timer.Track(start, "postgres.Read", totalCount)
	}(r)

	return r.src
}

// readRows reads a set of entries and sends them into the src channel.
// Returned rows starts at stamp > startTime (does not include the startTime bound).
// A maximum of maxRows rows are read.
// Return the maximum time stamp read, as well as the number of rows.
// TODO(daneroo) error handling
func (r *Reader) readRows(startTime time.Time) (time.Time, int) {
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
			// src <- types.Entry{Stamp: stamp.Time, Watt: watt}
			entry := types.Entry{Stamp: stamp.Time, Watt: watt}
			r.slice = append(r.slice, entry)
			if len(r.slice) == cap(r.slice) {
				r.src <- r.slice
				r.slice = make([]types.Entry, 0, r.Batch)
			}

		}
	}
	// TODO(daneroo) error handling
	// Any errors encountered by rows.Next or rows.Scan will be returned here
	// if rows.Err() != nil {
	// 	return rows.Err()
	// }
	return lastStamp, count
}
