package types

import (
	"time"
)

type Entry struct {
	Stamp time.Time
	Watt  int
}
