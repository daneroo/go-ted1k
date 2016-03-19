package main

import (
	// "github.com/daneroo/go-mysqltest/flux"
	"fmt"
	"log"

	"github.com/daneroo/go-mysqltest/flux"
	"github.com/daneroo/go-mysqltest/jsonl"
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

	tableNames := []string{"watt", "watt2", "watt3"}
	db := setup(tableNames)
	defer db.Close()

	// Setup the pipeline
	// create a read-only channel for source Entry(s)
	myReader := &mysql.Reader{
		TableName: "watt",
		DB:        db,
		// Epoch:     mysql.Recent,
		// Epoch:   mysql.SixMonths,
		// Epoch:     mysql.LastYear,
		Epoch:   mysql.AllTime,
		MaxRows: mysql.AboutADay,
	}
	log.Printf("mysql.Reader: %v", myReader)

	// Track the progress
	monitor := &progress.Monitor{
		Batch: progress.BatchByDay,
	}
	log.Printf("progress.Monitor: %v", monitor)

	// consume the channel with this sink
	myWriter := &mysql.Writer{
		TableName: "watt2",
		DB:        db,
	}
	log.Printf("mysql.Writer: %v", myWriter)

	fluxWriter := flux.DefaultWriter()
	log.Printf("flux.Writer: %v", fluxWriter)

	jsonlWriter := jsonl.DefaultWriter()
	log.Printf("jsonl.Writer: %v", jsonlWriter)

	// ignore.Write(monitor.Monitor(myReader.Read()))
	// myWriter.Write(monitor.Monitor(myReader.Read()))
	// fluxWriter.Write(monitor.Monitor(myReader.Read()))
	jsonlWriter.Write(monitor.Monitor(myReader.Read()))
	// jsonlWriter.Write(myReader.Read())

	// consume the channel with this sink
	// flux.WriteAll(src)
}

func setup(tableNames []string) *sqlx.DB {
	// Connect is Open and verify with a Ping
	db, err := sqlx.Connect("mysql", myCredentials)
	Checkerr(err)
	log.Println("Connected to MySQL")

	for _, tableName := range tableNames {
		createCopyTable(db, tableName)
	}
	totalCount(db)

	return db
}

func createCopyTable(db *sqlx.DB, tableName string) {
	ddlFormat := "CREATE TABLE IF NOT EXISTS %s ( stamp datetime NOT NULL DEFAULT '1970-01-01 00:00:00', watt int(11) NOT NULL DEFAULT '0',  PRIMARY KEY (`stamp`) )"
	ddl := fmt.Sprintf(ddlFormat, tableName)
	_, err := db.Exec(ddl)
	Checkerr(err)
}

func totalCount(db *sqlx.DB) {
	var totalCount int
	err := db.Get(&totalCount, "SELECT COUNT(*) FROM watt")
	Checkerr(err)
	log.Printf("Found %d entries in watt\n", totalCount)
}
