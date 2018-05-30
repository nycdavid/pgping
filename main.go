package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

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

type SysSqlDB struct {
	PingFunc func() error
}

func (sysDB *SysSqlDB) Ping() error {
	return sysDB.PingFunc()
}

func main() {
	// err = db.Ping()
	// if err != nil {
	// 	log.Print(err)
	// 	os.Exit(1)
	// }
	// log.Print("Postgres server READY and accepting connections...")
	// os.Exit(0)
	l := &SystemLogger{}
	os.Exit(
		realMain(l),
	)
}

func realMain(l Logger) int {
	_, err := openConnection(l)
	if err != nil {
		return 1
	}
	return 0
}

func openConnection(l Logger) (SqlDB, error) {
	addr := os.Getenv("PGCONN")
	if addr == "" {
		return nil, errors.New("PGCONN is empty")
	}
	return sql.Open("postgres", addr)
}
