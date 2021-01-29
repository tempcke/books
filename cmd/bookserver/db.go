package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/tempcke/books/internal"
)

const verboseLogging = false

func dbMigrateUp(dsn string, log *internal.Logger) error {
	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		return err
	}

	m.Log = log
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
