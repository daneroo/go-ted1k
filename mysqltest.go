package main

import (
	// "github.com/daneroo/go-mysqltest/flux"
	"log"

	"github.com/daneroo/go-mysqltest/ignore"
	"github.com/daneroo/go-mysqltest/mysql"
	"github.com/daneroo/go-mysqltest/progress"
	. "github.com/daneroo/go-mysqltest/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	myCredentials = "ted:secret@tcp(192.168.99.100:3306)/ted"
)

func main() {

	db := setup()
	defer db.Close()

	// Setup the pipeline
	// create a read-only channel for source Entry(s)
	myReader := &mysql.Reader{
		TableName: "watt",
		DB:        db,
		Epoch:     mysql.Recent,
		// Epoch:   mysql.SixMonths,
		MaxRows: mysql.AboutADay,
	}
	log.Printf("mysql.Reader: %v", myReader)

	// Track the progress
	monitor := &progress.Monitor{
		Batch: progress.BatchByDay,
	}
	// consume the channel with this sink
	myWriter := &mysql.Writer{
		TableName: "watt2",
		DB:        db,
	}
	log.Printf("mysql.Writer: %v", myWriter)

	ignore.Write(monitor.Monitor(myReader.Read()))
	// myWriter.Write(monitor.Monitor(myReader.Read()))

	// consume the channel with this sink
	// flux.WriteAll(src)
}

func setup() *sqlx.DB {
	// Connect is Open and verify with a Ping
	db, err := sqlx.Connect("mysql", myCredentials)
	Checkerr(err)
	log.Println("Connected to MySQL")

	createCopyTable(db)
	totalCount(db)

	return db
}

func createCopyTable(db *sqlx.DB) {
	// ddl:="create table if not exists watt2 like watt"
	ddl := "CREATE TABLE IF NOT EXISTS watt2 ( stamp datetime NOT NULL DEFAULT '1970-01-01 00:00:00', watt int(11) NOT NULL DEFAULT '0',  PRIMARY KEY (`stamp`) )"
	_, err := db.Exec(ddl)
	Checkerr(err)
}

func totalCount(db *sqlx.DB) {
	var totalCount int
	err := db.Get(&totalCount, "SELECT COUNT(*) FROM watt")
	Checkerr(err)
	log.Printf("Found %d entries in watt\n", totalCount)
}
