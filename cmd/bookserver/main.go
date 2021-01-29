package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/tempcke/books/api/rest"
	"github.com/tempcke/books/fake"
	"github.com/tempcke/books/internal"
	"github.com/tempcke/books/repository"
	"github.com/tempcke/books/usecase"
)

const dbDriver = "postgres"

func main() {
	log := internal.NewLogger()
	conf := NewConfigFromEnv()
	if err := run(conf, log); err != nil {
		log.Fatal("BookServer Error: " + err.Error())
	}
}

func run(conf Config, log *internal.Logger) error {
	if !conf.IsValid() {
		return fmt.Errorf("Invalid Config: %+v", conf)
	}

	repo, err := pgRepo(conf, log)
	if err != nil {
		return err
	}
	// repo := fakeRepo()

	server := rest.NewServer(repo, log)

	log.Info("Listening on " + conf.Port)
	fmt.Println("Listening on " + conf.Port)
	return http.ListenAndServe(":"+conf.Port, server)
}

func pgRepo(conf Config, log *internal.Logger) (usecase.BookReaderWriter, error) {

	if err := dbMigrateUp(conf.DSN, log); err != nil {
		return nil, err
	}

	db, err := sql.Open(dbDriver, conf.DSN)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to postgres: %s", err.Error())
	}

	repo := repository.NewPostgresRepo(db)
	return repo, nil
}

func fakeRepo() usecase.BookReaderWriter {
	return fake.NewBookRepo()
}
