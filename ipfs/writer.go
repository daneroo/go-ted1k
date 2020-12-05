package ipfs

import (
	"bufio"
	"bytes"
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

			err := w.writeOneEntry(entry)
			util.Checkerr(err)
		}
	}
	w.close()
	return count, nil
}

func (w *Writer) writeOneEntry(entry types.Entry) error {
	w.openFor(entry)
	return w.fw.WriteOneEntry(entry)
}

func (w *Writer) close() {
	path := pathFor(w.Grain, w.intvl)
	cid, _ := w.fw.Close()

	_, err := w.Dw.AddFile(path, cid)
	util.Checkerr(err)

	w.Dw.Close()
	// log.Printf("Final close: %s %s", w.intvl, w.Dw.Dir)
}

// Does 4 things; open File, buffer, encoder, Interval
// - if there is not yet an interval, set it
// - if we are past the end of the current interval
//   - close the the current writer
//   - open the next writer
func (w *Writer) openFor(entry types.Entry) {
	// could test Start==End (not initialized)
	// 1- Determine the interval the new Entry is in
	if w.intvl.Start.IsZero() {
		s := w.Grain.Floor(entry.Stamp)
		e := w.Grain.AddTo(s)
		w.intvl = timewalker.Interval{Start: s, End: e}
		// log.Printf("+Initial interval: %s : %s %s", w.Grain, entry.Stamp, w.intvl)
	}

	// 2- Close before Open - if we are in a new interval
	if !entry.Stamp.Before(w.intvl.End) {
		if w.fw.isOpen {
			path := pathFor(w.Grain, w.intvl)
			cid, _ := w.fw.Close()
			_, err := w.Dw.AddFile(path, cid)
			util.Checkerr(err)
			// log.Printf("Should close: %s", path)

			// new interval: for loop
			s := w.Grain.Floor(entry.Stamp)
			e := w.Grain.AddTo(s)
			w.intvl = timewalker.Interval{Start: s, End: e}
			// nupath := pathFor(w.Grain, w.intvl)
			// log.Printf("Should open: %s", nupath)
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
	dw.Dir = dir
	return dir, nil
}

// FWriter is
type FWriter struct {
	sh         *shell.Shell   // the ipfs shell
	isOpen     bool           // just to ensure no double close.
	w          *io.PipeWriter // where we write the encoded JSON - this is wrapped by the bufw
	Bufw       *bufio.Writer  // where we write the encoded JSON - this is exposed
	byteBuffer bytes.Buffer
	r          *io.PipeReader // the reader that is passed to sh.Add(r)
	cidChan    chan string    // will return the cid from the sh.Add(r) - or nothing
	errChan    chan error     // will receive the error from the sh.Add(r) - or nothing
}

// NewFWriter is a FWriter constructor
func NewFWriter(sh *shell.Shell, pin bool) *FWriter {
	r, w := io.Pipe()
	bufw := bufio.NewWriter(w) // buffer size does not matter much

	cidChan := make(chan string)
	errChan := make(chan error)

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
		w:       w,
		Bufw:    bufw,
		r:       r,
		cidChan: cidChan,
		errChan: errChan,
	}

}

// Close flushes the bufferd writer, closes the w:io.PipeWriter, and waits for cid,err from the spawned reader goroutine
func (fw *FWriter) Close() (string, error) {
	fw.Bufw.Write(fw.byteBuffer.Bytes())
	fw.byteBuffer.Reset()

	fw.Bufw.Flush()

	if err := fw.Bufw.Flush(); err != nil {
		return "", err
	}
	if err := fw.w.Close(); err != nil {
		return "", err
	}
	fw.isOpen = false
	select {
	case cid := <-fw.cidChan:
		// fmt.Println("received cid", cid)
		return cid, nil
	case err := <-fw.errChan:
		// fmt.Println("received error", err)
		return "", err
	}
}

// WriteOneEntry is ...
func (fw *FWriter) WriteOneEntry(entry types.Entry) error {
	s := fmt.Sprintf("{\"stamp\":\"%s\",\"watt\":%d}\n", entry.Stamp.Format(time.RFC3339Nano), entry.Watt)
	fw.byteBuffer.WriteString(s)

	if fw.byteBuffer.Len() >= 1024*4096 { // max speed 1e5, 1e4 is fine 991k/s vs 1.1M/s
		// log.Printf("Break the writer: %d bytes\n", fw.byteBuffer.Len())
		if _, err := fw.Bufw.Write(fw.byteBuffer.Bytes()); err != nil {
			return err
		}
		fw.byteBuffer.Reset()
	}
	return nil
}
