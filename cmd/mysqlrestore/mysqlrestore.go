package main

import (
	"log"

	"github.com/daneroo/go-ted1k/ephemeral"
	"github.com/daneroo/go-ted1k/jsonl"
	"github.com/daneroo/go-ted1k/logsetup"
	"github.com/daneroo/go-ted1k/merge"
	"github.com/daneroo/go-ted1k/mysql"
	"github.com/daneroo/go-ted1k/progress"
	"github.com/daneroo/go-ted1k/types"
	_ "github.com/go-sql-driver/mysql"
)

const (
	myCredentials = "ted:secret@tcp(0.0.0.0:3306)/ted"
)

func main() {
	logsetup.SetupFormat()
	log.Printf("Starting TED1K mysql restore\n") // TODO(daneroo): add version,buildDate

	db := mysql.Setup([]string{}, myCredentials)
	defer db.Close()

	// mysql
	if true {
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
