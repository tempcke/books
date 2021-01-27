package usecase

import "github.com/tempcke/books/entity/book"

// BookReader is used to fetch information about books
type BookReader interface {
	GetBookByID(id string) (book.Book, error)
	BookList() ([]book.Book, error)
}

// BookWriter is used to add and remove books
type BookWriter interface {
	AddBook(book.Book) error
	RemoveBook(id string) error
	UpdateBook(book.Book) error
}

// BookReaderWriter is used for updates
type BookReaderWriter interface {
	BookReader
	BookWriter
}

// AddBook is used to store a book
func AddBook(r BookWriter, book book.Book) error {
	if err := book.Validate(); err != nil {
		return err
	}
	return r.AddBook(book)
}

// GetBook gets a book by id
func GetBook(r BookReader, id string) (book.Book, error) {
	return r.GetBookByID(id)
}

// ListBooks lists all books from storage
func ListBooks(r BookReader) ([]book.Book, error) {
	return r.BookList()
}

// RemoveBook removes a book, error if does not exist or storage fails
func RemoveBook(r BookWriter, id string) error {
	return r.RemoveBook(id)
}

// ChangeBookStatus is used to modify the status of a book
func ChangeBookStatus(r BookReaderWriter, id string, status book.Status) (book.Book, error) {
	book, err := r.GetBookByID(id)
	if err != nil {
		return book, err
	}

	book.Status = status
	if err := book.Validate(); err != nil {
		return book, err
	}

	if err := r.UpdateBook(book); err != nil {
		return book, err
	}

	return book, nil
}

// ChangeBookRating is used to modify the rating of a book
func ChangeBookRating(r BookReaderWriter, id string, rating book.Rating) (book.Book, error) {
	book, err := r.GetBookByID(id)
	if err != nil {
		return book, err
	}

	book.Rating = rating
	if err := book.Validate(); err != nil {
		return book, err
	}

	if err := r.UpdateBook(book); err != nil {
		return book, err
	}

	return book, nil
}
