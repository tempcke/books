// Package fake contains fakes (mocks) used for testing only
package fake

import (
	"errors"

	"github.com/tempcke/books/entity/book"
)

// BookRepo is a fake book repository
type BookRepo struct {
	books map[string]book.Book
}

// NewBookRepo creates and returns a BookRepo
func NewBookRepo() BookRepo {
	return BookRepo{make(map[string]book.Book)}
}

// AddBook adds a book
func (r BookRepo) AddBook(book book.Book) error {
	r.books[book.ID] = book
	return nil
}

// RemoveBook removes a book
func (r BookRepo) RemoveBook(id string) error {
	if _, ok := r.books[id]; !ok {
		return errors.New("book not found")
	}
	delete(r.books, id)
	return nil
}

// GetBookByID gets a book by id
func (r BookRepo) GetBookByID(id string) (book.Book, error) {
	book, ok := r.books[id]
	if !ok {
		return book, errors.New("book not found")
	}
	return book, nil
}

// BookList lists books
func (r BookRepo) BookList() ([]book.Book, error) {
	list := make([]book.Book, len(r.books))
	i := 0
	for _, b := range r.books {
		list[i] = b
		i++
	}
	return list, nil
}

// UpdateBook updates a book record
func (r BookRepo) UpdateBook(book book.Book) error {
	r.books[book.ID] = book
	return nil
}
