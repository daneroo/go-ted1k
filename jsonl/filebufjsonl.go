package jsonl

// File Based (buffered) JSON Encoder
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
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
	fbje.bufw = bufio.NewWriter(fbje.file) // default size 4k
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
	return fbje.enc.Encode(v)
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
