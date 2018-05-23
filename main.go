package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Connection interface {
	Open(string, string) (*sql.DB, error)
}

type PostgresConnection struct {
}

func (pc *PostgresConnection) Open(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open("postgres", dataSourceName)
}

func main() {
	// pgConnStr := os.Getenv("PGCONN")
	// if pgConnStr == "" {
	// }
	//
	// db, err := sql.Open("postgres", pgConnStr)
	// if err != nil {
	// 	log.Print(err)
	// 	os.Exit(1)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	log.Print(err)
	// 	os.Exit(1)
	// }
	// log.Print("Postgres server READY and accepting connections...")
	// os.Exit(0)
	os.Exit(realMain(&PostgresConnection{}))
}

func realMain(c Connection) int {
	cs := os.Getenv("PGCONN")
	if cs == "" {
		log.Print("PGCONN environment variable cannot be blank")
		return 1
	}
	_, err := c.Open("postgres", cs)
	if err != nil {
		return 0
	}
	return 1
}
