package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/daneroo/go-ted1k/ephemeral"
	"github.com/daneroo/go-ted1k/jsonl"
	"github.com/daneroo/go-ted1k/merge"
	"github.com/daneroo/go-ted1k/mysql"
	"github.com/daneroo/go-ted1k/postgres"
	"github.com/daneroo/go-ted1k/progress"
	"github.com/daneroo/go-ted1k/types"
	_ "github.com/go-sql-driver/mysql"
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
	db := mysql.Setup(tableNames, myCredentials)
	defer db.Close()
	conn := postgres.Setup(context.Background(), tableNames, pgCredentials)
	defer conn.Close(context.Background())

	if false {
		time.Sleep(100 * time.Millisecond)
		log.Println("-= ephemeral -> ephemeral")
		ephemeral.NewWriter().Write(progress.Monitor("ephemeral->ephemeral", ephemeral.NewReader().Read()))

		time.Sleep(100 * time.Millisecond)
		log.Println("-= ephemeral -> jsonl")
		jsonl.NewWriter().Write(progress.Monitor("ephemeral->jsonl", ephemeral.NewReader().Read()))

		time.Sleep(100 * time.Millisecond)
		log.Println("-= jsonl -> ephemeral")
		ephemeral.NewWriter().Write(progress.Monitor("jsonl -> ephemeral", jsonl.NewReader().Read()))

		time.Sleep(100 * time.Millisecond)
		log.Println("-= ephemeral -> postgres")
		postgres.NewWriter(conn, tableNames[0]).Write(progress.Monitor("ephemeral->postgres", ephemeral.NewReader().Read()))
		// postgres.NewWriter(conn, tableNames[0]).Write(progress.Monitor("ephemeral->postgres(2)", ephemeral.NewReader().Read()))

		time.Sleep(100 * time.Millisecond)
		log.Println("-= postgres -> ephemeral")
		ephemeral.NewWriter().Write(progress.Monitor("postgres -> ephemeral", postgres.NewReader(conn, tableNames[0]).Read()))

		time.Sleep(100 * time.Millisecond)
		log.Println("-= jsonl -> postgres")
		postgres.NewWriter(conn, tableNames[0]).Write(progress.Monitor("jsonl->postgres", jsonl.NewReader().Read()))

		verify("ephemeral<->ephemeral", ephemeral.NewReader().Read(), ephemeral.NewReader().Read())
		verify("jsonl<->ephemeral", jsonl.NewReader().Read(), ephemeral.NewReader().Read())
		verify("postgres<->ephemeral", postgres.NewReader(conn, tableNames[0]).Read(), ephemeral.NewReader().Read())
		verify("postgres<->jsonl", postgres.NewReader(conn, tableNames[0]).Read(), jsonl.NewReader().Read())
	}

	if false {

		time.Sleep(100 * time.Millisecond)
		log.Println("-= ephemeral -> mysql")
		mysql.NewWriter(db, tableNames[0]).Write(progress.Monitor("ephemeral->mysql", ephemeral.NewReader().Read()))
		mysql.NewWriter(db, tableNames[0]).Write(progress.Monitor("ephemeral->mysql(2)", ephemeral.NewReader().Read()))

		time.Sleep(100 * time.Millisecond)
		log.Println("-= jsonl -> mysql")
		mysql.NewWriter(db, tableNames[0]).Write(progress.Monitor("jsonl->mysql", jsonl.NewReader().Read()))

		time.Sleep(100 * time.Millisecond)
		log.Println("-= mysql -> ephemeral")
		ephemeral.NewWriter().Write(progress.Monitor("mysql -> ephemeral", mysql.NewReader(db, tableNames[0]).Read()))

		verify("mysql<->ephemeral", mysql.NewReader(db, tableNames[0]).Read(), ephemeral.NewReader().Read())
	}

	// verify("mysql<->ephemeral", mysql.NewReader(db, tableNames[0]).Read(), ephemeral.NewReader().Read())
	// gaps(ephemeral.NewReader().Read())
	// gaps(mysql.NewReader(db, tableNames[0]).Read())

	// log.Println("-= ephemeral -> mysql -> mysql.watt2")
	// mysql.NewWriter(db, tableNames[0]).Write(progress.Monitor("ephemeral->mysql", ephemeral.NewReader().Read()))
	// mysql.NewWriter(db, tableNames[1]).Write(progress.Monitor("ephemeral->mysql2", ephemeral.NewReader().Read()))
	mysql.NewWriter(db, tableNames[1]).Write(progress.Monitor("mysql->mysql2", mysql.NewReader(db, tableNames[0]).Read()))
	verify("mysql<->mysql", mysql.NewReader(db, tableNames[0]).Read(), mysql.NewReader(db, tableNames[1]).Read())

	//  ** Mysql -> Mysql
	// 137k/s (~200M entries , SSD, empty destination) - dirac rate ~ 34k/s count: 223M
	// 24k/s (~200M entries , SSD, full destination) - dirac rate ~ 13k/s count: 86M too slow stopped
	// pipeToMysql(fromMysql(db), "watt2", db)

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

func gaps(in <-chan []types.Entry) {
	ephemeral.NewWriter().Write(progress.Gaps(in))
}
