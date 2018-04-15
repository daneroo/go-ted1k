package main

import (
	"fmt"
	"log"
	"time"

	"github.com/daneroo/go-ted1k/flux"
	"github.com/daneroo/go-ted1k/ignore"
	"github.com/daneroo/go-ted1k/jsonl"
	"github.com/daneroo/go-ted1k/merge"
	"github.com/daneroo/go-ted1k/mysql"
	"github.com/daneroo/go-ted1k/progress"
	"github.com/daneroo/go-ted1k/types"
	"github.com/daneroo/timewalker"
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
	db := mysql.Setup(tableNames, myCredentials)
	defer db.Close()

	// gaps(myReader.Read())

	// ** to Ignore
	// 342k entries/s (~200M entries , SSD)
	// pipeToIgnore(fromMysql(db))
	// 300k entries/s (~200M entries , SSD)
	// pipeToIgnore(fromJsonl())

	// 190k entries/s (~200M entries , SSD)
	// pipeToJsonl(fromMysql(db))

	// 137k entries/s (~200M entries , SSD, empty destination)
	// 24k entries/s (~200M entries , SSD, full destination)
	// pipeToMysql(fromMysql(db), "watt2", db)

	// 130k entries/s (~200M entries , SSD)
	// pipeToMysql(fromJsonl(), "watt2", db)

	// 116k entries/s (~200M entries , SSD, empty or full)
	// pipeToFlux(fromMysql(db))

	// 197k entries/s (~200M entries , SSD)
	verify(db)
}

func verify(db *sqlx.DB) {
	monitor := &progress.Monitor{Batch: progress.BatchByDay}

	vv := merge.Verify(fromJsonl(), monitor.Monitor(fromMysql(db)))
	log.Printf("Verified:\n")
	for _, v := range vv {
		log.Println(v)
	}

}

func fromJsonl() <-chan types.Entry {
	jsonlReader := jsonl.DefaultReader()
	jsonlReader.Grain = timewalker.Month
	return jsonlReader.Read()
}
func fromMysql(db *sqlx.DB) <-chan types.Entry {
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
	return myReader.Read()
}

func gaps(in <-chan types.Entry) {
	monitor := &progress.Monitor{Batch: progress.BatchByDay}
	ignore.Write(monitor.Gaps(in))
}

func pipeToIgnore(in <-chan types.Entry) {
	monitor := &progress.Monitor{Batch: progress.BatchByDay}
	ignore.Write(monitor.Monitor(in))
}

func pipeToJsonl(in <-chan types.Entry) {
	jsonlWriter := jsonl.DefaultWriter()
	jsonlWriter.Grain = timewalker.Month

	monitor := &progress.Monitor{Batch: progress.BatchByDay}
	jsonlWriter.Write(monitor.Monitor(in))
}

func pipeToMysql(in <-chan types.Entry, tableName string, db *sqlx.DB) {
	// consume the channel with this sink
	myWriter := &mysql.Writer{TableName: "watt2", DB: db}
	// log.Printf("mysql.Writer: %v", myWriter)

	monitor := &progress.Monitor{Batch: progress.BatchByDay}
	myWriter.Write(monitor.Monitor(in))
	myWriter.Close()
}

func pipeToFlux(in <-chan types.Entry) {
	fluxWriter := flux.DefaultWriter()
	// log.Printf("flux.Writer: %v", fluxWriter)
	monitor := &progress.Monitor{Batch: progress.BatchByDay}
	fluxWriter.Write(monitor.Monitor(in))
}
