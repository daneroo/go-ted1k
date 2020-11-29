package merge

import (
	"reflect"
	"testing"
	"time"

	"github.com/daneroo/go-ted1k/types"
)

func TestMergeTypeString(t *testing.T) {
	var data = []struct {
		m Type
		s string // expected
	}{
		{
			m: Equal,
			s: "Equal",
		}, {
			m: Conflict,
			s: "Conflict",
		}, {
			m: MissingInA,
			s: "MissingInA",
		}, {
			m: MissingInB,
			s: "MissingInB",
		}, {
			m: -1,
			s: "Unknown",
		}, {
			m: 100,
			s: "Unknown",
		},
	}

	for idx, tt := range data {

		s := tt.m.String()

		if !reflect.DeepEqual(s, tt.s) {
			t.Errorf("Expected srting(%d) to be\n%v, but it was \n%v\ninstead.", idx, tt.s, s)
		}
	}
}

func TestVerify(t *testing.T) {
	var data = []struct {
		a   <-chan types.Entry
		b   <-chan types.Entry
		msg []string // expected
	}{
		{
			a:   chanFromSlice([]int{}),
			b:   chanFromSlice([]int{}),
			msg: []string{},
		}, {
			a:   chanFromSlice([]int{1000}),
			b:   chanFromSlice([]int{1000}),
			msg: []string{"[2016-01-01T00:00:00Z, 2016-01-01T00:00:00Z] Equal"},
		}, {
			a:   chanFromSlice([]int{1000}),
			b:   chanFromSlice([]int{2000}),
			msg: []string{"[2016-01-01T00:00:00Z, 2016-01-01T00:00:00Z] Conflict"},
		}, {
			a:   chanFromSlice([]int{1000, 2000}),
			b:   chanFromSlice([]int{1000, 2000}),
			msg: []string{"[2016-01-01T00:00:00Z, 2016-01-01T00:00:01Z] Equal"},
		}, {
			a: chanFromSlice([]int{1000, 2000}),
			b: chanFromSlice([]int{1000, 2000, 3000}),
			msg: []string{
				"[2016-01-01T00:00:00Z, 2016-01-01T00:00:01Z] Equal",
				"[2016-01-01T00:00:02Z, 2016-01-01T00:00:02Z] MissingInA",
			},
		}, {
			a: chanFromSlice([]int{1000, 2000, 3000}),
			b: chanFromSlice([]int{1000, 2000}),
			msg: []string{
				"[2016-01-01T00:00:00Z, 2016-01-01T00:00:01Z] Equal",
				"[2016-01-01T00:00:02Z, 2016-01-01T00:00:02Z] MissingInB",
			},
		}, {
			a: chanFromSlice([]int{1000, -1, 3000}),
			b: chanFromSlice([]int{1000, 2000, 3000}),
			msg: []string{
				"[2016-01-01T00:00:00Z, 2016-01-01T00:00:00Z] Equal",
				"[2016-01-01T00:00:01Z, 2016-01-01T00:00:01Z] MissingInA",
				"[2016-01-01T00:00:02Z, 2016-01-01T00:00:02Z] Equal",
			},
		}, {
			a: chanFromSlice([]int{1000, 2000, 3000}),
			b: chanFromSlice([]int{1000, -1, 3000}),
			msg: []string{
				"[2016-01-01T00:00:00Z, 2016-01-01T00:00:00Z] Equal",
				"[2016-01-01T00:00:01Z, 2016-01-01T00:00:01Z] MissingInB",
				"[2016-01-01T00:00:02Z, 2016-01-01T00:00:02Z] Equal",
			},
		},
	}

	for idx, tt := range data {

		msg := Verify(tt.a, tt.b)

		if !reflect.DeepEqual(msg, tt.msg) {
			t.Errorf("Expected msg(%d) to be\n%v, but it was \n%v\ninstead.", idx, tt.msg, msg)
		}
	}
}

func TestZip(t *testing.T) {
	var data = []struct {
		a    <-chan types.Entry
		b    <-chan types.Entry
		zips []zipEntry // expected
	}{
		{
			a:    chanFromSlice([]int{}),
			b:    chanFromSlice([]int{}),
			zips: []zipEntry{},
		}, {
			a: chanFromSlice([]int{1000}),
			b: chanFromSlice([]int{1000}),
			zips: []zipEntry{
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:00Z"), Watt: 1000}, match: Equal},
			},
		}, {
			a: chanFromSlice([]int{1000}),
			b: chanFromSlice([]int{2000}),
			zips: []zipEntry{
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:00Z"), Watt: 1000}, match: Conflict},
			},
		}, {
			a: chanFromSlice([]int{1000, 2000}),
			b: chanFromSlice([]int{1000, 2000}),
			zips: []zipEntry{
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:00Z"), Watt: 1000}, match: Equal},
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:01Z"), Watt: 2000}, match: Equal},
			},
		}, {
			a: chanFromSlice([]int{1000, 2000}),
			b: chanFromSlice([]int{1000, 2000, 3000}),
			zips: []zipEntry{
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:00Z"), Watt: 1000}, match: Equal},
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:01Z"), Watt: 2000}, match: Equal},
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:02Z"), Watt: 3000}, match: MissingInA},
			},
		}, {
			a: chanFromSlice([]int{1000, 2000, 3000}),
			b: chanFromSlice([]int{1000, 2000}),
			zips: []zipEntry{
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:00Z"), Watt: 1000}, match: Equal},
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:01Z"), Watt: 2000}, match: Equal},
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:02Z"), Watt: 3000}, match: MissingInB},
			},
		}, {
			a: chanFromSlice([]int{1000, -1, 3000}),
			b: chanFromSlice([]int{1000, 2000, 3000}),
			zips: []zipEntry{
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:00Z"), Watt: 1000}, match: Equal},
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:01Z"), Watt: 2000}, match: MissingInA},
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:02Z"), Watt: 3000}, match: Equal},
			},
		}, {
			a: chanFromSlice([]int{1000, 2000, 3000}),
			b: chanFromSlice([]int{1000, -1, 3000}),
			zips: []zipEntry{
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:00Z"), Watt: 1000}, match: Equal},
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:01Z"), Watt: 2000}, match: MissingInB},
				{entry: types.Entry{Stamp: fromStamp("2016-01-01T00:00:02Z"), Watt: 3000}, match: Equal},
			},
		},
	}

	for _, tt := range data {

		zips := sliceFromChan(zip(tt.a, tt.b))

		if !reflect.DeepEqual(zips, tt.zips) {
			t.Errorf("Expected zips to be \n%v, but it was \n%v\ninstead.", tt.zips, zips)
		}
	}
}

func fromStamp(s string) time.Time {
	stamp, _ := time.Parse(time.RFC3339, s)
	return stamp
}

func sliceFromChan(zip <-chan zipEntry) []zipEntry {
	zips := []zipEntry{}
	for ze := range zip {
		zips = append(zips, ze)
	}
	return zips
}
func chanFromSlice(ww []int) <-chan types.Entry {
	src := make(chan types.Entry)
	stamp, _ := time.Parse(time.RFC3339, "2016-01-01T00:00:00Z")
	go func() {
		for _, w := range ww {
			if w >= 0 {
				src <- types.Entry{Stamp: stamp, Watt: w}
			}
			stamp = stamp.Add(time.Second)
		}
		close(src)
	}()
	return src
}
