package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
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

	rowCount, offset := 10, 0
	for offset = 0; offset <= totalCount; offset += rowCount {
		rows, err := db.Query("SELECT stamp,watt FROM watt order by stamp asc LIMIT ? OFFSET ?", rowCount, offset)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			// var stamp string
			var stamp mysql.NullTime
			var watt int
			err = rows.Scan(&stamp, &watt)
			if err != nil {
				log.Println(err)
			}
			log.Printf(" %v: %v", stamp, watt)
		}
		err = rows.Err() // get any error encountered during iteration

	}

}
