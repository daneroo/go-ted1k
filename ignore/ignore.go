package ignore

import (
	"time"

	. "github.com/daneroo/go-ted1k/types"
	. "github.com/daneroo/go-ted1k/util"
)

func Write(src <-chan Entry) {
	start := time.Now()
	count := 0
	for _ = range src {
		count++
	}
	TimeTrack(start, "ignore.Write", count)
}
