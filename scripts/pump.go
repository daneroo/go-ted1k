package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/daneroo/go-ted1k/ephemeral"
	"github.com/daneroo/go-ted1k/flux"
	"github.com/daneroo/go-ted1k/ignore"
	"github.com/daneroo/go-ted1k/jsonl"
	"github.com/daneroo/go-ted1k/merge"
	"github.com/daneroo/go-ted1k/mysql"
	"github.com/daneroo/go-ted1k/postgres"
	"github.com/daneroo/go-ted1k/progress"
	"github.com/daneroo/go-ted1k/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	myCredentials = "ted:secret@tcp(0.0.0.0:3306)/ted"
	// pgCredentials    = "postgres://postgres:secret@127.0.0.1:5432/ted"
	pgCredentials    = "postgres://postgres:secret@0.0.0.0:5432/ted"
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

	tableNames := []string{"watt", "watt2"}
	// db := mysql.Setup(tableNames, myCredentials)
	// defer db.Close()
	conn := postgres.Setup(context.Background(), tableNames, pgCredentials)
	defer conn.Close(context.Background())

	if false {

		time.Sleep(100 * time.Millisecond)
		log.Println("-= ephemeral -> ephemeral")
		// ephemeral.NewWriter().Write(ephemeral.NewReader().Read()) // unmonitored
		ephemeral.NewWriter().Write(ephemeral.Monitor("ephemeral->ephemeral", ephemeral.NewReader().Read())) // monitored

		time.Sleep(100 * time.Millisecond)
		log.Println("-= ephemeral -> jsonl")
		// jsonl.NewWriter().Write(ephemeral.NewReader().Read()) // unmonitored
		jsonl.NewWriter().Write(ephemeral.Monitor("ephemeral->jsonl", ephemeral.NewReader().Read())) // monitored

		time.Sleep(100 * time.Millisecond)
		log.Println("-= jsonl -> ephemeral")
		// ephemeral.NewWriter().Write(jsonl.NewReader().Read()) // unmonitored
		ephemeral.NewWriter().Write(ephemeral.Monitor("jsonl -> ephemeral", jsonl.NewReader().Read())) // monitored

		time.Sleep(100 * time.Millisecond)
		log.Println("-= ephemeral -> postgres")
		// postgres.NewWriter().Write(ephemeral.NewReader().Read())                                        // unmonitored
		postgres.NewWriter(conn, tableNames[0]).Write(ephemeral.Monitor("ephemeral->postgres", ephemeral.NewReader().Read()))    // monitored
		postgres.NewWriter(conn, tableNames[0]).Write(ephemeral.Monitor("ephemeral->postgres(2)", ephemeral.NewReader().Read())) // monitored

		time.Sleep(100 * time.Millisecond)
		log.Println("-= postgres -> ephemeral")
		// ephemeral.NewWriter().Write(jsonl.NewReader().Read())                                          // unmonitored
		ephemeral.NewWriter().Write(ephemeral.Monitor("postgres -> ephemeral", postgres.NewReader(conn, tableNames[0]).Read())) // monitored
	}

	verify("ephemeral<->ephemeral", ephemeral.NewReader().Read(), ephemeral.NewReader().Read())
	verify("jsonl<->ephemeral", jsonl.NewReader().Read(), ephemeral.NewReader().Read())
	verify("postgres<->ephemeral", postgres.NewReader(conn, tableNames[0]).Read(), ephemeral.NewReader().Read())
	verify("postgres<->jsonl", postgres.NewReader(conn, tableNames[0]).Read(), jsonl.NewReader().Read())

	// dirac = rate ~ 71k/s count: 31M - empty destination - withMultipleInsert
	// dirac = rate ~ 109k/s count: 31M - full destination - withMultipleInsert
	// dirac = rate ~ 159k/s count: 31M - empty destination - writeWithCopyFrom / withMultipleInsert as fallback
	// dirac = rate ~ 100k/s count: 31M - full destination - writeWithCopyFrom / withMultipleInsert as fallback
	// dirac - size ~3.0G
	// pipeToPostgres(fromSynth(), "watt", conn)

	// e2e: verify postgres,fromSynth
	// dirac = took 2m rate ~ 240k/s count: 31M
	// {
	// 	log.Println("Verifying synth<->postgres")
	// 	monitor := &progress.Monitor{Batch: progress.BatchByDay * 10}
	// 	verify(fromSynth(), monitor.Monitor(fromPostgres(conn)))
	// }

	// gaps(fromMysql(db))
	// dirac - rate ~ 801k/s count: 31M
	// gaps(fromSynth())

	// ** to Ignore
	// 342k/s (~200M entries , SSD) - dirac rate ~ 149k/s count: 223M
	// pipeToIgnore(fromMysql(db))

	// proxmox - 5m10s - rate ~ 718k/s count: 223M
	// proxmox - 6m12s - rate ~ 598k/s count: 223M (hyper ch:0)
	// proxmox - 4m33s - rate ~ 816k/s count: 223M (hyper ch:100)
	// pipeToIgnore(fromPostgres(conn))

	// 300k/s (~200M entries , SSD) - dirac rate ~ 193k/s count: 223M
	// proxmox - 14m - rate ~ 275k/s count: 223M
	// proxmox - 7m37s - rate ~ 488k/s count: 223M chCap:100 buf:32k
	// pipeToIgnore(fromJsonl())
	// dirac - rate ~ 814k/s count: 31M
	// proxmox - 19s - rate ~ 1564k/s count: 31M
	// pipeToIgnore(fromSynth())

	//  ** to Jsonl
	// 190k/s (~200M entries , SSD) - dirac rate ~ 97k/s count: 223M
	// pipeToJsonl(fromMysql(db))

	//  ** Mysql -> Mysql
	// 137k/s (~200M entries , SSD, empty destination) - dirac rate ~ 34k/s count: 223M
	// 24k/s (~200M entries , SSD, full destination) - dirac rate ~ 13k/s count: 86M too slow stopped
	// pipeToMysql(fromMysql(db), "watt2", db)

	//  ** Jsonl -> Mysql
	// 130k/s (~200M entries , SSD) - dirac - rate ~ 34702.6/s count: 223101124
	// pipeToMysql(fromJsonl(), "watt", db)

	//  ** Jsonl -> Postgres
	// dirac - took 1h30m rate ~ 41k/s count: 223M - empty destination (hyper)
	// dirac - took 1h28mm rate ~ 42k/s count: 223M - full destination (hyper)
	// proxmox - 20m5s  - rate ~ 185k/s count: 223M - empty destination (non-hyper)
	// proxmox - 24m57s - rate ~ 148k/s count: 223M - empty destination (hyper)
	// proxmox - 30m54s - rate ~ 120k/s count: 223M - full destination (hyper)
	// pipeToPostgres(fromJsonl(), "watt", conn)

	// proxmox - 21m - rate ~ 175k/s count: 223M
	// proxmox - 23m - rate ~ 160k/s count: 223M (hyper)
	// proxmox - 10m12s - rate ~ 364k/s count: 223M (hyper - new jsonl)
	// proxmox - 9m45s - rate ~ 381k/s count: 223M (hyper - new jsonl new pgReader ch:100)
	// proxmox - 9m39s - rate ~ 385k/s count: 223M (hyper - new jsonl new pgReader ch:1000)
	// {
	// 	monitor := &progress.Monitor{Batch: progress.BatchByDay}
	// 	verify(fromJsonl(), monitor.Monitor(fromPostgres(conn)))
	// }

	//  ** Mysql -> Flux
	// 116k/s (~200M entries , SSD, empty or full)
	// pipeToFlux(fromMysql(db))

	// 197k/s (~200M entries , SSD) - dirac rate ~ 102/s count: 223M
	// {
	// 	monitor := &progress.Monitor{Batch: progress.BatchByDay}
	// 	verify(fromJsonl(), monitor.Monitor(fromMysql(db)))
	// }

}

func verify(name string, a, b <-chan []types.Entry) {
	vv := merge.Verify(a, b)
	log.Printf("Verified %s:\n", name)
	for _, v := range vv {
		log.Println(v)
	}
}

// func fromSynth() <-chan types.Entry {
// 	// math.PI * 1e7 ~ 1 year!
// 	synthReader := &synth.Reader{Epoch: synth.ThisYear, TotalRows: 3.1415926e7}
// 	return synthReader.Read()
// }

// func fromJsonl() <-chan types.Entry {
// 	jsonlReader := jsonl.NewReader()
// 	return jsonlReader.Read()
// }

// func fromPostgres(conn *pgx.Conn) <-chan types.Entry {
// 	// create a read-only channel for source Entry(s)
// 	pgReader := &postgres.Reader{
// 		TableName: "watt",
// 		Conn:      conn,
// 		// Epoch:     mysql.ThisYear,
// 		// Epoch: mysql.Recent,
// 		// Epoch: mysql.SixMonths,
// 		//  About a 10M rows for ted.watt.2016-02-14-1555.sql.bz2
// 		// Epoch: time.Date(2015, time.October, 1, 0, 0, 0, 0, time.UTC),
// 		// Epoch: mysql.LastYear,
// 		Epoch:   mysql.AllTime,
// 		MaxRows: mysql.AboutADay,
// 	}
// 	return pgReader.Read()
// }

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

// func pipeToJsonl(in <-chan types.Entry) {
// 	jsonlWriter := jsonl.NewWriter()
// 	monitor := &progress.Monitor{Batch: progress.BatchByDay}
// 	jsonlWriter.Write(monitor.Monitor(in))
// }

func pipeToMysql(in <-chan types.Entry, tableName string, db *sqlx.DB) {
	// consume the channel with this sink
	myWriter := &mysql.Writer{TableName: tableName, DB: db}
	// log.Printf("mysql.Writer: %v", myWriter)

	monitor := &progress.Monitor{Batch: progress.BatchByDay * 10}
	myWriter.Write(monitor.Monitor(in))
	myWriter.Close()
}

// func pipeToPostgres(in <-chan types.Entry, tableName string, conn *pgx.Conn) {
// 	// consume the channel with this sink
// 	pgWriter := &postgres.Writer{TableName: tableName, Conn: conn}
// 	// log.Printf("mysql.Writer: %v", myWriter)

// 	monitor := &progress.Monitor{Batch: progress.BatchByDay * 10}
// 	pgWriter.Write(monitor.Monitor(in))
// 	pgWriter.Close()
// }

func pipeToFlux(in <-chan types.Entry) {
	fluxWriter := flux.DefaultWriter()
	// log.Printf("flux.Writer: %v", fluxWriter)
	monitor := &progress.Monitor{Batch: progress.BatchByDay}
	fluxWriter.Write(monitor.Monitor(in))
}
