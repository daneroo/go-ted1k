package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

// Setup is ...
func Setup(ctx context.Context, tableNames []string, credentials string) *pgx.Conn {
	conn, err := pgx.Connect(ctx, credentials)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		panic(err)
	}
	log.Println("Connected to Postgres")

	for _, tableName := range tableNames {
		createTable(ctx, conn, tableName)
	}
	// totalCount(ctx, conn)

	return conn
}

func createTable(ctx context.Context, conn *pgx.Conn, tableName string) {
	ddlFormat := "CREATE TABLE IF NOT EXISTS %s (stamp TIMESTAMP WITHOUT TIME ZONE NOT NULL PRIMARY KEY,watt integer NOT NULL DEFAULT '0');"
	ddl := fmt.Sprintf(ddlFormat, tableName)
	_, err := conn.Exec(ctx, ddl)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	sqlCreateHyperFormat := "SELECT create_hypertable('%s', 'stamp')"
	sqlCreateHyper := fmt.Sprintf(sqlCreateHyperFormat, tableName)
	log.Println(sqlCreateHyper)
	_, err = conn.Exec(ctx, sqlCreateHyper)
	if err != nil {
		log.Println(err)
		// panic(err)
	}
}

func totalCount(ctx context.Context, conn *pgx.Conn) {
	var totalCount int
	err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM watt").Scan(&totalCount)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Printf("Found %d entries in watt\n", totalCount)
}
