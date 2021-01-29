package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/tempcke/books/api/rest"
	"github.com/tempcke/books/entity/book"
	"github.com/tempcke/books/fake"
)

type jsonMap map[string]interface{}

var (
	logger = &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			TimestampFormat: "15:04:05",
		},
		Hooks: make(logrus.LevelHooks),
	}
	repo   = fake.NewBookRepo()
	server = rest.NewServer(repo, logger)
)

var (
	author  = "john smith"
	pubdate = "2020-01-01"
	rating  = book.RateOne
	status  = book.StatusCheckedIn
)

var bookJsonTemplate = `{"title":"%v","author":"%v","pubdate":"%v","rating":%v,"status":"%v"}`

var dateFormat = "2006-01-02"

func makeBookJson(title string) string {
	return fmt.Sprintf(
		bookJsonTemplate,
		title,
		author,
		pubdate,
		rating,
		status,
	)
}

func TestPostBook(t *testing.T) {
	t.Run("expect 200 and book stored in repo", func(t *testing.T) {
		title := "post book"

		// post new book
		json := makeBookJson(title)
		rr := httptestPost("/book", json) // responseRecorder

		// parse response
		data := getJsonMapFromResponseBody(t, rr)
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.NotEmpty(t, data["id"])
		id := data["id"].(string)

		// confirm record in repo
		b, err := repo.GetBookByID(id)
		assert.NoError(t, err)
		assert.Equal(t, title, b.Title)
		assert.Equal(t, author, b.Author)
		assert.Equal(t, pubdate, b.PubDate.Format(dateFormat))
		assert.Equal(t, rating, b.Rating)
		assert.Equal(t, status, b.Status)

		// check response data structure
		assertDataMatchesBook(t, data, b)
	})

	t.Run("post book with empty title, expect 400", func(t *testing.T) {
		json := makeBookJson("")
		rr := httptestPost("/book", json)
		data := getJsonMapFromResponseBody(t, rr)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.NotEmpty(t, data["error"])
	})

	t.Run("post with invalid json", func(t *testing.T) {
		rr := httptestPost("/book", `{"title":"t","author":"a","pubdate":"2020-01-01","rating":1,"status":"CheckedIn",}`) // trailing comma is invalid
		data := getJsonMapFromResponseBody(t, rr)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.NotEmpty(t, data["error"])
	})

	t.Run("post with invalid pubdate format", func(t *testing.T) {
		json := fmt.Sprintf(
			bookJsonTemplate,
			"invalid pubdate",
			author,
			"01/01/2020",
			rating,
			status,
		)
		rr := httptestPost("/book", json)
		data := getJsonMapFromResponseBody(t, rr)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.NotEmpty(t, data["error"])
	})
}

func TestGetBook(t *testing.T) {
	t.Run("sunny day", func(t *testing.T) {
		b := makeBook("GET book")
		repo.AddBook(b)
		rr := httptestGet("/book/" + b.ID)
		data := getJsonMapFromResponseBody(t, rr)
		assert.Equal(t, http.StatusOK, rr.Code)
		assertDataMatchesBook(t, data, b)
	})

	t.Run("book not found", func(t *testing.T) {
		b := makeBook("GET book not found")
		rr := httptestGet("/book/" + b.ID)
		data := getJsonMapFromResponseBody(t, rr)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.NotEmpty(t, data["error"])
	})
}

func TestListBooks(t *testing.T) {
	t.Run("expect empty set when no books exist", func(t *testing.T) {
		// reset repo and server
		repo = fake.NewBookRepo()
		server = rest.NewServer(repo, logger)

		getResponse := httptestGet("/book")
		assert.Equal(t, http.StatusOK, getResponse.Code)

		var bookList struct {
			Items []map[string]interface{} `json:"items"`
		}
		err := json.Unmarshal(getResponse.Body.Bytes(), &bookList)
		assert.Nil(t, err)
		assert.Len(t, bookList.Items, 0)
	})

	t.Run("list two books", func(t *testing.T) {
		// create two books in repo
		b1 := makeBook("list book 1")
		b2 := makeBook("list book 2")
		repo.AddBook(b1)
		repo.AddBook(b2)

		// list via API
		getResponse := httptestGet("/book")
		assert.Equal(t, http.StatusOK, getResponse.Code)

		var bookList struct {
			Items []map[string]interface{} `json:"items"`
		}
		err := json.Unmarshal(getResponse.Body.Bytes(), &bookList)
		assert.Nil(t, err)

		// ensure two results
		assert.Len(t, bookList.Items, 2)

		// ensure those results are among the books added
		for _, b := range bookList.Items {
			id := b["id"]
			if id != b1.ID && id != b2.ID {
				t.Errorf("book id %v was not added?", id)
				break
			}
		}

		// make sure the same book wasn't just listed twice...
		assert.NotEqual(t, bookList.Items[0]["id"], bookList.Items[1]["id"])
	})
}

func TestDeleteBook(t *testing.T) {
	b := makeBook("del book")

	// should a restful DELETE on a resource that does not exist
	// result in a 404 or not?
	// https://stackoverflow.com/a/16632048/2683059
	// a lot of conflicting answers on this one, I'm going to chose no
	// for now because I can't think of a reason why the client should care
	t.Run("unknown book, expect 204", func(t *testing.T) {
		rr := httptestDelete("/book/" + b.ID)
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("delete existing book, expect 204", func(t *testing.T) {
		repo.AddBook(b)

		rr := httptestDelete("/book/" + b.ID)
		assert.Equal(t, http.StatusNoContent, rr.Code)

		// record should not be retrievable by repo anymore
		_, err := repo.GetBookByID(b.ID)
		assert.Error(t, err)
	})
}

// http request helper functions
func httptestPost(uri, jsonStr string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodPost, uri, jsonReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	return execReq(req)
}

func httptestGet(uri string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	return execReq(req)
}

func httptestDelete(uri string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodDelete, uri, nil)
	return execReq(req)
}

func execReq(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)
	return rr
}

// json helper functions
func jsonReader(jsonStr string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(jsonStr))
}

func getJsonMapFromResponseBody(t *testing.T, r *httptest.ResponseRecorder) jsonMap {
	t.Helper()
	var m jsonMap
	err := json.Unmarshal(r.Body.Bytes(), &m)
	assert.Nil(t, err)
	return m
}

func makeBook(title string) book.Book {
	return book.NewBook(title, "john smith", time.Now(), book.RateOne, book.StatusCheckedIn)
}

// custom assertions
func assertDataMatchesBook(t *testing.T, data jsonMap, b book.Book) {
	t.Helper()
	assert.Equal(t, b.ID, data["id"])
	assert.Equal(t, b.Title, data["title"])
	assert.Equal(t, b.Author, data["author"])
	assert.Equal(t, b.PubDate.Format(dateFormat), data["pubdate"])
	assert.Equal(t, float64(b.Rating), data["rating"])
	assert.Equal(t, b.Status.String(), data["status"])
}
