package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/daneroo/go-ted1k/timer"
	"github.com/daneroo/go-ted1k/types"
	"github.com/jackc/pgx/v4"

	// register mysql driver
	_ "github.com/go-sql-driver/mysql"
)

const (
	writeBatchSize = 10000 // could move to Writer struct
)

// Writer is a ...
type Writer struct {
	Conn      *pgx.Conn
	TableName string
}

// Close frees prepared Statements
func (w *Writer) Close() {
	log.Printf("Closing things.. (not connection")
}

func (w *Writer) Write(src <-chan types.Entry) {
	start := time.Now()
	count := 0
	entries := make([]types.Entry, 0, writeBatchSize)

	for entry := range src {

		entries = append(entries, entry)

		count++
		if (len(entries) % writeBatchSize) == 0 {
			w.flush(entries)
			entries = make([]types.Entry, 0, writeBatchSize)
		}

	}
	// last flush
	w.flush(entries)
	timer.Track(start, "mysql.Write", count)
}

// perform the actual batch insert
func (w *Writer) flush(entries []types.Entry) {
	w.writeWithMultipleInsert(entries)
}

func (w *Writer) writeWithMultipleInsert(entries []types.Entry) {
	if len(entries) == 0 {
		return
	}

	// log.Printf("flush: would have flushed %d entries", len(entries))
	sql := w.makeSQL(len(entries))

	// flatten the parameters into a single value array
	vals := []interface{}{}
	for _, entry := range entries {
		vals = append(vals, entry.Stamp, entry.Watt)
	}

	// commandTag, err := w.Conn.Exec(context.Background(), sql, vals...)
	_, err := w.Conn.Exec(context.Background(), sql, vals...)
	if err != nil {
		log.Printf("Unable to execute write: %v\n", err)
	}
}

// make multiple value insert sql statement
func (w *Writer) makeSQL(length int) string {
	if length == 0 {
		return ""
	}
	var sql strings.Builder
	insertSQLFormat := "INSERT INTO %s (stamp, watt) VALUES ($1,$2)"
	onConflict := " ON CONFLICT (stamp) DO NOTHING"

	sql.WriteString(fmt.Sprintf(insertSQLFormat, w.TableName))
	for i := 1; i < length; i++ {
		sql.WriteString(fmt.Sprintf(",($%d,$%d)", i*2+1, i*2+2))
	}
	sql.WriteString(onConflict)

	return sql.String()
}
