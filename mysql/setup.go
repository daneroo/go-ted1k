package mysql

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// Setup is ...
func Setup(tableNames []string, credentials string) *sqlx.DB {
	// Connect is Open and verify with a Ping
	db, err := sqlx.Connect("mysql", credentials)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println("Connected to MySQL")

	for _, tableName := range tableNames {
		createCopyTable(db, tableName)
	}
	totalCount(db)

	return db
}

func createCopyTable(db *sqlx.DB, tableName string) {
	ddlFormat := "CREATE TABLE IF NOT EXISTS %s ( stamp datetime NOT NULL DEFAULT '1970-01-01 00:00:00', watt int(11) NOT NULL DEFAULT '0',  PRIMARY KEY (`stamp`) )  ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;"
	ddl := fmt.Sprintf(ddlFormat, tableName)
	_, err := db.Exec(ddl)
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func totalCount(db *sqlx.DB) {
	var totalCount int
	err := db.Get(&totalCount, "SELECT COUNT(*) FROM watt")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Printf("Found %d entries in watt\n", totalCount)
}
