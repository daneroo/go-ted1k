package mysql

import (
	"fmt"
	"time"

	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	insertSqlFormat = "INSERT IGNORE INTO %s (stamp, watt) VALUES (?,?)"
	writeBatchSize  = 12 * 3600
)

type Writer struct {
	DB         *sqlx.DB
	TableName  string
	tx         *sqlx.Tx
	insertStmt *sqlx.Stmt
}

func (w *Writer) Write(src <-chan Entry) {
	start := time.Now()
	w.commitAndBeginTx(true)

	count := 0
	for entry := range src {

		w.writeOneRow(entry.Stamp, entry.Watt)
		// log.Printf("Write %v, %d  (%d)\n", entry.Stamp, entry.Watt, count)

		count++
		if (count % writeBatchSize) == 0 {
			// log.Printf("Write::checkpoint at %d records %v", count, entry.Stamp)
			w.commitAndBeginTx(true)
		}

	}

	// commit but don't start another transaction
	w.commitAndBeginTx(false)
	TimeTrack(start, "mysql.Write", count)
}

// close stmt, commit, then start a tx, and prepare stmt
func (w *Writer) commitAndBeginTx(beginAgain bool) {
	if w.insertStmt != nil {
		w.insertStmt.Close()
		w.insertStmt = nil
		// log.Println("Closed statement")
	}
	if w.tx != nil {
		w.tx.Commit()
		w.tx = nil
		// log.Println("Commited transaction")
	}
	if beginAgain {
		var err error
		w.tx, err = w.DB.Beginx()
		Checkerr(err)
		insertSql := fmt.Sprintf(insertSqlFormat, w.TableName)
		w.insertStmt, err = w.tx.Preparex(insertSql)
		Checkerr(err)
		// log.Println("Prepared insert statement (in a transaction)")
	}
}

func (w *Writer) writeOneRow(stamp time.Time, watt int) {
	// log.Printf("Write %v, %d\n", stamp, watt)
	_, err := w.insertStmt.Exec(stamp, watt)
	Checkerr(err)
}
