package main

import (
	"fmt"
	"log"
	"time"

	"github.com/daneroo/go-ted1k/flux"
	"github.com/daneroo/go-ted1k/ignore"
	"github.com/daneroo/go-ted1k/jsonl"
	"github.com/daneroo/go-ted1k/mysql"
	"github.com/daneroo/go-ted1k/progress"
	. "github.com/daneroo/go-ted1k/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	// myCredentials = "ted:secret@tcp(192.168.99.100:3306)/ted"
	myCredentials = "ted:secret@tcp(0.0.0.0:3306)/ted"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format("2006-01-02T15:04:05.0000Z") + " - " + string(bytes))
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
	log.Printf("Starting TED1K pump\n") // version,buildDate

	tableNames := []string{"watt", "watt2", "watt3"}
	db := setup(tableNames)
	defer db.Close()

	// Setup the pipeline
	// create a read-only channel for source Entry(s)
	myReader := &mysql.Reader{
		TableName: "watt",
		DB:        db,
		// Epoch:     mysql.ThisYear,
		// Epoch: mysql.Recent,
		// Epoch: mysql.SixMonths,
		// Epoch: time.Date(2015, time.November, 1, 0, 0, 0, 0, time.UTC),
		// Epoch: mysql.LastYear,
		Epoch:   mysql.AllTime,
		MaxRows: mysql.AboutADay,
	}

	// Track the progress
	monitor := &progress.Monitor{Batch: progress.BatchByDay}

	// consume the channel with this sink
	myWriter := &mysql.Writer{TableName: "watt2", DB: db}
	log.Printf("mysql.Writer: %v", myWriter)

	fluxWriter := flux.DefaultWriter()
	log.Printf("flux.Writer: %v", fluxWriter)

	jsonlReader := jsonl.DefaultReader()
	// jsonlReader.Grain = timewalker.Month
	log.Printf("jsonl.Reader: %v", jsonlReader)

	jsonlWriter := jsonl.DefaultWriter()
	// jsonlWriter.Grain = timewalker.Month
	log.Printf("jsonl.Writer: %v", jsonlWriter)

	// 320k entries/s
	// ignore.Write(monitor.Monitor(myReader.Read()))
	ignore.Write(monitor.Gaps(myReader.Read()))

	// 3.5k entries/s
	// myWriter.Write(monitor.Monitor(myReader.Read()))

	// x.xk entries/s
	// fluxWriter.Write(monitor.Monitor(myReader.Read()))

	// 120k entries/s
	// jsonlWriter.Write(monitor.Monitor(myReader.Read()))

	// 230k entries/s
	// ignore.Write(monitor.Monitor(jsonlReader.Read()))

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
	ddlFormat := "CREATE TABLE IF NOT EXISTS %s ( stamp datetime NOT NULL DEFAULT '1970-01-01 00:00:00', watt int(11) NOT NULL DEFAULT '0',  PRIMARY KEY (`stamp`) )  ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;"
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
