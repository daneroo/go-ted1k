package types

import (
	"time"
)

// Entry is the basic datatype for an energy measurement
type Entry struct {
	Stamp time.Time `json:"stamp"`
	Watt  int       `json:"watt"`
}
