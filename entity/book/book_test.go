package book_test

import (
	"testing"
	"time"

	"github.com/tempcke/books/entity/book"
)

// dateFormat is the format I would expect a date value to be passed as a string
// the production code doesn't require this format
const dateFormat = "2006-01-02"

const (
	title   = "Refactoring"
	author  = "Martin Fowler"
	pubdate = "1999-06-28"
	rating  = book.RateThree
	status  = book.StatusCheckedIn
)

func TestBook(t *testing.T) {
	pubDate, err := time.Parse(dateFormat, pubdate)
	if err != nil {
		t.Error(err)
	}
	b := book.NewBook(title, author, pubDate, rating, status)
	assertEqual(t, title, b.Title)
	assertEqual(t, author, b.Author)
	assertEqual(t, pubDate, b.PubDate)
	assertEqual(t, rating, b.Rating)
	assertEqual(t, status, b.Status)
	assertEqual(t, nil, b.Validate())
}

func TestBookValidation(t *testing.T) {
	pubDate, err := time.Parse(dateFormat, pubdate)
	if err != nil {
		t.Error(err)
	}

	t.Run("Empty Title", func(t *testing.T) {
		b := book.NewBook("", author, pubDate, rating, status)
		assertEqual(t, book.ErrTitleIsRequired, b.Validate())
	})

	t.Run("Empty Author", func(t *testing.T) {
		b := book.NewBook(title, "", pubDate, rating, status)
		assertEqual(t, book.ErrAuthorIsRequired, b.Validate())
	})

	t.Run("Zero PubDate", func(t *testing.T) {
		b := book.NewBook(title, author, time.Time{}, rating, status)
		assertEqual(t, book.ErrPubDateIsRequired, b.Validate())
	})

	t.Run("Invalid Rating", func(t *testing.T) {
		b := book.NewBook(title, author, pubDate, 42, status)
		assertEqual(t, book.ErrRatingInvalid, b.Validate())
	})

	t.Run("Invalid Status", func(t *testing.T) {
		b := book.NewBook(title, author, pubDate, rating, "SomeInvalidStatus")
		assertEqual(t, book.ErrStatusInvalid, b.Validate())
	})
}

func TestBookID(t *testing.T) {
	b1 := book.NewBook("first book", "john smith", time.Now(), rating, status)
	b2 := book.NewBook("second one", "john smith", time.Now(), rating, status)
	// book id should be a uuid string, which is 36 chars long
	assertEqual(t, 36, len(b1.ID))
	assertEqual(t, 36, len(b2.ID))
	assertNotEqual(t, b1.ID, b2.ID)
}

func assertEqual(t *testing.T, want, got interface{}) {
	t.Helper()
	if got != want {
		t.Errorf(
			"Not Equal!\nWant: %v\t%T\nGot:  %v\t%T",
			want, want,
			got, got)
	}
}

func assertNotEqual(t *testing.T, a, b interface{}) {
	t.Helper()
	if a == b {
		t.Errorf("Should not be equal!\nA: %v\t%T\nB: %v\t%T", a, a, b, b)
	}
}
