package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	pgConnStr := os.Getenv("PGCONN")
	if pgConnStr == "" {
		log.Print("PGCONN environment variable cannot be blank")
		os.Exit(1)
	}
	db, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Print("Postgres server READY and accepting connections...")
	os.Exit(0)
}
