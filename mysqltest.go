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

	const maxCountPerChunk = 3600 * 24
	rowCount := 0
	startTimeExcl := time.Date(2007, time.January, 0, 0, 0, 0, 0, time.UTC)
	// startTimeExcl := time.Date(2037, time.January, 0, 0, 0, 0, 0, time.UTC)
	for {
		chunkRowCount := 0
		startTimeExcl, chunkRowCount = oneChunk(db, startTimeExcl, maxCountPerChunk)
		rowCount += chunkRowCount
		if chunkRowCount == 0 {
			break
		}
	}
	log.Printf("Fetched a total of %d rows (%d before iteration)", rowCount, totalCount)

}

func oneChunk(db *sql.DB, startTimeExcl time.Time, maxCountPerChunk int) (time.Time, int) {
	sql := "SELECT stamp,watt FROM watt where stamp>? ORDER BY stamp ASC LIMIT ?"
	// sql := "SELECT stamp,watt FROM watt where stamp<? ORDER BY stamp DESC LIMIT ?"
	rows, err := db.Query(sql, startTimeExcl, maxCountPerChunk)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	avgWatt := 0
	chunkRowCount := 0
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
		chunkRowCount++
		if stamp.Valid {
			lastStamp = stamp.Time
		}
		// log.Printf(" %v: %v", stamp, watt)
	}
	err = rows.Err() // get any error encountered during iteration
	if err != nil {
		log.Println(err)
	}
	if chunkRowCount != 0 {
		avgWatt /= chunkRowCount
	}
	log.Printf("average between (%v - %v]: %v (%v)", startTimeExcl, lastStamp, avgWatt, chunkRowCount)
	return lastStamp, chunkRowCount
}
