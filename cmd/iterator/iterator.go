package main

import (
	"fmt"
	"log"
	"time"

	"github.com/daneroo/go-ted1k/ephemeral"
	"github.com/daneroo/go-ted1k/logsetup"
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
	performIteration("Slice Channel Adaptor Iterator", newSliceIterator())

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

func performIteration(name string, iter iterator) {
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

// Iterator is ... - first we do this with a struct, then with a closure
type iterator interface {
	Next() bool
	Value() types.Entry
	Error() error
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

type sliceIterator struct {
	reader *ephemeral.Reader
	src    <-chan []types.Entry
	slice  []types.Entry // current slice
	err    error
}

func newSliceIterator() *sliceIterator {
	reader := ephemeral.NewReader()
	src := reader.Read()

	return &sliceIterator{
		reader: reader,
		src:    src,
	}
}

// for slice := range src {
// 	for _, entry := range slice { // index,entry
// 		if entry.Watt >= 0 {
// 			count++
// 		}
// 	}
// }

func (i *sliceIterator) Next() bool {
	// if current i.slice is non empty, return true (there are more items)
	// otherwise, get next slice from src channel (until we have a non-empty one)
	// if the channel is closed, we have no more items
	// note: the zero element (initialized in the struct) is an empty slice
	for len(i.slice) == 0 {
		// fetch the nxt slice if there is one
		if slice, ok := <-i.src; ok {
			i.slice = slice // strore is struct state
		} else {
			return false
		}
	}
	return true
}

func (i *sliceIterator) Value() types.Entry {
	// shift/pop front:	x, a = a[0], a[1:]
	head, slice := i.slice[0], i.slice[1:]
	i.slice = slice
	return head
}

func (i *sliceIterator) Error() error {
	return i.err
}
