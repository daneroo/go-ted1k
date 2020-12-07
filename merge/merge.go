package merge

import (
	"fmt"
	"time"

	"github.com/daneroo/go-ted1k/iterator"
	"github.com/daneroo/go-ted1k/types"
)

// Verify compares a zipEntry stream to combine consecutive match types
func Verify(aa <-chan []types.Entry, bb <-chan []types.Entry) []string {
	zip := compare(iterator.NewSliceIterator(aa), iterator.NewSliceIterator(bb))
	vv := combineConsecutiveZipEntries(zip)
	return vv
}

func combineConsecutiveZipEntries(zipslice <-chan []zipEntry) []string {
	vv := make([]string, 0)
	first := zipEntry{match: -1}
	last := zipEntry{match: -1}
	consecutive := 0

	for zip := range zipslice {
		for _, ze := range zip {
			if first.match == -1 {
				first = ze
				consecutive = 0
			}
			if ze.match != first.match {
				vv = append(vv, fmt.Sprintf("[%s, %s](%d) %s",
					first.entry.Stamp.Format(time.RFC3339),
					last.entry.Stamp.Format(time.RFC3339),
					consecutive,
					first.match.String(),
				))

				first = ze
				consecutive = 0
			}
			consecutive++
			last = ze
		}
	}

	if last.match != -1 {
		vv = append(vv, fmt.Sprintf("[%s, %s](%d) %s",
			first.entry.Stamp.Format(time.RFC3339),
			last.entry.Stamp.Format(time.RFC3339),
			consecutive,
			first.match.String(),
		))
	} // else (no entries at all?)

	return vv
}

// Type classifies each entry in Equal,Conflict,MissingInA ot MissingInB
type Type int

const (
	// Equal denotes twow entries have identical values (for the same timestamp)
	Equal Type = iota // 0
	// Conflict denotes twow entries have conficting values (for the same timestamp)
	Conflict // 1
	// MissingInA denotes the entry is not present in the A channel
	MissingInA // 2
	// MissingInB denotes the entry is not present in the B channel
	MissingInB // 3
)

func (m Type) String() string {
	switch m {
	case Equal:
		return "Equal"
	case Conflict:
		return "Conflict"
	case MissingInA:
		return "MissingInA"
	case MissingInB:
		return "MissingInB"
	}
	return "Unknown"
}

type zipEntry struct {
	entry types.Entry
	match Type
}

const zipEntrySliceCapacity = 100

// compare compares two Entry channels, expecting the values to be time sorted,
// classifying each distinct time ordered entry as either, equal,conflict,MissingInA,MissingInB
func compare(aa, bb iterator.Entry) <-chan []zipEntry {
	zipslice := make([]zipEntry, 0, zipEntrySliceCapacity) // current slice

	zip := make(chan []zipEntry)
	go func() {
		a, aOk := aa.ValueIfPresent()
		b, bOk := bb.ValueIfPresent()
		for aOk && bOk {
			if a.Stamp.Equal(b.Stamp) {
				if a.Watt == b.Watt {
					zipslice = append(zipslice, zipEntry{entry: a, match: Equal})
				} else {
					zipslice = append(zipslice, zipEntry{entry: a, match: Conflict})
				}
				a, aOk = aa.ValueIfPresent()
				b, bOk = bb.ValueIfPresent()
			} else if a.Stamp.Before(b.Stamp) {
				zipslice = append(zipslice, zipEntry{entry: a, match: MissingInB})
				a, aOk = aa.ValueIfPresent()
			} else if a.Stamp.After(b.Stamp) {
				zipslice = append(zipslice, zipEntry{entry: b, match: MissingInA})
				b, bOk = bb.ValueIfPresent()
			}
			if len(zipslice) >= zipEntrySliceCapacity {
				zip <- zipslice
				zipslice = make([]zipEntry, 0, zipEntrySliceCapacity)
			}
		}
		for aOk { // drain A channel
			zipslice = append(zipslice, zipEntry{entry: a, match: MissingInB})
			if len(zipslice) >= zipEntrySliceCapacity {
				zip <- zipslice
				zipslice = make([]zipEntry, 0, zipEntrySliceCapacity)
			}
			a, aOk = aa.ValueIfPresent()

		}
		for bOk { // drain B channel
			zipslice = append(zipslice, zipEntry{entry: b, match: MissingInA})
			if len(zipslice) >= zipEntrySliceCapacity {
				zip <- zipslice
				zipslice = make([]zipEntry, 0, zipEntrySliceCapacity)
			}
			b, bOk = bb.ValueIfPresent()
		}
		// flush the slice to the channel
		if len(zipslice) > 0 {
			zip <- zipslice
		}
		close(zip)
	}()
	return zip
}

type entryIterator struct {
	src   <-chan types.Entry
	entry types.Entry
	err   error
}

func newEntryIterator(src <-chan types.Entry) *entryIterator {
	return &entryIterator{
		src: src,
	}
}

func (i *entryIterator) Next() bool {
	if entry, ok := <-i.src; ok {
		i.entry = entry // strore in struct state
	} else {
		return false
	}
	return true
}

func (i *entryIterator) Value() types.Entry {
	return i.entry
}

func (i *entryIterator) Error() error {
	return i.err
}

type sliceIterator struct {
	src   <-chan []types.Entry
	slice []types.Entry // current slice
	err   error
}

// func newSliceIterator(src <-chan []types.Entry) *sliceIterator {
// 	return &sliceIterator{
// 		src: src,
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
