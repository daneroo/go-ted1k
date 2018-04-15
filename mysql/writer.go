package mysql

import (
	"bytes"
	"fmt"
	"log"
	"time"

	. "github.com/daneroo/go-ted1k/types"
	. "github.com/daneroo/go-ted1k/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	insertSQLFormat = "INSERT IGNORE INTO %s (stamp, watt) VALUES (?,?)"
	writeBatchSize  = 10000 // could move to Writer struct
)

// Writer is a ...
type Writer struct {
	DB        *sqlx.DB
	TableName string
	prepStmts map[string]*sqlx.Stmt // TODO(daneroo): these need to be closed..
}

// Close frees prepared Statements
func (w *Writer) Close() {
	for _, stmt := range w.prepStmts {
		log.Printf("Closing prepared statement")
		_ = stmt.Close()
	}
	w.prepStmts = make(map[string]*sqlx.Stmt)
}

func (w *Writer) Write(src <-chan Entry) {
	start := time.Now()
	count := 0
	entries := make([]Entry, 0, writeBatchSize)

	for entry := range src {

		entries = append(entries, entry)

		count++
		if (len(entries) % writeBatchSize) == 0 {
			w.flush(entries)
			entries = make([]Entry, 0, writeBatchSize)
		}

	}
	// last flush
	w.flush(entries)
	TimeTrack(start, "mysql.Write", count)
}

// perform the actual batch insert
func (w *Writer) flush(entries []Entry) {
	if len(entries) == 0 {
		return
	}
	// log.Printf("flush: would have flushed %d entries", len(entries))
	sql := w.makeSQL(len(entries))

	vals := []interface{}{}
	for _, entry := range entries {
		vals = append(vals, entry.Stamp, entry.Watt)
	}
	stmt := w.makeStmt(sql)
	stmt.MustExec(vals...)
	// log.Printf("res: %v", res)
}

func (w *Writer) makeStmt(sql string) *sqlx.Stmt {
	if w.prepStmts == nil {
		w.prepStmts = make(map[string]*sqlx.Stmt)
	}

	// Prepare query, if necessary
	if _, ok := w.prepStmts[sql]; !ok {
		if stmt, err := w.DB.Preparex(sql); err != nil {
			log.Println(err)
			panic(err)
		} else {
			w.prepStmts[sql] = stmt
		}
	}
	return w.prepStmts[sql]
}

// make multiple value insert sql statement
func (w *Writer) makeSQL(length int) string {
	if length == 0 {
		return ""
	}
	var sql bytes.Buffer
	sql.WriteString(fmt.Sprintf(insertSQLFormat, w.TableName))
	for i := 0; i < length-1; i++ {
		sql.WriteString(",(?,?)")
	}

	return sql.String()
}
