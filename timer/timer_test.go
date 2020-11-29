package timer

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	var data = []struct {
		elapsed time.Duration
		count   int
		fmt     string // expected
	}{
		{
			elapsed: time.Minute,
			count:   60,
			fmt:     "thing took 1m0s, rate ~ 1.0/s count: 60",
		}, {
			elapsed: 0,
			count:   10,
			fmt:     "thing took 0s, rate ~ +Inf/s count: 10",
		}, {
			elapsed: time.Minute,
			count:   0,
			fmt:     "thing took 1m0s, rate ~ 0.0/s count: 0",
		}, {
			elapsed: time.Duration(1124) * time.Second, //18m44
			count:   212737945,
			fmt:     "thing took 18m44s, rate ~ 189.3k/s count: 212737945",
		}, {
			elapsed: time.Duration(1124702) * time.Millisecond, //18m44.702
			count:   212737945,
			fmt:     "thing took 18m44.702s, rate ~ 189.2k/s count: 212737945",
		}, { // test rounding elapsed to millisecond
			elapsed: time.Duration(1124702134088) * time.Nanosecond, //18m44.702134088,
			count:   212737945,
			fmt:     "thing took 18m44.702s, rate ~ 189.2k/s count: 212737945",
		},
	}

	for idx, tt := range data {
		fmt := format(tt.elapsed, "thing", tt.count)
		if fmt != tt.fmt {
			t.Errorf("Expected format(%d) to be %v, but it was %v instead.", idx, tt.fmt, fmt)
		}
	}
}

func TestRate(t *testing.T) {
	var data = []struct {
		elapsed time.Duration
		count   int
		fmt     string // expected
	}{
		{
			elapsed: time.Minute,
			count:   60,
			fmt:     "rate ~ 1.0/s",
		}, {
			elapsed: time.Minute,
			count:   0,
			fmt:     "rate ~ 0.0/s",
		}, {
			elapsed: 0,
			count:   10,
			fmt:     "rate ~ +Inf/s",
		}, {
			elapsed: 0,
			count:   0,
			fmt:     "rate ~ 0.0/s",
		}, {
			elapsed: time.Second,
			count:   123,
			fmt:     "rate ~ 123.0/s",
		}, {
			elapsed: time.Second,
			count:   123456,
			fmt:     "rate ~ 123.5k/s",
		}, {
			elapsed: time.Second,
			count:   123456789,
			fmt:     "rate ~ 123.5M/s",
		}, {
			elapsed: time.Duration(1124) * time.Second, //18m44
			count:   212737945,
			fmt:     "rate ~ 189.3k/s",
		}, {
			elapsed: time.Duration(1124702) * time.Millisecond, //18m44.702
			count:   212737945,
			fmt:     "rate ~ 189.2k/s",
		}, { // test rounding elapsed to millisecond
			elapsed: time.Duration(1124702134088) * time.Nanosecond, //18m44.702134088,
			count:   212737945,
			fmt:     "rate ~ 189.2k/s",
		},
	}

	for idx, tt := range data {
		fmt := Rate(tt.elapsed, tt.count)
		if fmt != tt.fmt {
			t.Errorf("Expected format(%d) to be %v, but it was %v instead.", idx, tt.fmt, fmt)
		}
	}
}
