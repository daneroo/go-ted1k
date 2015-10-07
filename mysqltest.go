package main

import (
	"database/sql"
	"github.com/daneroo/go-mysqltest/flux"
	// "github.com/daneroo/go-mysqltest/sink"
	"github.com/daneroo/go-mysqltest/source"
	. "github.com/daneroo/go-mysqltest/util"
	_ "github.com/go-sql-driver/mysql"
	// "os"
	// "github.com/jmoiron/sqlx"
	"log"
)

const (
	// myCredentials = "daniel@tcp(192.168.5.105:3306)/ted"
	myCredentials = "ted:secret@tcp(192.168.99.100:3306)/ted"
)

func main() {

	// flux.Try()
	// os.Exit(0)

	db := setup()
	defer db.Close()

	// create a read-only channel for source Entry(s)
	src := source.ReadAll(db)

	// consume the channel with this sink
	// sink.IgnoreAll(db, src)
	// sink.WriteAll(db, src)
	// flux.IgnoreAll(src)
	flux.WriteAll(src)
}

func setup() *sql.DB {
	db, err := sql.Open("mysql", myCredentials)
	Checkerr(err)
	log.Println("Connected to MySQL")

	createCopyTable(db)

	var totalCount int
	row := db.QueryRow("SELECT COUNT(*) FROM watt")
	err = row.Scan(&totalCount)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		panic(err)
	}
	log.Printf("Found %d entries in watt\n", totalCount)

	return db
}

func createCopyTable(db *sql.DB) {
	// ddl:="create table if not exists watt2 like watt"
	ddl := "CREATE TABLE IF NOT EXISTS watt2 ( stamp datetime NOT NULL DEFAULT '1970-01-01 00:00:00', watt int(11) NOT NULL DEFAULT '0',  PRIMARY KEY (`stamp`) )"
	_, err := db.Exec(ddl)
	Checkerr(err)
	// log.Printf("%v\n", result)
}
