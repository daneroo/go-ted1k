package mysql

import (
	"log"
	"time"

	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	insertSql      = "INSERT IGNORE INTO watt2 (stamp, watt) VALUES (?,?)"
	writeBatchSize = 12 * 3600
)

type Writer struct {
	DB         *sqlx.DB
	TableName  string
	tx         *sqlx.Tx
	insertStmt *sqlx.Stmt
}

func (w *Writer) Write(src <-chan Entry) {
	start := time.Now()
	var err error
	w.tx, err = w.DB.Beginx()
	Checkerr(err)
	w.insertStmt, err = w.tx.Preparex(insertSql)
	Checkerr(err)
	log.Println("Prepared insert statement (in a transaction)")

	count := 0
	for entry := range src {

		w.writeOneRow(entry.Stamp, entry.Watt)
		// log.Printf("Write %v, %d  (%d)\n", entry.Stamp, entry.Watt, count)

		count++
		if (count % writeBatchSize) == 0 {
			TimeTrack(start, "mysql.Write.checkpoint", count)
			// log.Printf("Write::checkpoint at %d records %v", count, entry.Stamp)
			w.commitAndBeginTx(w.DB)
		}

	}

	// final Close
	w.insertStmt.Close()

	// final Tx.commit
	err = w.tx.Commit() // not quite right..
	Checkerr(err)

	TimeTrack(start, "sink.WriteAll", count)
}

func (w *Writer) commitAndBeginTx(db *sqlx.DB) {
	w.insertStmt.Close()
	var err error
	w.tx.Commit()
	w.tx, err = w.DB.Beginx()
	Checkerr(err)
	w.insertStmt, err = w.tx.Preparex(insertSql)
	Checkerr(err)
}

func (w *Writer) writeOneRow(stamp time.Time, watt int) {
	// log.Printf("Write %v, %d\n", stamp, watt)
	_, err := w.insertStmt.Exec(stamp, watt)
	Checkerr(err)
}
