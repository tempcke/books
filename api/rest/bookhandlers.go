package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/tempcke/books/entity/book"
	"github.com/tempcke/books/usecase"
)

var dateFormat = "2006-01-02"

func addBook(bookRepo usecase.BookReaderWriter, log *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := BookModel{}
		if err := decodeRequestData(w, r.Body, &data); err != nil {
			log.Error(err)
			return
		}

		pDate, err := time.Parse(dateFormat, data.PubDate)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, "pubdate must be in yyyy-mm-dd format")
			log.Error(err)
			return
		}

		b := book.NewBook(
			data.Title,
			data.Author,
			pDate,
			book.Rating(data.Rating),
			book.Status(data.Status),
		)

		if err := usecase.AddBook(bookRepo, b); err != nil {
			log.Debug(err)
			errorResponse(w, http.StatusBadRequest, "Missing or invalid fields")
			return
		}

		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, NewBookModel(b))
	}
}

func getBook(bookRepo usecase.BookReader, log *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")
		b, err := usecase.GetBook(bookRepo, bookID)
		if err != nil {
			errorResponse(w, http.StatusNotFound, "bookId not found")
			log.Debug("getBook handler, id not found: " + bookID)
			return
		}
		jsonResponse(w, NewBookModel(b))
	}
}

func listBooks(bookRepo usecase.BookReader, log *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := usecase.ListBooks(bookRepo)
		if err != nil {
			log.Error(err)
			errorResponse(w, http.StatusInternalServerError, "Error fetching list")
			return
		}
		jsonResponse(w, NewBookListModel(books...))
	}
}

func deleteBook(bookRepo usecase.BookReaderWriter, log *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")
		err := usecase.RemoveBook(bookRepo, bookID)
		if err != nil {
			// what should a RESTful DELETE endpoint do
			// when the resource does not exist?
			// for now I vote nothing, the client wants it gone
			// and it isn't there ... so client should be happy

			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
