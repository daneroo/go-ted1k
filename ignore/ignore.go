package ignore

import (
	"time"

	. "github.com/daneroo/go-mysqltest/types"
	. "github.com/daneroo/go-mysqltest/util"
)

func Write(src <-chan Entry) {
	start := time.Now()
	count := 0
	for _ = range src {
		count++
	}
	TimeTrack(start, "ignore.Write", count)
}
