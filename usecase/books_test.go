package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tempcke/books/entity/book"
	"github.com/tempcke/books/usecase"
)

func TestAddBook(t *testing.T) {
	repo := newBookRepo()
	goodBook := makeBook("add book")
	badBook := makeBook("") // empty title will not validate
	assert.NoError(t, usecase.AddBook(repo, goodBook))
	assert.Error(t, usecase.AddBook(repo, badBook))
}

func TestGetBook(t *testing.T) {
	repo := newBookRepo()
	bIn := makeBook("get book")

	t.Run("expect error if book not found", func(t *testing.T) {
		_, err := usecase.GetBook(repo, bIn.ID)
		assert.Error(t, err)
	})

	t.Run("get a book that exists", func(t *testing.T) {
		repo.AddBook(bIn)
		bOut, err := usecase.GetBook(repo, bIn.ID)
		assert.NoError(t, err)
		assert.Equal(t, bIn, bOut)
	})
}

func TestRemoveBook(t *testing.T) {
	repo := newBookRepo()
	b := makeBook("remove book")

	// this is not a great test as the error is currently coming from the repo
	// will have to be sure to test this on the real repo itself
	t.Run("expect error if book not found", func(t *testing.T) {
		err := usecase.RemoveBook(repo, b.ID)
		assert.Error(t, err)
	})

	t.Run("remove an existing book", func(t *testing.T) {
		repo.AddBook(b)
		err := usecase.RemoveBook(repo, b.ID)
		assert.NoError(t, err)
		_, err = repo.GetBookByID(b.ID)
		assert.Error(t, err)
	})
}

func TestListBooks(t *testing.T) {
	repo := newBookRepo()
	a, b := makeBook("A"), makeBook("B")
	repo.AddBook(a)
	repo.AddBook(b)
	books, err := usecase.ListBooks(repo)

	// error should only happen on a db connection or query error
	// we are using a fake repo so it won't happen but check it anyway?
	assert.NoError(t, err)

	// we added 2, so there should be 2
	assert.Len(t, books, 2)

	// ensure we didnt get the same book repeated twice
	assert.NotEqual(t, books[0].ID, books[1].ID)

	// ensure we got back the books we added
	for _, book := range books {
		assert.True(t, book.ID == a.ID || book.ID == b.ID)
	}
}

func TestUpdateBookStatus(t *testing.T) {
	repo := newBookRepo()
	a := makeBook("update status")

	t.Run("expect error when book does not exist", func(t *testing.T) {
		_, err := usecase.ChangeBookStatus(repo, a.ID, book.StatusCheckedOut)
		assert.Error(t, err)
	})

	repo.AddBook(a)

	t.Run("expect error on invalid status", func(t *testing.T) {
		_, err := usecase.ChangeBookStatus(repo, a.ID, "invalid-status")
		assert.Error(t, err)
	})

	t.Run("should update the status", func(t *testing.T) {
		b, err := usecase.ChangeBookStatus(repo, a.ID, book.StatusCheckedOut)
		assert.NoError(t, err)
		assert.Equal(t, book.StatusCheckedOut, b.Status)
	})
}

func TestUpdateBookRating(t *testing.T) {
	repo := newBookRepo()
	a := makeBook("update rating")

	t.Run("expect error when book does not exist", func(t *testing.T) {
		_, err := usecase.ChangeBookRating(repo, a.ID, book.RateTwo)
		assert.Error(t, err)
	})

	repo.AddBook(a)

	t.Run("expect error on invalid rating", func(t *testing.T) {
		_, err := usecase.ChangeBookRating(repo, a.ID, 42)
		assert.Error(t, err)
	})

	t.Run("should update the rating", func(t *testing.T) {
		b, err := usecase.ChangeBookRating(repo, a.ID, book.RateTwo)
		assert.NoError(t, err)
		assert.Equal(t, book.RateTwo, b.Rating)
	})
}

func makeBook(title string) book.Book {
	return book.NewBook(title, "john smith", time.Now(), book.RateOne, book.StatusCheckedIn)
}

type bookRepo struct {
	books map[string]book.Book
}

func newBookRepo() bookRepo {
	return bookRepo{make(map[string]book.Book)}
}

func (r bookRepo) AddBook(book book.Book) error {
	r.books[book.ID] = book
	return nil
}
func (r bookRepo) RemoveBook(id string) error {
	if _, ok := r.books[id]; !ok {
		return errors.New("book not found")
	}
	delete(r.books, id)
	return nil
}
func (r bookRepo) GetBookByID(id string) (book.Book, error) {
	book, ok := r.books[id]
	if !ok {
		return book, errors.New("book not found")
	}
	return book, nil
}
func (r bookRepo) BookList() ([]book.Book, error) {
	list := make([]book.Book, len(r.books))
	i := 0
	for _, b := range r.books {
		list[i] = b
		i++
	}
	return list, nil
}
func (r bookRepo) UpdateBook(book book.Book) error {
	r.books[book.ID] = book
	return nil
}
