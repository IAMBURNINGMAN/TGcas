package storage

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect(dbURL string) *sql.DB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("db.Ping: %v", err)
	}
	return db
}
