package storage

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) {
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("migrate.New: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate.Up: %v", err)
	}
	log.Println("миграции применены")
}
