package ipfs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/daneroo/go-ted1k/types"
	"github.com/daneroo/go-ted1k/util"
	"github.com/daneroo/timewalker"
	shell "github.com/ipfs/go-ipfs-api"
)

// Writer is ...
type Writer struct {
	Grain timewalker.Duration
	intvl timewalker.Interval
	Dw    *DWriter
	fw    *FWriter
}

// NewWriter is a constructor for the Writer struct
func NewWriter(sh *shell.Shell) *Writer {
	return &Writer{
		Grain: timewalker.Month,
		Dw:    NewDWriter(sh),
		fw:    NewFWriter(sh, false), // no pin
	}
}

// Write consumes an Entry channel - returns (count,error)
// preforming batched writes (of size writeBatchSize)
func (w *Writer) Write(src <-chan []types.Entry) (int, error) {
	count := 0

	for slice := range src {
		for _, entry := range slice {
			count++

			w.openFor(entry)
			// err := w.fw.Enc.Encode(&entry)
			_, err := fmt.Fprintf(w.fw.W, "{\"stamp\":\"%s\",\"watt\":%d}\n", entry.Stamp.Format(time.RFC3339Nano), entry.Watt)
			util.Checkerr(err)
		}
	}
	w.close()
	return count, nil
}

func (w *Writer) close() {
	log.Printf("Final close: %s", w.intvl)
	path := pathFor(w.Grain, w.intvl)
	cid, _ := w.fw.Close()
	_, err := w.Dw.AddFile(path, cid)
	util.Checkerr(err)

	w.Dw.Close()
	log.Printf("Final close: %s", w.intvl)
	log.Printf("Final close: %s", w.Dw.Dir)
}

// Does 4 things; open File, buffer, encoder, Interval
func (w *Writer) openFor(entry types.Entry) {
	// could test Start==End (not initialized)
	// 1- Determine the interval the new Entry is in
	if !w.intvl.Start.IsZero() {
		// log.Printf("-I: %s : %s %s", w.Grain, entry.Stamp, w.intvl)
	} else {
		s := w.Grain.Floor(entry.Stamp)
		e := w.Grain.AddTo(s)
		w.intvl = timewalker.Interval{Start: s, End: e}
		log.Printf("+I: %s : %s %s", w.Grain, entry.Stamp, w.intvl)
	}

	// 2- Close before Open - if we are in a new interval
	if !entry.Stamp.Before(w.intvl.End) {
		if w.fw.isOpen {
			log.Printf("Should close: %s", w.intvl)
			path := pathFor(w.Grain, w.intvl)
			cid, _ := w.fw.Close()
			_, err := w.Dw.AddFile(path, cid)
			util.Checkerr(err)

			// new interval: for loop
			s := w.Grain.Floor(entry.Stamp)
			e := w.Grain.AddTo(s)
			w.intvl = timewalker.Interval{Start: s, End: e}
			log.Printf("Should open: %s", w.intvl)
			w.fw = NewFWriter(w.Dw.sh, false) // no pin
		}
	}
}

// DWriter is ...
type DWriter struct {
	sh  *shell.Shell // the ipfs shell
	Dir string       // the current directory cid
}

// NewDWriter is a DWriter constructor
func NewDWriter(sh *shell.Shell) *DWriter {
	dir, err := sh.NewObject("unixfs-dir")
	if err != nil {
		log.Fatalf("unable to create ipfs directory (unixfs-dir)")
	}
	log.Printf("- building dir: %s\n", dir)
	return &DWriter{
		sh:  sh,
		Dir: dir,
	}
}

// Close simply pins the current directory cid (dw.Dir)
func (dw *DWriter) Close() error {
	return dw.sh.Pin(dw.Dir)
}

// AddFile adds a file to the current directory (dw.Dir) and update the cid of the current  directory (dw.Dir)
//  for convenience, also return the current Dir cid
func (dw *DWriter) AddFile(path, cid string) (string, error) {
	create := true
	dir, err := dw.sh.PatchLink(dw.Dir, path, cid, create)
	if err != nil {
		return "", err
	}
	log.Printf("+ building dir: %s\n", dir)
	dw.Dir = dir
	return dir, nil
}

// FWriter is ...
type FWriter struct {
	sh      *shell.Shell   // the ipfs shell
	isOpen  bool           // just to ensure no double close.
	W       *io.PipeWriter // where we write the encoded JSON - this is exposed
	Enc     *json.Encoder  // a json encoder, which write to the W writer
	r       *io.PipeReader // the reader that is passed to sh.Add(r)
	cidChan chan string    // will return the cid from the sh.Add(r) - or nothing
	errChan chan error     // will receive the error from the sh.Add(r) - or nothing
}

// NewFWriter is a FWriter constructor
func NewFWriter(sh *shell.Shell, pin bool) *FWriter {
	r, w := io.Pipe()
	cidChan := make(chan string)
	errChan := make(chan error)
	enc := json.NewEncoder(w)

	go func() {
		cid, err := sh.Add(r, shell.Pin(pin))
		if err != nil {
			log.Println("sh.Add() error")
			close(cidChan)
			errChan <- err
			close(errChan)
			return
		}
		cidChan <- cid
		close(cidChan)
		close(errChan)
	}()

	return &FWriter{
		sh:      sh,
		isOpen:  true,
		W:       w,
		Enc:     enc,
		r:       r,
		cidChan: cidChan,
		errChan: errChan,
	}

}

// Close closes the W:io.PipeWriter, and waits for cid,err from the spawner goroutine
func (fw *FWriter) Close() (string, error) {
	err := fw.W.Close()
	if err != nil {
		return "", err
	}
	fw.isOpen = false
	select {
	case cid := <-fw.cidChan:
		// fmt.Println("received cid", cid)
		return cid, nil
	case err = <-fw.errChan:
		// fmt.Println("received error", err)
		return "", err
	}
}
