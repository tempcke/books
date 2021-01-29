package rest

import (
	"github.com/tempcke/books/entity/book"
)

// ErrorResponse response model
type ErrorResponse struct {
	Error string `json:"error"`
}

// BookList response model
type BookList struct {
	Items []BookModel `json:"items"`
}

// NewBookListModel constructs a BookList model from a set of books
func NewBookListModel(bookList ...book.Book) BookList {
	pl := BookList{
		Items: make([]BookModel, len(bookList)),
	}
	for i, p := range bookList {
		pl.Items[i] = NewBookModel(p)
	}
	return pl
}

// BookModel is a response model for a book
type BookModel struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	PubDate string `json:"pubdate"`
	Rating  int    `json:"rating"`
	Status  string `json:"status"`
}

// NewBookModel is the BookModel constructor
func NewBookModel(book book.Book) BookModel {
	return BookModel{
		ID:      book.ID,
		Title:   book.Title,
		Author:  book.Author,
		PubDate: book.PubDate.Format(dateFormat),
		Rating:  int(book.Rating),
		Status:  string(book.Status),
	}
}
