package postgres

import (
	// "github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestMakeSQL(t *testing.T) {
	var data = []struct {
		length int    // input
		sql    string // expected
	}{
		{
			length: 1,
			sql:    "INSERT INTO zzz (stamp, watt) VALUES ($1,$2) ON CONFLICT (stamp) DO NOTHING",
		}, {
			length: 2,
			sql:    "INSERT INTO zzz (stamp, watt) VALUES ($1,$2),($3,$4) ON CONFLICT (stamp) DO NOTHING",
		}, {
			length: 0,
			sql:    "",
		},
	}

	for _, tt := range data {
		myWriter := &Writer{TableName: "zzz"}

		sql := myWriter.makeSQL(tt.length)

		if sql != tt.sql {
			t.Errorf("Expected sql to be %s, but it was %s instead.", tt.sql, sql)
		}
	}
}
