package repository_test

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
	"github.com/tempcke/books/entity/book"
	"github.com/tempcke/books/repository"
	"github.com/tempcke/books/usecase"
)

const dateFormat = "2006-01-02"

var pgRepo repository.Postgres

func TestPostgresRepo(t *testing.T) {
	r := pgRepo

	t.Run("ensure PostgresRepository is a BookReadWriter", func(t *testing.T) {
		assert.Implements(t, (*usecase.BookReaderWriter)(nil), pgRepo)
	})

	t.Run("GetBookByID should return error when book not found", func(t *testing.T) {
		b := makeBook("non existing book")
		_, err := r.GetBookByID(b.ID)
		assert.Error(t, err)
	})

	t.Run("add and get book", func(t *testing.T) {
		b := makeBook("add book")

		if err := r.AddBook(b); err != nil {
			t.Fatal(err)
		}

		bOut, err := r.GetBookByID(b.ID)
		assert.NoError(t, err)
		assert.Equal(t, b.ID, bOut.ID)
		assert.Equal(t, b.Title, bOut.Title)
		assert.Equal(t, b.Author, bOut.Author)
		assert.Equal(t, b.PubDate.Format(dateFormat), bOut.PubDate.Format(dateFormat))
		assert.Equal(t, b.Rating, bOut.Rating)
		assert.Equal(t, b.Status, bOut.Status)
	})

	t.Run("list books", func(t *testing.T) {
		// create books
		a := makeBook("list book A")
		b := makeBook("list book B")
		c := makeBook("list book C")

		books := map[string]book.Book{
			a.ID: a,
			b.ID: b,
			c.ID: c,
		}

		// store books
		r.AddBook(a)
		r.AddBook(b)
		r.AddBook(c)

		// list entities, this is what we want to test!
		bookList, err := r.BookList()
		assert.NoError(t, err)

		// iterate over list counting the times each id is seen
		seen := make(map[string]int, 3)
		for _, bk := range bookList {
			delete(books, bk.ID)
			seen[bk.ID]++
		}

		// ensure the list does not repeat any books
		for _, n := range seen {
			assert.Equal(t, 1, n, "book ids are repeated in the list")
		}

		// books are removed from map as they are found so there should be none left
		assert.Len(t, books, 0)
	})

	t.Run("remove book", func(t *testing.T) {
		// create and store book
		b := makeBook("remove book")

		t.Run("delete a book that does not exist should error", func(t *testing.T) {
			err := r.RemoveBook(b.ID)
			assert.Equal(t, repository.ErrRecordNotFound, err)
		})

		t.Run("add then remove book", func(t *testing.T) {
			r.AddBook(b)

			// remove book
			if err := r.RemoveBook(b.ID); err != nil {
				t.Fatal(err)
			}

			// try to retrieve Book
			if _, err := r.GetBookByID(b.ID); err == nil {
				t.Fatal("book found when it should have been deleted")
			}
		})
	})

	t.Run("update book", func(t *testing.T) {
		b := makeBook("update book")

		t.Run("can not update book that does not exist", func(t *testing.T) {
			err := r.UpdateBook(b)
			assert.Equal(t, repository.ErrRecordNotFound, err)
		})

		t.Run("add then update book", func(t *testing.T) {
			r.AddBook(b)
			b.Status = book.StatusCheckedOut
			b.Rating = book.RateTwo
			err := r.UpdateBook(b)
			assert.NoError(t, err)
		})
	})
}

func makeBook(title string) book.Book {
	return book.NewBook(title, "john smith", time.Now(), book.RateOne, book.StatusCheckedIn)
}

func loadMigrations(db *sql.DB) error {
	query := `
	  CREATE TABLE IF NOT EXISTS books (
			id         VARCHAR(36)  PRIMARY KEY,
			title      VARCHAR(128),
			author     VARCHAR(128),
			pubdate    date,
			rating     INT,
			status     VARCHAR(16),
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ
		)
	`
	_, err := db.Exec(query)
	return err
}

// this setup uses dockertest to create a postgres instance in docker for testing
// it also constructs the pgRepo var used by all the tests
// this way all the tests can test aginst a real postgres instance
func TestMain(m *testing.M) {
	// Setup docker postgres instance
	rand.Seed(time.Now().UnixNano())
	var (
		db *sql.DB

		user     = "postgres"
		password = "secret"
		dbname   = "postgres"
		port     = strconv.Itoa(54300 + rand.Intn(100))
		dialect  = "postgres"
		dsn      = fmt.Sprintf(
			"postgres://%s:%s@localhost:%v/%s?sslmode=disable",
			user, password, port, dbname,
		)

		containerTTL uint = 60 // seconds
	)

	pool, err := dockertest.NewPool("")

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbname,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	resource.Expire(containerTTL)
	defer resource.Close()

	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}

	if err = pool.Retry(func() error {
		db, err = sql.Open(dialect, dsn)
		if err != nil {
			return err
		}
		if err = db.Ping(); err != nil {
			return err
		}
		if err = loadMigrations(db); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	// deferred cleanup
	defer func() {
		db.Close()

		// cleanup
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	// construct repository
	pgRepo = repository.NewPostgresRepo(db)

	// run the tests
	os.Exit(m.Run())
}
