package types

import (
	"time"
)

type Entry struct {
	Stamp time.Time `json:"stamp"`
	Watt  int       `json:"watt"`
}
