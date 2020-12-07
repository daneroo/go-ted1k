package main

import (
	"fmt"
	"log"
	"time"

	"github.com/daneroo/go-ted1k/ephemeral"
	"github.com/daneroo/go-ted1k/iterator"
	"github.com/daneroo/go-ted1k/logsetup"
	"github.com/daneroo/go-ted1k/merge"
	"github.com/daneroo/go-ted1k/progress"
	"github.com/daneroo/go-ted1k/timer"
	"github.com/daneroo/go-ted1k/types"
)

func main() {
	logsetup.SetupFormat()
	log.Printf("Starting iterator test\n")

	totalCount := int(31415926)
	// totalCount := int(3)
	performIteration("Reference Entry Iterator", newEntryIterator(totalCount))
	eatSlices("Explicit Slice Channel Iterator")
	performIteration("Slice Channel Adaptor Iterator", iterator.NewSliceIterator(ephemeral.NewReader().Read()))
	verify("eph <-> eph", ephemeral.NewReader(), ephemeral.NewReader())
	verify("eph <-> eph", ephemeral.NewReader(), ephemeral.NewReader())
	verify("eph <-> eph", ephemeral.NewReader(), ephemeral.NewReader())
	verify("eph <-> eph", ephemeral.NewReader(), ephemeral.NewReader())
}

func verify(name string, a, b types.EntryReader) {
	log.Printf("-=- %s\n", name)
	vv := merge.Verify(a.Read(), progress.Monitor(name, b.Read()))
	log.Printf("Verified %s:\n", name)
	for _, v := range vv {
		log.Println(v)
	}
}

func eatSlices(name string) {
	start := time.Now()
	count := 0
	reader := ephemeral.NewReader()
	src := reader.Read()
	for slice := range src {
		for _, entry := range slice { // index,entry
			if entry.Watt >= 0 {
				count++
			}
		}
	}

	timer.Track(start, fmt.Sprintf("%33s", name), count)
}

func performIteration(name string, iter iterator.Entry) {
	start := time.Now()
	count := 0
	for iter.Next() {
		entry := iter.Value()
		if entry.Watt >= 0 {
			count++
		}
	}
	if iter.Error() != nil {
		log.Fatalf("%s, error: %s\n", name, iter.Error())
	}
	timer.Track(start, fmt.Sprintf("%33s", name), count)
}

type entryIterator struct {
	Epoch      time.Time
	TotalCount int
	// These are my state variables
	currentCount int
	stamp        time.Time
	// ustamp       int64
	err error
}

func newEntryIterator(totalCount int) *entryIterator {
	var err error
	if totalCount < 0 {
		err = fmt.Errorf("'totalCount' is %d, should be >= 0", totalCount)
	}
	stamp := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	// ustamp := stamp.Unix()
	return &entryIterator{
		Epoch:      stamp,
		TotalCount: totalCount,
		stamp:      stamp,
		// ustamp:     ustamp,
		err: err,
	}
}

func (i *entryIterator) Next() bool {
	if i.err != nil {
		return false
	}
	i.currentCount++
	i.stamp = i.stamp.Add(time.Second)
	// i.ustamp++

	return i.currentCount <= i.TotalCount
}

// ValueIfPresent implements the Entry interface
// ValueIfPresent is like: a, aOk := <-aa
func (i *entryIterator) ValueIfPresent() (types.Entry, bool) {
	ok := i.Next()
	if ok {
		return i.Value(), ok // true
	}
	return types.Entry{}, ok // false
}

func (i *entryIterator) Value() types.Entry {
	if i.err != nil || !(i.currentCount <= i.TotalCount) {
		panic("Value is not valid after iterator finished")
	}
	entry := types.Entry{Stamp: i.stamp, Watt: int(i.stamp.Unix() % 5000)}
	// stamp := time.Unix(i.ustamp, 0)
	// entry := types.Entry{Stamp: stamp, Watt: int(i.ustamp % 5000)}
	return entry
}

func (i *entryIterator) Error() error {
	return i.err
}
