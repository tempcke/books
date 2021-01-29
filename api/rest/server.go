package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tempcke/books/internal"
	"github.com/tempcke/books/usecase"
)

// Server is used to expose appliaction over a RESTful API
type Server struct {
	http.Handler
	bookRepo usecase.BookReaderWriter
	log      *internal.Logger
}

// NewServer constructs a Server
func NewServer(bookRepo usecase.BookReaderWriter, logger *internal.Logger) *Server {
	server := new(Server)
	server.bookRepo = bookRepo
	server.log = logger
	server.initRouter()
	return server
}

func (s *Server) initRouter() {
	r := chi.NewRouter()
	r.Route("/book", func(r chi.Router) {
		r.Post("/", addBook(s.bookRepo, s.log))
		r.Get("/", listBooks(s.bookRepo, s.log))
		r.Route("/{bookID}", func(r chi.Router) {
			r.Get("/", getBook(s.bookRepo, s.log))
			r.Delete("/", deleteBook(s.bookRepo, s.log))
			r.Put("/status/{status}", putBookStatus(s.bookRepo, s.log))
			r.Put("/rating/{rating}", putBookRating(s.bookRepo, s.log))
		})
	})
	s.Handler = r
}
