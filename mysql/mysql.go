package mysql

import (
	"log"
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
	LastYear  = time.Date(2015, time.July, 1, 0, 0, 0, 0, time.UTC)
	AllTime   = time.Date(1970, time.January, 0, 0, 0, 0, 0, time.UTC)
	FarFuture = time.Date(2037, time.January, 0, 0, 0, 0, 0, time.UTC)
)

type Mysqler struct {
	db      *sqlx.DB
	epoch   time.Time
	maxRows int
}

func New(db *sqlx.DB, epoch time.Time, maxRows int) (*Mysqler, error) {
	my := &Mysqler{db: db, epoch: epoch, maxRows: maxRows}
	if maxRows <= 0 {
		my.maxRows = AboutADay
	}
	log.Printf("source.options: %v", my)
	return my, nil
}

// ReadAll creates and return a channel of Entry
func (my Mysqler) Read() <-chan Entry {
	src := make(chan Entry)

	go func(my Mysqler) {
		start := time.Now()

		totalCount := 0
		startTime := my.epoch
		for {
			lastStamp, rowCount := my.readRows(startTime, src)

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
	}(my)

	return src
}

// readRows reads a set of entries and sends them into the src channel.
// Returned rows starts at stamp > startTime (does not include the startTime bound).
// A maximum of maxRows rows are read.
// Return the maximum time stamp read, as well as the number of rows.
func (my Mysqler) readRows(startTime time.Time, src chan<- Entry) (time.Time, int) {
	start := time.Now()
	sql := "SELECT stamp,watt FROM watt where stamp>? ORDER BY stamp ASC LIMIT ?"

	rows, err := my.db.Query(sql, startTime, my.maxRows)
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
