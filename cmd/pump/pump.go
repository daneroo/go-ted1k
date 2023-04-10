package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/daneroo/go-ted1k/ephemeral"
	"github.com/daneroo/go-ted1k/ipfs"
	"github.com/daneroo/go-ted1k/jsonl"
	"github.com/daneroo/go-ted1k/logsetup"
	"github.com/daneroo/go-ted1k/merge"
	"github.com/daneroo/go-ted1k/mysql"
	"github.com/daneroo/go-ted1k/postgres"
	"github.com/daneroo/go-ted1k/progress"
	"github.com/daneroo/go-ted1k/types"
	_ "github.com/go-sql-driver/mysql"
	shell "github.com/ipfs/go-ipfs-api"
)

const (
	// myCredentials = "ted:secret@tcp(0.0.0.0:3306)/ted"
	myCredentials = "root@tcp(darwin.imetrical.com:3306)/ted"

	// pgCredentials    = "postgres://postgres:secret@127.0.0.1:5432/ted"
	pgCredentialsDefault = "postgres://postgres:secret@0.0.0.0:5432/ted"
)

func main() {
	var since time.Duration
	defaultDuration := 100 * 24 * time.Hour // 100 days
	flag.DurationVar(&since, "since", defaultDuration, "duration since the event (e.g., 2400h)")
	var skipCopyFrom bool
	flag.BoolVar(&skipCopyFrom, "skip-copy-from", false, "Always use writeWithMultipleInserts, do not use writeWithCopyFrom ")
	flag.Parse()

	logsetup.SetupFormat()

	// just for prettier output: I wish flag.DurationVar() would handle days directly
	sinceDays := since.Hours() / 24
	log.Printf("Starting TED1K pump (since=%.1f days skip-copy-from=%v)\n", sinceDays, skipCopyFrom) // TODO(daneroo): add version,buildDate

	// tableNames := []string{"watt", "watt2"}
	tableNames := []string{"watt"}
	db := mysql.Setup(tableNames, myCredentials)
	defer db.Close()

	pgCredentials := os.Getenv("PGCONN")
	if pgCredentials == "" {
		pgCredentials = pgCredentialsDefault
	}

	conn := postgres.Setup(context.Background(), tableNames, pgCredentials)
	defer conn.Close(context.Background())
	sh := shell.NewShell("localhost:5001")
	if sh == nil {
		// just to use the sh variable
		log.Fatal("IPFS shell is nil")
	}

	// // ipfs -> postgres
	// if true {
	// 	dirCid := "QmSLJPEZocdPZ99pazEkiJTaf3B1zeBmAQWEr7n9fSNgEu"
	// 	doTest("ipfs <-> postgres", ipfs.NewReader(sh, dirCid), postgres.NewWriter(conn, "watt"))
	// }
	// os.Exit(0)

	// mysql (remote) -> postgres
	if true {
		myReader := mysql.NewReader(db, "watt")
		myReader.Epoch = time.Now().Add(-since)
		log.Printf("Reading MySQL since %s\n", myReader.Epoch)
		pgReader := postgres.NewReader(conn, "watt")
		pgReader.Epoch = myReader.Epoch
		verify("mysql <-> postgres", myReader, pgReader)

		// actually insert
		// Move this to subscribeAndSync
		doTest(fmt.Sprintf("mysql -> postgres (skipCopyFrom=%v)", skipCopyFrom), myReader, postgres.NewWriter(conn, "watt", skipCopyFrom))
	}
	os.Exit(0)

	// postgres to ipfs
	// fmt.Println()
	// iw := ipfs.NewWriter(sh)
	// doTest("postgres -> ipfs", postgres.NewReader(conn, "watt"), iw)
	// dirCid := iw.Dw.Dir
	// log.Printf("CID: %s\n", dirCid)
	// verify("postgres <-> ipfs", postgres.NewReader(conn, "watt"), ipfs.NewReader(sh, dirCid))

	// ephemeral
	if true {
		fmt.Println()
		doTest("ephemeral -> ephemeral", ephemeral.NewReader(), ephemeral.NewWriter())
		verify("ephemeral <-> ephemeral", ephemeral.NewReader(), ephemeral.NewReader())
	}

	// jsonl
	if true {
		fmt.Println()
		doTest("ephemeral -> jsonl", ephemeral.NewReader(), jsonl.NewWriter())
		doTest("jsonl -> ephemeral", jsonl.NewReader(), ephemeral.NewWriter())
		verify("ephemeral <-> jsonl", ephemeral.NewReader(), jsonl.NewReader())
	}

	// ipfs
	if true {
		fmt.Println()
		iw := ipfs.NewWriter(sh)
		doTest("ephemeral -> ipfs", ephemeral.NewReader(), iw)
		dirCid := iw.Dw.Dir
		// dirCid := "QmYEZzGXRwzWArokCyEqpJnLrbp3F2WEUY46huWtu6TqL6"
		doTest("ipfs -> ephemeral", ipfs.NewReader(sh, dirCid), ephemeral.NewWriter())
		verify("ephemeral <-> ipfs", ephemeral.NewReader(), ipfs.NewReader(sh, dirCid))
	}

	// postgres
	if true {
		fmt.Println()
		doTest("ephemeral -> postgres", ephemeral.NewReader(), postgres.NewWriter(conn, tableNames[0], skipCopyFrom))
		doTest("postgres -> ephemeral", postgres.NewReader(conn, tableNames[0]), ephemeral.NewWriter())
		verify("ephemeral <-> postgres", ephemeral.NewReader(), postgres.NewReader(conn, tableNames[0]))
	}

	// mysql
	if true {
		fmt.Println()
		doTest("ephemeral -> mysql", ephemeral.NewReader(), mysql.NewWriter(db, tableNames[0]))
		doTest("mysql -> ephemeral", mysql.NewReader(db, tableNames[0]), ephemeral.NewWriter())
		verify("ephemeral <-> mysql", ephemeral.NewReader(), mysql.NewReader(db, tableNames[0]))
	}

	//  ** Mysql -> Flux
	// 116k/s (~200M entries , SSD, empty or full)
	// pipeToFlux(fromMysql(db))

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
