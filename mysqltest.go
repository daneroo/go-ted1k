package main

import (
	"database/sql"
	"github.com/daneroo/go-mysqltest/source"
	. "github.com/daneroo/go-mysqltest/util"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	// myCredentials = "daniel@tcp(192.168.5.105:3306)/ted"
	myCredentials = "ted:secret@tcp(192.168.99.100:3306)/ted"
	insertSql     = "INSERT IGNORE INTO watt2 (stamp, watt) VALUES (?,?)"
)

var (
	db         *sql.DB
	tx         *sql.Tx
	insertStmt *sql.Stmt
)

func main() {
	log.Printf("Just getting %s\n", "started")

	var err error
	db, err = sql.Open("mysql", myCredentials)
	Checkerr(err)
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
	src := source.ReadAll(db)
	// consume the channel with this sink
	writeAll(src)
}

func createCopyTable() {
	// ddl:="create table if not exists watt2 like watt"
	ddl := "CREATE TABLE IF NOT EXISTS watt2 ( stamp datetime NOT NULL DEFAULT '1970-01-01 00:00:00', watt int(11) NOT NULL DEFAULT '0',  PRIMARY KEY (`stamp`) )"
	_, err := db.Exec(ddl)
	Checkerr(err)
	// log.Printf("%v\n", result)
}
