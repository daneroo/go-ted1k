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
	"github.com/daneroo/go-ted1k/synth"
	"github.com/daneroo/go-ted1k/types"
	"github.com/daneroo/timewalker"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	// myCredentials = "ted:secret@tcp(192.168.99.100:3306)/ted"
	// myCredentials = "ted:secret@tcp(127.0.0.1:3306)/ted"
	myCredentials    = "ted:secret@tcp(0.0.0.0:3306)/ted"
	fmtRFC3339Millis = "2006-01-02T15:04:05.000Z07:00"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format(fmtRFC3339Millis) + " - " + string(bytes))
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
	log.Printf("Starting TED1K pump\n") // TODO(daneroo): add version,buildDate

	tableNames := []string{"watt", "watt2", "watt3"}
	db := mysql.Setup(tableNames, myCredentials)
	defer db.Close()

	// gaps(fromMysql(db))
	// dirac - rate ~ 801k/s count: 31M
	// gaps(fromSynth())

	// ** to Ignore
	// 342k/s (~200M entries , SSD) - dirac rate ~ 149k/s count: 223M
	// pipeToIgnore(fromMysql(db))
	// 300k/s (~200M entries , SSD) - dirac rate ~ 193k/s count: 223M
	// pipeToIgnore(fromJsonl())
	// dirac - rate ~ 814k/s count: 31M
	// pipeToIgnore(fromSynth())

	// 190k/s (~200M entries , SSD) - dirac rate ~ 97k/s count: 223M
	// pipeToJsonl(fromMysql(db))

	// 137k/s (~200M entries , SSD, empty destination) - dirac rate ~ 34k/s count: 223M
	// 24k/s (~200M entries , SSD, full destination) - dirac rate ~ 13k/s count: 86M too slow stopped
	// pipeToMysql(fromMysql(db), "watt2", db)

	// 130k/s (~200M entries , SSD) - dirac - rate ~ 34702.6/s count: 223101124
	// pipeToMysql(fromJsonl(), "watt", db)

	// 116k/s (~200M entries , SSD, empty or full)
	// pipeToFlux(fromMysql(db))

	// 197k/s (~200M entries , SSD) - dirac rate ~ 102/s count: 223M
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

func fromSynth() <-chan types.Entry {
	// math.PI * 1e7 ~ 1 year!
	synthReader := &synth.Reader{Epoch: synth.ThisYear, TotalRows: 3.1415926e7}
	return synthReader.Read()
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
		//  About a 10M rows for ted.watt.2016-02-14-1555.sql.bz2
		// Epoch: time.Date(2015, time.October, 1, 0, 0, 0, 0, time.UTC),
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
	monitor := &progress.Monitor{Batch: progress.BatchByDay * 10}
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
	myWriter := &mysql.Writer{TableName: tableName, DB: db}
	// log.Printf("mysql.Writer: %v", myWriter)

	monitor := &progress.Monitor{Batch: progress.BatchByDay * 10}
	myWriter.Write(monitor.Monitor(in))
	myWriter.Close()
}

func pipeToFlux(in <-chan types.Entry) {
	fluxWriter := flux.DefaultWriter()
	// log.Printf("flux.Writer: %v", fluxWriter)
	monitor := &progress.Monitor{Batch: progress.BatchByDay}
	fluxWriter.Write(monitor.Monitor(in))
}
