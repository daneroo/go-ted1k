package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func main() {
	log.Printf("Just getting %s\n", "started")
	for i := 0; i < 2; i++ {
		log.Printf("working %v\n", i)
	}
	db, err := sql.Open("mysql", "daniel@tcp(192.168.5.105:3306)/ted")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	log.Println("Survived Opening")

	var totalCount int
	row := db.QueryRow("SELECT COUNT(*) FROM watt")
	err = row.Scan(&totalCount)

	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
	}
	log.Printf("Found %v entries\n", totalCount)

	rowCount, offset := 1000, 0
	startTimeExcl := time.Date(2007, time.January, 0, 0, 0, 0, 0, time.UTC)
	for offset = 0; offset <= totalCount; offset += rowCount {
		startTimeExcl = oneChunk(db, startTimeExcl, rowCount)
	}

}

func oneChunk(db *sql.DB, startTimeExcl time.Time, rowCount int) time.Time {
	rows, err := db.Query("SELECT stamp,watt FROM watt where stamp>? ORDER BY stamp ASC LIMIT ?", startTimeExcl, rowCount)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	avgWatt := 0
	count := 0
	var lastStamp time.Time
	for rows.Next() {
		// var stamp string
		var stamp mysql.NullTime
		var watt int
		err = rows.Scan(&stamp, &watt)
		if err != nil {
			log.Println(err)
		}
		avgWatt += watt
		count++
		if stamp.Valid {
			lastStamp = stamp.Time
		}
		// log.Printf(" %v: %v", stamp, watt)
	}
	err = rows.Err() // get any error encountered during iteration
	if err != nil {
		log.Println(err)
	}
	avgWatt /= count
	log.Printf("average between (%v - %v]: %v", startTimeExcl, lastStamp, avgWatt)
	return lastStamp
}
