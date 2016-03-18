package sink

import (
	"log"
	"time"

	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
	"github.com/daneroo/go-mysqltest/vendor/github.com/jmoiron/sqlx"
)

const (
	// insertSql      = "INSERT IGNORE INTO watt2 (stamp, watt) VALUES (?,?)"
	writeBatchSize = 12 * 3600
)

var (
	tx         *sqlx.Tx
	insertStmt *sqlx.Stmt
)

func WriteAll(db *sqlx.DB, src <-chan Entry) {
	start := time.Now()
	var err error
	tx, err = db.Beginx()
	Checkerr(err)
	insertStmt, err = tx.Preparex(insertSql)
	Checkerr(err)
	log.Println("Prepared insert statement (in a transaction)")

	count := 0
	for entry := range src {

		writeOneRow(entry.Stamp, entry.Watt)
		// log.Printf("Write %v, %d  (%d)\n", entry.stamp, entry.watt, count)

		count++
		if (count % writeBatchSize) == 0 {
			log.Printf("Write::checkpoint at %d records %v", count, entry.Stamp)
			commitAndBeginTx(db)
		}

	}

	// final Close
	insertStmt.Close()

	// final Tx.commit
	err = tx.Commit() // not quite right..
	Checkerr(err)

	TimeTrack(start, "sink.WriteAll", count)
}

func commitAndBeginTx(db *sqlx.DB) {
	insertStmt.Close()
	var err error
	tx.Commit()
	tx, err = db.Beginx()
	Checkerr(err)
	insertStmt, err = tx.Preparex(insertSql)
	Checkerr(err)
}

func writeOneRow(stamp time.Time, watt int) {
	// log.Printf("Write %v, %d\n", stamp, watt)
	_, err := insertStmt.Exec(stamp, watt)
	Checkerr(err)
}
