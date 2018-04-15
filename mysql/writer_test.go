package mysql

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
			sql:    "INSERT IGNORE INTO zzz (stamp, watt) VALUES (?,?)",
		}, {
			length: 2,
			sql:    "INSERT IGNORE INTO zzz (stamp, watt) VALUES (?,?),(?,?)",
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
