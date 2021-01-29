package book

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Validation Errors
var (
	ErrTitleIsRequired   = errors.New("Title is required")
	ErrAuthorIsRequired  = errors.New("Author is required")
	ErrPubDateIsRequired = errors.New("PubDate is required")
	ErrRatingInvalid     = errors.New("Rating value is not supported")
	ErrStatusInvalid     = errors.New("Status value is not supported")
)

// Rating is a book rating from 1 to 3
type Rating int

// Int returns the int type
func (r Rating) Int() int {
	return int(r)
}
func (r Rating) String() string {
	return strconv.Itoa(r.Int())
}

// Rating values
const (
	RateOne = Rating(iota + 1)
	RateTwo
	RateThree
)

// Status is the status of the book, checked in or out
type Status string

func (s Status) String() string {
	return string(s)
}

// Status values
const (
	StatusCheckedIn  = Status("CheckedIn")
	StatusCheckedOut = Status("CheckedOut")
)

// Book entity
type Book struct {
	ID      string
	Title   string
	Author  string
	PubDate time.Time
	Rating  Rating
	Status  Status
}

// NewBook creates a new Book
func NewBook(
	title, author string,
	pubdate time.Time,
	rating Rating,
	status Status,
) Book {
	return Book{
		ID:      uuid.New().String(),
		Title:   title,
		Author:  author,
		PubDate: pubdate,
		Rating:  rating,
		Status:  status,
	}
}

// Validate the Book object
func (b Book) Validate() error {
	var zeroTime time.Time
	if len(b.Title) == 0 {
		return ErrTitleIsRequired
	}
	if len(b.Author) == 0 {
		return ErrAuthorIsRequired
	}
	if b.PubDate == zeroTime {
		return ErrPubDateIsRequired
	}
	if err := b.validateRating(); err != nil {
		return err
	}
	return b.validateStatus()
}

func (b Book) validateRating() error {
	switch b.Rating {
	case RateOne, RateTwo, RateThree:
		return nil
	}
	return ErrRatingInvalid
}

func (b Book) validateStatus() error {
	switch b.Status {
	case StatusCheckedIn, StatusCheckedOut:
		return nil
	}
	return ErrStatusInvalid
}
