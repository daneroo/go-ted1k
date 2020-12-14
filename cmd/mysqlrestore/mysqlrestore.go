package main

import (
	"context"
	"log"

	"github.com/daneroo/go-ted1k/ephemeral"
	"github.com/daneroo/go-ted1k/jsonl"
	"github.com/daneroo/go-ted1k/logsetup"
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
	pgCredentials = "postgres://postgres:secret@0.0.0.0:5432/ted"
)

func main() {
	logsetup.SetupFormat()
	log.Printf("Starting TED1K mysql restore\n") // TODO(daneroo): add version,buildDate

	db := mysql.Setup([]string{}, myCredentials)
	defer db.Close()

	conn := postgres.Setup(context.Background(), []string{"watt"}, pgCredentials)
	defer conn.Close(context.Background())

	// phase-2 (restore mysql/watt to postgres)
	if true {
		// pre-phase-2:json->postgres (jsonl-ted-rollup.20150928.1006)
		// doTest("jsonl -> postgres", jsonl.NewReader(), postgres.NewWriter(conn, "watt"))
		// verify("jsonl <-> postgres", jsonl.NewReader(), postgres.NewReader(conn, "watt"))

		// Phase-2
		// verify and copy to postgres
		// verify("mysql(watt) <-> postgres", mysql.NewReader(db, "watt"), postgres.NewReader(conn, "watt"))
		// doTest("mysql(watt) -> postgres", mysql.NewReader(db, "watt"), postgres.NewWriter(conn, "watt"))

		// verify only
		verify("mysql(watt) <-> postgres", mysql.NewReader(db, "watt"), postgres.NewReader(conn, "watt"))

	}

	// phase-1 (restore compare to json to watt and ted_native)
	if false {
		// pre phase-1
		// merge ted_native over watt (from ted.20150928.1006.sql.bz2)
		// doTest("mysql(ted_native) -> mysql(watt)", mysql.NewReader(db, "ted_native"), mysql.NewWriter(db, "watt"))
		// now dump the combined jsonl
		// doTest("mysql(watt) -> jsonl", mysql.NewReader(db, "watt"), jsonl.NewWriter())
		verify("jsonl <-> mysql(watt)", jsonl.NewReader(), mysql.NewReader(db, "watt"))
		verify("jsonl <-> mysql(ted_native)", jsonl.NewReader(), mysql.NewReader(db, "ted_native"))
		// Ignore ted_service
	}

}

func doTest(name string, r types.EntryReader, w types.EntryWriter) (int, error) {
	log.Printf("-=- %s\n", name)
	return w.Write(progress.Monitor(name, r.Read()))
}

func verify(name string, a, b types.EntryReader) {
	log.Printf("-=- %s\n", name)
	vv := merge.Verify(a.Read(), progress.Monitor(name, b.Read()))
	log.Printf("Verified %s:\n", name)
	for _, v := range vv {
		log.Println(v)
	}
}

func gaps(in <-chan []types.Entry) {
	ephemeral.NewWriter().Write(progress.Gaps(in))
}
