package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/daneroo/go-ted1k/ephemeral"
	"github.com/daneroo/go-ted1k/jsonl"
	"github.com/daneroo/go-ted1k/logsetup"
	"github.com/daneroo/go-ted1k/merge"
	"github.com/daneroo/go-ted1k/postgres"
	"github.com/daneroo/go-ted1k/progress"
	"github.com/daneroo/go-ted1k/types"
	_ "github.com/go-sql-driver/mysql"
)

const (
	// myCredentials = "ted:secret@tcp(0.0.0.0:3306)/ted"
	// myCredentials = "root@tcp(darwin.imetrical.com:3306)/ted"

	pgCredentialsDefault = "postgres://postgres:secret@0.0.0.0:5432/ted"
	// pgCredentialsDefault = "postgres://postgres:secret@d1-px1:5432/ted"
)

func main() {
	logsetup.SetupFormat()
	log.Printf("Starting TED1K dump jsonl\n") // TODO(daneroo): add version,buildDate

	// db := mysql.Setup([]string{}, myCredentials)
	// defer db.Close()

	pgCredentials := os.Getenv("PGCONN")
	if pgCredentials == "" {
		pgCredentials = pgCredentialsDefault
	}
	conn := postgres.Setup(context.Background(), []string{"watt"}, pgCredentials)
	defer conn.Close(context.Background())

	// phase 1: jsonl -> postgres
	if true {
		verify("jsonl <-> postgres", jsonl.NewReader(), postgres.NewReader(conn, "watt"))
		skipCopyFrom := true
		doTest(fmt.Sprintf("jsonl -> postgres (skipCopyFrom=%v)", skipCopyFrom), jsonl.NewReader(), postgres.NewWriter(conn, "watt", skipCopyFrom))
	}

}

func doTest(name string, r types.EntryReader, w types.EntryWriter) (int, error) {
	log.Printf("-=- %s\n", name)
	return w.Write(progress.Monitor(name, r.Read()))
}

func verify(name string, a, b types.EntryReader) {
	log.Printf("-=- %s\n", name)
	// vv := merge.Verify(a.Read(), progress.Monitor(name+" (B)", b.Read()))
	vv := merge.Verify(progress.Monitor(name+" (A)", a.Read()), progress.Monitor(name+" (B)", b.Read()))
	log.Printf("Verified %s:\n", name)
	for _, v := range vv {
		log.Println(v)
	}
}

func gaps(in <-chan []types.Entry) {
	ephemeral.NewWriter().Write(progress.Gaps(in))
}
