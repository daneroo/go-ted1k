package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/daneroo/go-ted1k/timer"
	"github.com/daneroo/go-ted1k/types"
	"github.com/jackc/pgconn"
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
	if len(entries) == 0 {
		return
	}
	// w.writeWithMultipleInsert(entries)
	w.writeWithCopyFrom(entries)
}

// writeWithCopyFrom is the fastest way to insert, but has no "ON CONFLICT (stamp) DO NOTHING" mechanism,
// so it falls back to writeWithMultipleInsert when we can assert the specific error
// TODO(daneroo) move error handling and fallback to flush() method, or new wrapper
func (w *Writer) writeWithCopyFrom(entries []types.Entry) {
	rows := [][]interface{}{}
	for _, entry := range entries {
		rows = append(rows, []interface{}{entry.Stamp, entry.Watt})
	}

	copyCount, err := w.Conn.CopyFrom(
		context.Background(),
		pgx.Identifier{w.TableName},
		[]string{"stamp", "watt"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		if pge, ok := err.(*pgconn.PgError); ok {
			// We know we have a pgconn.PgError
			// pgconn.PgError: ERROR: duplicate key value violates unique constraint "watt_pkey" (SQLSTATE 23505)
			if pge.Code == "23505" {
				// log.Printf("Retrying withMultipleInsert: Code: %v TableName: %v ConstraintName: %v\n", pge.Code, pge.TableName, pge.ConstraintName)
				w.writeWithMultipleInsert(entries)
				// early return
				return
			}
		}
		log.Printf("Unable to insert (copyFrom): %v\n", err)

	}
	if copyCount != writeBatchSize {
		log.Printf("writeWithCopyFrom inserted %d rows\n", copyCount)
	}
}

func (w *Writer) writeWithMultipleInsert(entries []types.Entry) {

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
		log.Printf("Unable to execute multiple insert: %v\n", err)
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
