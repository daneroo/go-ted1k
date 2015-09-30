package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
)

const (
	// myCredentials = "daniel@tcp(192.168.5.105:3306)/ted"
	myCredentials    = "ted:secret@tcp(192.168.99.100:3306)/ted"
	insertSql        = "INSERT IGNORE INTO watt2 (stamp, watt) VALUES (?,?)"
	maxCountPerChunk = 3600 * 24
)

var (
	db         *sql.DB
	tx         *sql.Tx
	insertStmt *sql.Stmt
	epoch      = time.Date(2015, time.September, 27, 0, 0, 0, 0, time.UTC)
	// epoch         = time.Date(2007, time.January, 0, 0, 0, 0, 0, time.UTC)
	// epoch = time.Date(2037, time.January, 0, 0, 0, 0, 0, time.UTC)
)

func main() {
	log.Printf("Just getting %s\n", "started")

	var err error
	db, err = sql.Open("mysql", myCredentials)
	checkErr(err)
	defer db.Close()
	log.Println("Survived Opening")

	createCopyTable()

	// insertStmt, err = db.Prepare(insertSql)
	// defer insertStmt.Close()
	log.Println("Prepared insert statement (in a transaction")

	var totalCount int
	row := db.QueryRow("SELECT COUNT(*) FROM watt")
	err = row.Scan(&totalCount)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		panic(err)
	}
	log.Printf("Found %d entries in watt\n", totalCount)

	tx, err = db.Begin()
	checkErr(err)
	insertStmt, err = tx.Prepare(insertSql)
	checkErr(err)

	rowCount := 0
	startTimeExcl := epoch
	for {
		chunkRowCount := 0
		startTimeExcl, chunkRowCount = oneChunk(db, startTimeExcl, maxCountPerChunk)
		rowCount += chunkRowCount

		commitAndBeginTx()

		if chunkRowCount == 0 {
			break
		}
	}

	// final Close
	insertStmt.Close()
	// final Tx.commit
	err = tx.Commit() // not quite right..
	checkErr(err)

	log.Printf("Fetched a total of %d rows (%d before iteration)", rowCount, totalCount)

}

func oneChunk(db *sql.DB, startTimeExcl time.Time, maxCountPerChunk int) (time.Time, int) {
	defer timeTrack(time.Now(), "oneChunk", maxCountPerChunk)
	sql := "SELECT stamp,watt FROM watt where stamp>? ORDER BY stamp ASC LIMIT ?"
	// sql := "SELECT stamp,watt FROM watt where stamp<? ORDER BY stamp DESC LIMIT ?"
	rows, err := db.Query(sql, startTimeExcl, maxCountPerChunk)
	checkErr(err)
	defer rows.Close()

	avgWatt := 0
	chunkRowCount := 0
	var lastStamp time.Time
	for rows.Next() {
		// var stamp string
		var stamp mysql.NullTime
		var watt int
		err = rows.Scan(&stamp, &watt)
		if err != nil {
			log.Println(err)
		}
		avgWatt += watt
		chunkRowCount++
		if stamp.Valid {
			lastStamp = stamp.Time
			writeOneRow(stamp.Time, watt)
		}
		// log.Printf(" %v: %v", stamp, watt)
	}
	err = rows.Err() // get any error encountered during iteration
	checkErr(err)

	if chunkRowCount != 0 {
		avgWatt /= chunkRowCount
	}
	log.Printf("average between (%v - %v]: %v (%v)", startTimeExcl, lastStamp, avgWatt, chunkRowCount)
	return lastStamp, chunkRowCount
}

func createCopyTable() {
	// ddl:="create table if not exists watt2 like watt"
	ddl := "CREATE TABLE IF NOT EXISTS watt2 ( stamp datetime NOT NULL DEFAULT '1970-01-01 00:00:00', watt int(11) NOT NULL DEFAULT '0',  PRIMARY KEY (`stamp`) )"
	_, err := db.Exec(ddl)
	checkErr(err)
	// log.Printf("%v\n", result)
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
	// 	log.Printf("id:%d affected:%d", id, affected)
	// }
}

func timeTrack(start time.Time, name string, count int) {
	elapsed := time.Since(start)
	if count > 0 {
		rate := float64(count) / elapsed.Seconds()
		log.Printf("%s took %s, rate ~ %.1f/s", name, elapsed, rate)
	} else {
		log.Printf("%s took %s", name, elapsed)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
