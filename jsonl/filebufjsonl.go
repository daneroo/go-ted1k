package jsonl

// File Based (buffered) JSON Encoder
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/daneroo/go-ted1k/types"
)

// FBJE File Buffered JSon Encoder
type FBJE struct {
	isOpen bool
	file   io.WriteCloser
	bufw   *bufio.Writer
	enc    *json.Encoder
}

// Open is ...
func (fbje *FBJE) Open(fileName string) error {
	if fbje.isOpen {
		fbje.Close()
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	fbje.file = file
	fbje.bufw = bufio.NewWriter(fbje.file) // default size 4k is fine
	fbje.enc = json.NewEncoder(fbje.bufw)
	fbje.isOpen = true

	return nil
}

// Encode is ...
func (fbje *FBJE) Encode(v interface{}) error {
	if !fbje.isOpen {
		err := fmt.Errorf("FBJE: Encoder is not open")
		return err
	}
	//
	// return fbje.enc.Encode(v)
	// This is faster: 1.1M/s vs 670k/s for both Day and Month Grain,
	entry := v.(*types.Entry)
	_, err := fmt.Fprintf(fbje.bufw, "{\"stamp\":\"%s\",\"watt\":%d}\n", entry.Stamp.Format(time.RFC3339Nano), entry.Watt)
	return err
}

// Close is ...
func (fbje *FBJE) Close() error {
	if fbje.isOpen {
		fbje.bufw.Flush()
		fbje.file.Close()
		fbje.isOpen = false
	}
	//TODO(daneroo) Should we return an error if not open ??
	return nil
}
