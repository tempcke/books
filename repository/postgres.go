package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/tempcke/books/entity/book"
)

// Errors
var (
	ErrRecordNotFound  = errors.New("Record not found")
	ErrRecordNotUnique = errors.New("Record not unique")
)

// Postgres repository should NOT be used in production
type Postgres struct {
	db *sql.DB
}

// NewPostgresRepo constructs an Postgres repository
func NewPostgresRepo(db *sql.DB) Postgres {
	if err := db.Ping(); err != nil {
		log.Fatal("Could not connect to db: " + err.Error())
	}

	return Postgres{
		db: db,
	}
}

// AddBook persists a book
func (r Postgres) AddBook(b book.Book) error {
	// custom error in case record already exists
	if _, err := r.GetBookByID(b.ID); err == nil {
		return ErrRecordNotUnique
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO books
		(id, title, author, pubdate, rating, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		b.ID,
		b.Title,
		b.Author,
		b.PubDate,
		b.Rating,
		b.Status,
		time.Now(),
		time.Now(),
	)

	return err
}

// RemoveBook removes a previously stored book
func (r Postgres) RemoveBook(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM books WHERE id = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)

	if err != nil {
		return err
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// GetBookByID returns a previously stored book
func (r Postgres) GetBookByID(id string) (b book.Book, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, author, pubdate, rating, status
		FROM books WHERE id = $1
	`

	err = r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.Title, &b.Author,
		&b.PubDate, &b.Rating, &b.Status,
	)

	return b, err
}

// BookList returns all books previously stored
// idealy this would be filterable at some point...
func (r Postgres) BookList() ([]book.Book, error) {
	bookList := make([]book.Book, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, author, pubdate, rating, status
		FROM books
	`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return bookList, err
	}
	defer rows.Close()

	for rows.Next() {
		b := book.Book{}

		err = rows.Scan(
			&b.ID, &b.Title, &b.Author,
			&b.PubDate, &b.Rating, &b.Status,
		)
		if err != nil {
			return bookList, err
		}

		bookList = append(bookList, b)
	}

	return bookList, nil
}

// UpdateBook updates a previously stored book record
func (r Postgres) UpdateBook(b book.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE books
		SET rating = $2,
				status = $3,
				updated_at = $4
		WHERE id = $1;
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		b.ID,
		b.Rating,
		b.Status,
		time.Now(),
	)

	if err != nil {
		return err
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return ErrRecordNotFound
	}

	return nil
}
