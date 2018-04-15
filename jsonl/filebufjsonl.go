package jsonl

// File Based (bufferd) JSON Encoder
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// File Buffered JSon Encoder

type FBJE struct {
	isOpen bool
	file   io.WriteCloser
	bufw   *bufio.Writer
	enc    *json.Encoder
}

func (fbje *FBJE) Open(fileName string) error {
	if fbje.isOpen {
		fbje.Close()
	}

	if file, err := os.Create(fileName); err != nil {
		return err
	} else {
		fbje.file = file
	}

	fbje.bufw = bufio.NewWriter(fbje.file) // default size 4k
	fbje.enc = json.NewEncoder(fbje.bufw)
	fbje.isOpen = true

	return nil
}

func (fbje *FBJE) Encode(v interface{}) error {
	if !fbje.isOpen {
		return fmt.Errorf("FBJE: Encoder is not open")
	}
	return fbje.enc.Encode(v)
}

func (fbje *FBJE) Close() error {
	if fbje.isOpen {
		fbje.bufw.Flush()
		fbje.file.Close()
		fbje.isOpen = false
	}
	//TODO(daneroo) Should we return an error if not open ??
	return nil
}
