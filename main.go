package main

import (
	"database/sql"
	"errors"
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

type DBHandle interface {
	Ping() error
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
	err := checkConnectionStr()
	if err != nil {
		l.Print(err.Error())
		return 1
	}

	h, err := openConnection(c)
	if err != nil {
		l.Print(err.Error())
		return 1
	}

	err = testLine(h)
	if err != nil {
		l.Print(err.Error())
		return 1
	}
	return 0
}

func testLine(h DBHandle) error {
	return h.Ping()
}

func checkConnectionStr() error {
	cs := os.Getenv("PGCONN")
	if cs == "" {
		return errors.New("PGCONN environment variable cannot be blank")
	}
	return nil
}

func openConnection(c Connection) (DBHandle, error) {
	cs := os.Getenv("PGCONN")
	h, err := c.Open("postgres", cs)
	if err != nil {
		return nil, err
	}
	return h, nil
}
