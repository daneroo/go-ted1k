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

func createCopyTable() {
	// ddl:="create table if not exists watt2 like watt"
	ddl := "CREATE TABLE IF NOT EXISTS watt2 ( stamp datetime NOT NULL DEFAULT '1970-01-01 00:00:00', watt int(11) NOT NULL DEFAULT '0',  PRIMARY KEY (`stamp`) )"
	_, err := db.Exec(ddl)
	checkErr(err)
	// log.Printf("%v\n", result)
}
