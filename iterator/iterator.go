package iterator

import "github.com/daneroo/go-ted1k/types"

// Implements a pattern found in this article on iteration patters:
// iteration article - <https://blog.kowalczyk.info/article/1Bkr/3-ways-to-iterate-in-go.html>

// Entry is the interface for iterating over entries
type Entry interface {
	Next() bool
	// ValueIfPresent is like: a, aOk := <-aa
	ValueIfPresent() (types.Entry, bool)
	Value() types.Entry
	Error() error
}

// SliceIterator is the receiver for out slice entry iterator
type SliceIterator struct {
	src   <-chan []types.Entry
	slice []types.Entry // current slice
	err   error
}

// NewSliceIterator constructs a new Entry iterator from a channel of entry slices
func NewSliceIterator(src <-chan []types.Entry) *SliceIterator {
	return &SliceIterator{
		src: src,
	}
}

// Next implements the Entry interface
func (i *SliceIterator) Next() bool {
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

// ValueIfPresent implements the Entry interface
// ValueIfPresent is like: a, aOk := <-aa
func (i *SliceIterator) ValueIfPresent() (types.Entry, bool) {
	ok := i.Next()
	if ok {
		return i.Value(), ok // true
	}
	return types.Entry{}, ok // false
}

// Value implements the Entry interface
func (i *SliceIterator) Value() types.Entry {
	// shift/pop front:	x, a = a[0], a[1:]
	head, slice := i.slice[0], i.slice[1:]
	i.slice = slice
	return head
}

// Error implements the Entry interface
func (i *SliceIterator) Error() error {
	return i.err
}
