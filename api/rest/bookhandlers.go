package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/tempcke/books/entity/book"
	"github.com/tempcke/books/internal"
	"github.com/tempcke/books/usecase"
)

var dateFormat = "2006-01-02"

func addBook(bookRepo usecase.BookReaderWriter, log *internal.Logger) http.HandlerFunc {
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

func getBook(bookRepo usecase.BookReader, log *internal.Logger) http.HandlerFunc {
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

func listBooks(bookRepo usecase.BookReader, log *internal.Logger) http.HandlerFunc {
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

func deleteBook(bookRepo usecase.BookReaderWriter, log *internal.Logger) http.HandlerFunc {
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

func putBookStatus(bookRepo usecase.BookReaderWriter, log *internal.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")
		b, err := usecase.GetBook(bookRepo, bookID)
		if err != nil {
			errorResponse(w, http.StatusNotFound, "bookId not found")
			log.Debug("putBookStatus handler, id not found: " + bookID)
			return
		}

		status := chi.URLParam(r, "status")

		b, err = usecase.ChangeBookStatus(bookRepo, bookID, book.Status(status))
		if err != nil {
			log.Debug(err)
			errorResponse(w, http.StatusBadRequest, "Failed to update book, are you passing a valid status?")
			return
		}

		w.WriteHeader(http.StatusOK)
		jsonResponse(w, NewBookModel(b))
	}
}

func putBookRating(bookRepo usecase.BookReaderWriter, log *internal.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")
		b, err := usecase.GetBook(bookRepo, bookID)
		if err != nil {
			errorResponse(w, http.StatusNotFound, "bookId not found")
			log.Debug("putBookStatus handler, id not found: " + bookID)
			return
		}

		rating := chi.URLParam(r, "rating")
		value, err := strconv.Atoi(rating)
		if err != nil {
			log.Debug("putBookStatus handler, could not convert rating to int: " + rating)
			errorResponse(w, http.StatusBadRequest, "Invalid rating, could not convert to int")
		}

		b, err = usecase.ChangeBookRating(bookRepo, bookID, book.Rating(value))
		if err != nil {
			log.Debug(err)
			errorResponse(w, http.StatusBadRequest, "Failed to update book, are you passing a valid rating?")
			return
		}

		w.WriteHeader(http.StatusOK)
		jsonResponse(w, NewBookModel(b))
	}
}
