package util

import (
	// "github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func Add(x, y int) int {
	return x + y
}

var numberset = []struct {
	x      int
	y      int
	result int
}{
	{1, 2, 3},
	{2, 2, 4},
	{3, 3, 6},
	{-1, 1, 0},
}

func TestSomething(t *testing.T) {
	for _, set := range numberset {
		aresult := Add(set.x, set.y)
		if aresult != set.result {
			t.Errorf("Expected %d+%d==%d, got %d", set.x, set.y, set.result, aresult)
		}
	}
}
