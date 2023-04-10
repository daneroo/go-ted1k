package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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

// This is the only method that has been error-ified!
// SetupPool is ...
// TODO(daneroo): The same error treatment should be applied everywhere.
func SetupPool(ctx context.Context, tableNames []string, credentials string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, credentials)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	log.Println("Connected to Postgres")

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Printf("Unable to acquire connection from pool: %v\n", err)
		return nil, err
	}
	defer conn.Release()
	for _, tableName := range tableNames {
		createTable(ctx, conn.Conn(), tableName)
	}
	// totalCount(ctx, conn)
	return pool, nil
}

func createTable(ctx context.Context, conn *pgx.Conn, tableName string) {
	ddlFormat := "CREATE TABLE IF NOT EXISTS %s (stamp TIMESTAMPTZ NOT NULL PRIMARY KEY,watt integer NOT NULL DEFAULT '0');"
	ddl := fmt.Sprintf(ddlFormat, tableName)
	_, err := conn.Exec(ctx, ddl)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	// Confirm that the database is in UTC:
	// SELECT current_setting('TIMEZONE'); ==> UTC
	var currentTZ string
	err = conn.QueryRow(context.Background(), "SELECT current_setting('TIMEZONE')").Scan(&currentTZ)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	if currentTZ != "UTC" {
		log.Printf("Warning: Current timezone: %s != UTC\n", currentTZ)
	} else {
		log.Printf("Confirmed Current timezone: %s == UTC\n", currentTZ)
	}

	// Only create the hypertable if it doesn't already exist.
	// previous unconditional: sqlCreateHyperFormat := "SELECT create_hypertable('%s', 'stamp')"
	sqlCreateHyperFormat := "SELECT create_hypertable('%s', 'stamp') WHERE NOT EXISTS (SELECT 1 FROM _timescaledb_catalog.hypertable WHERE table_name = '%s')"
	sqlCreateHyper := fmt.Sprintf(sqlCreateHyperFormat, tableName, tableName)
	_, err = conn.Exec(ctx, sqlCreateHyper)
	if err != nil {
		log.Println(sqlCreateHyper)
		log.Println(err)
		panic(err)
	}
}

func TotalCount(ctx context.Context, conn *pgx.Conn) {
	var totalCount int
	err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM watt").Scan(&totalCount)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Printf("Found %d entries in watt\n", totalCount)
}
