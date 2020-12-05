package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/daneroo/go-ted1k/ipfs"
	"github.com/daneroo/go-ted1k/timer"
	"github.com/daneroo/go-ted1k/types"
	shell "github.com/ipfs/go-ipfs-api"
)

const (
	fmtRFC3339Millis = "2006-01-02T15:04:05.000Z07:00"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format(fmtRFC3339Millis) + " - " + string(bytes))
}

// Clean everything:
// ipfs pin ls --type recursive | cut -d' ' -f1 | xargs -n1 ipfs pin rm
// ipfs repo gc

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
	log.Printf("Starting IPFS bootstrap test\n") // TODO(daneroo): add version,buildDate

	// Where your local node is running on localhost:5001
	sh := shell.NewShell("localhost:5001")

	if true {
		cid, err := addOneFile(sh, 1, true)
		if err != nil {
			log.Fatalf("error adding file: %s", err)
		}
		log.Printf("added file: %s\n", cid)
		r, err := getOneFile(sh, cid)
		if err != nil {
			log.Fatalf("error getting file: %s", err)
		}
		defer r.Close()
		// showReader(r, cid)
	}

	if false {
		dircid, err := addOneDir(sh)
		if err != nil {
			log.Fatalf("error adding directory: %s", err)
		}
		log.Printf("added dir: %s\n", dircid)

		getDirectory(sh, dircid)

	}
}

func getDirectory(sh *shell.Shell, cid string) (*shell.UnixLsObject, error) {
	// objects := make(map[string]*shell.UnixLsObject)
	objects, err := sh.FileList(cid)
	if err != nil {
		return objects, nil
	}
	log.Printf("Dir: %+v\n", objects)
	for idx, link := range objects.Links {
		log.Printf("Link %d: %+v\n", idx, link)
		if link.Type == "File" { // File or Directory
			r, err := getOneFile(sh, link.Hash)
			if err != nil {
				log.Fatalf("error getting file: %s", err)
			}
			defer r.Close()
			// showReader(r, link.Hash)
		} else {
			getDirectory(sh, link.Hash)
		}
	}
	return objects, nil

}

func showReader(r io.ReadCloser, name string) {
	fmt.Printf("--- reading: %s -----\n", name)
	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatalf("error reading: %s", err)
	}
	fmt.Printf("\n--- done reading: %s\n", name)

}
func getOneFile(sh *shell.Shell, cid string) (io.ReadCloser, error) {
	return sh.Cat(cid)
}

func addOneFile(sh *shell.Shell, day int, pin bool) (string, error) {
	fw := ipfs.NewFWriter(sh, pin)
	var b bytes.Buffer

	stamp := time.Date(2020, time.January, day, 0, 0, 0, 0, time.UTC)
	start := time.Now()
	length := 2678400
	// size := 86400
	for i := 0; i < length; i++ {
		entry := types.Entry{Stamp: stamp, Watt: int(stamp.Unix() % 5000)}
		// fmt.Fprintf(fw.W, "{\"stamp\":\"%s\",\"watt\":%d}\n", entry.Stamp.Format(time.RFC3339Nano), entry.Watt)

		s := fmt.Sprintf("{\"stamp\":\"%s\",\"watt\":%d}\n", entry.Stamp.Format(time.RFC3339Nano), entry.Watt)
		b.WriteString(s)

		if i%1e5 == 0 { // max speed 1e5, 1e4 is fine 991k/s vs 1.1M/s
			// log.Printf("Break the writer: %d bytes\n", len(b.Bytes()))
			fw.Bufw.Write(b.Bytes())
			b.Reset()
		}

		stamp = stamp.Add(time.Second)
	}

	fw.Bufw.Write(b.Bytes())
	b.Reset()
	fw.Bufw.Flush()
	timer.Track(start, "sh.Add", length)

	// r := strings.NewReader(content)
	// cid, err := sh.Add(r, shell.Pin(pin))
	cid, err := fw.Close()
	log.Printf("Added content(%d): %s\n", length, cid)
	return cid, err
}

func addOneDir(sh *shell.Shell) (string, error) {
	dw := ipfs.NewDWriter(sh)
	for day := range []int{1, 2, 3} {
		cid, err := addOneFile(sh, day, false)
		if err != nil {
			return "", err
		}
		path := fmt.Sprintf("day/2020-07-0%dT00:00:00Z.json", day)
		if _, err := dw.AddFile(path, cid); err != nil {
			return "", err
		}
		log.Printf("created file: %s : %s\n", path, cid)
	}
	// Finally we need to pin the top level dircid
	if err := dw.Close(); err != nil {
		return "", err
	}
	log.Printf("created dir: %s\n", dw.Dir)

	return dw.Dir, nil
}
