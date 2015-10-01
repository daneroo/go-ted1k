package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Entry struct {
	stamp time.Time
	watt  int
}

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
	log.Println("Prepared insert statement (in a transaction)")

	var totalCount int
	row := db.QueryRow("SELECT COUNT(*) FROM watt")
	err = row.Scan(&totalCount)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		panic(err)
	}
	log.Printf("Found %d entries in watt\n", totalCount)

	// create a read-only channel for source Entry(s)
	src := readAll()
	// consume the channel with this sink
	writeAll(src)
}

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
