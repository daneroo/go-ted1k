package main

import (
	"log"
	"time"
)

func writeAll(src <-chan Entry) {
	var err error
	tx, err = db.Begin()
	checkErr(err)
	insertStmt, err = tx.Prepare(insertSql)
	checkErr(err)

	count := 0
	for entry := range src {

		writeOneRow(entry.stamp, entry.watt)
		// log.Printf("Write %v, %d  (%d)\n", entry.stamp, entry.watt, count)

		count++
		if (count % 10000) == 0 {
			log.Printf("Commit checkpoint at %d records", count)
			commitAndBeginTx()
		}

	}

	// final Close
	insertStmt.Close()
	// final Tx.commit
	err = tx.Commit() // not quite right..
	checkErr(err)

}

func commitAndBeginTx() {
	insertStmt.Close()
	var err error
	tx.Commit()
	tx, err = db.Begin()
	checkErr(err)
	insertStmt, err = tx.Prepare(insertSql)
	checkErr(err)
}
func writeOneRow(stamp time.Time, watt int) {
	// log.Printf("Write %v, %d\n", stamp, watt)
	_, err := insertStmt.Exec(stamp, watt)
	checkErr(err)
	// id, _ := result.LastInsertId()
	// affected, _ := result.RowsAffected()
	// if affected > 0 {
	//  log.Printf("id:%d affected:%d", id, affected)
	// }
}
