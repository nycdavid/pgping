package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type SqlDB interface {
	Ping() error
}

type Connection interface {
	Open(string, string) (SqlDB, error)
}

type Logger interface {
	Print(...interface{})
}

type SystemLogger struct{}

func (sl *SystemLogger) Print(v ...interface{}) {
	log.Print(v)
}

type SystemConnection struct{}

func (sc *SystemConnection) Open(driverName, dataSourceName string) (SqlDB, error) {
	return sql.Open(driverName, dataSourceName)
}

func main() {
	l := &SystemLogger{}
	c := &SystemConnection{}
	os.Exit(realMain(l, c))
}

func realMain(l Logger, c Connection) int {
	db, err := openConnection(c)
	if err != nil {
		l.Print(err.Error())
		return 1
	}

	limit, err := strconv.Atoi(os.Getenv("PINGLIMIT"))
	if err != nil {
		limit = 1
	}

	i := 0
	err = db.Ping()
	for i < limit-1 && err != nil {
		err = db.Ping()
		i++
	}

	if err != nil {
		l.Print(err.Error())
		return 1
	}
	l.Print("Postgres server READY and accepting connections...")
	return 0
}

func openConnection(c Connection) (SqlDB, error) {
	addr := os.Getenv("PGCONN")
	if addr == "" {
		return nil, errors.New("PGCONN is empty")
	}
	return c.Open("postgres", addr)
}
