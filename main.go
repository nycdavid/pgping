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

type Logger interface {
	Print(...interface{})
}

type SystemLogger struct{}

func (sl *SystemLogger) Print(v ...interface{}) {
	log.Print(v)
}

type PostgresConnection struct{}

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
	os.Exit(
		realMain(
			&PostgresConnection{},
			&SystemLogger{},
		),
	)
}

func realMain(c Connection, l Logger) int {
	cs := os.Getenv("PGCONN")
	if cs == "" {
		l.Print("PGCONN environment variable cannot be blank")
		return 1
	}
	_, err := c.Open("postgres", cs)
	if err != nil {
		return 0
	}
	return 1
}
