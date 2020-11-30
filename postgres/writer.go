package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/daneroo/go-ted1k/types"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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
	log.Printf("Closing things.. (not connection)")
}

// NewWriter is a constructor for the Writer struct
func NewWriter(conn *pgx.Conn, tableName string) *Writer {
	return &Writer{
		Conn:      conn,
		TableName: tableName,
	}
}

// Write consumes an Entry channel - returns (count,error)
// The returned count is the number of entries that were processed.
// Note: We could return the affected rows by summing the returned values from w.flush()
// Note: We could terminate early in case of error (from w.flush), which we currently silently ignore
func (w *Writer) Write(src <-chan []types.Entry) (int, error) {
	count := 0
	entries := make([]types.Entry, 0, writeBatchSize)

	for slice := range src {
		for _, entry := range slice {

			entries = append(entries, entry)

			count++
			if (len(entries) % writeBatchSize) == 0 {
				w.flush(entries)
				entries = make([]types.Entry, 0, writeBatchSize)
			}

		}
	}
	// last flush
	w.flush(entries)
	return count, nil
}

// flush performs the actual batch insert.
func (w *Writer) flush(entries []types.Entry) (int, error) {
	if len(entries) == 0 {
		return 0, nil
	}
	return w.writeWithFallback(entries)
}

// writeWithFallback inserts rows with the fastest method available
// It first attempts to use writeWithCopyFrom,
// but it falls back to writeWithMultipleInsert when we can assert the specific `duplicate key value` error
func (w *Writer) writeWithFallback(entries []types.Entry) (int, error) {
	copyCount, err := w.writeWithCopyFrom(entries)
	// fallback to writeWithMultipleInsert if we have the specific `duplicate key value` error
	if err != nil && isUniqueKeyError(err) {
		// log.Printf("Retrying withMultipleInsert: err: %v\n", err)
		return w.writeWithMultipleInsert(entries)
	}
	return copyCount, err
}

// isUniqueKeyError detects the specific error (unique key violation)
func isUniqueKeyError(err error) bool {
	if err != nil {
		// Detect pgconn.PgError
		if pge, ok := err.(*pgconn.PgError); ok {
			// We know we have a pgconn.PgError
			// pgconn.PgError: ERROR: duplicate key value violates unique constraint "watt_pkey" (SQLSTATE 23505)
			if pge.Code == "23505" {
				// log.Printf("Retrying withMultipleInsert: Code: %v TableName: %v ConstraintName: %v\n", pge.Code, pge.TableName, pge.ConstraintName)
				return true
			}
		}
		log.Printf("Unable to insert (copyFrom): %v\n", err)
	}
	return false
}

// writeWithCopyFrom is the fastest way to insert, but has no "ON CONFLICT (stamp) DO NOTHING" mechanism,
func (w *Writer) writeWithCopyFrom(entries []types.Entry) (int, error) {
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
	return int(copyCount), err
}

// Writes to postgres using a multiple value insert,
// using `ON CONFLICT DO NOTHING`
// The returned count is the affected rows, not necessarily the number of entries passed in.
func (w *Writer) writeWithMultipleInsert(entries []types.Entry) (int, error) {

	// log.Printf("flush: would have flushed %d entries", len(entries))
	sql := w.makeSQL(len(entries))

	// flatten the parameters into a single value array
	vals := []interface{}{}
	for _, entry := range entries {
		vals = append(vals, entry.Stamp, entry.Watt)
	}

	// commandTag, err := w.Conn.Exec(context.Background(), sql, vals...)
	commandTag, err := w.Conn.Exec(context.Background(), sql, vals...)
	if err != nil {
		log.Printf("Unable to execute multiple insert: %v\n", err)
		return 0, err
	}
	return int(commandTag.RowsAffected()), nil
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
