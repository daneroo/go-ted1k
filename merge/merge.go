package merge

import (
	"fmt"
	"time"

	"github.com/daneroo/go-ted1k/types"
)

// Verify compares a zipEntry stream to consecutive match types
func Verify(aa <-chan types.Entry, bb <-chan types.Entry) []string {
	zip := zip(aa, bb)

	vv := make([]string, 0)
	first := zipEntry{match: -1}
	last := zipEntry{match: -1}

	for ze := range zip {

		if first.match == -1 {
			first = ze
		}
		if ze.match != first.match {
			vv = append(vv, fmt.Sprintf("[%s, %s] %s",
				first.entry.Stamp.Format(time.RFC3339),
				last.entry.Stamp.Format(time.RFC3339),
				first.match.String(),
			))

			first = ze
		}
		last = ze
	}

	if last.match != -1 {
		vv = append(vv, fmt.Sprintf("[%s, %s] %s",
			first.entry.Stamp.Format(time.RFC3339),
			last.entry.Stamp.Format(time.RFC3339),
			first.match.String(),
		))
	}

	return vv
}

// Type classifies each entry in Equal,Conflic,MissingInA ot MissingInB
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
	entry types.Entry //anonymous field Car
	match Type
}

// Zip compares two Entry channels, expecting the values to be time sorted,
// classifying each distinct time ordered entry as either, equal,conflict,MissingInA,MissingInB
// the done channel signifies termination
func zip(aa <-chan types.Entry, bb <-chan types.Entry) <-chan zipEntry {
	zip := make(chan zipEntry)
	go func() {
		a, aOk := <-aa
		b, bOk := <-bb
		for aOk && bOk {
			// fmt.Printf("a: %v %v b: %v %v\n", a, aOk, b, bOk)
			if a.Stamp.Equal(b.Stamp) {
				if a.Watt == b.Watt {
					zip <- zipEntry{entry: a, match: Equal}
				} else {
					zip <- zipEntry{entry: a, match: Conflict}
				}
				a, aOk = <-aa
				b, bOk = <-bb
			} else if a.Stamp.Before(b.Stamp) {
				zip <- zipEntry{entry: a, match: MissingInB}
				a, aOk = <-aa
			} else if a.Stamp.After(b.Stamp) {
				zip <- zipEntry{entry: b, match: MissingInA}
				b, bOk = <-bb
			}
		}
		for aOk { // drain A channel
			zip <- zipEntry{entry: a, match: MissingInB}
			a, aOk = <-aa
		}
		for bOk { // drain B channel
			zip <- zipEntry{entry: b, match: MissingInA}
			b, bOk = <-bb
		}
		close(zip)
	}()
	return zip
}
