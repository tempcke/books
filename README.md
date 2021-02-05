# BookServer Challenge

## Challenge Requirments
Golang/REST Code Challenge

Using Go as your language, create a CRUD API to manage a list of Books, fulfilling the

### following requirements:
1. Books should have the following Attributes:
  - Title
  - Author
  - Publisher
  - Publish Date
  - Rating (1-3)
  - Status (CheckedIn, CheckedOut)
2. Each endpoint should have test coverage of both successful and failed (due to user error) requests.
3. Use a data store of your choice.
4. Include unit tests that exercise the endpoints such that both 200-level and 400-level responses are induced.
5. The app should be stood up in Docker, and client code (such as cURL requests and your unit tests) should be executed on the host machine, against the containerized app.
6. Send the project along as a git repository.
7. Please do not use go-swagger to generate the server side code. Part of the goal of this challenge is to see your coding skills. :)

-----
## Notes on decisions made

### database for testing
If production is going to hit a real postgres instance then I want the tests to hit a postgres instance.  There is not a reliable in-memory substutue to test postgreSQL queries.  Therefore I'm using the dockertest library which results in a 2 to 5 second lag time for the test as it spins up the container, but it is worth it.  Sometimes I use build tags to only run those integration tests on travis or circle etc so they do not slow down my normal test runs during development.

### database update queries
I personally do not like mutating db objects if it can be avoided.  For this reason if I had more time I would remove status and rating from the book table and instead keep the values in another table with a timestamp, the most recent one would be pulled with the book entity however the change history would be retained in the data store and mutations of the book row would not be needed.

## Setup and execution instructions

### Environment Variables
.env is used automaticly by docker-compose to pass environment variables into the containers however it is ignored by git via .gitignore so an example file has been provided in .env.example

This example file will be copied to .env when you issue `make .env` or `make run`

For development you can simply leave .env alone unless you desire the appliaction expose itself on different ports to avoid local port conflicts.

### Run the server
`make run` will ensure .env exists and then do `docker-compose up`  so long as you have docker and docker-composed install everything *should* work just fine.  Please create an issue letting me know if something does not work as expected

## RESTful API requests
### Add Book
```
curl -X POST "http://localhost:8080/book" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d '{
  "title": "Refactoring",
  "author": "Martin Fowler",
  "pubdate": "1999-06-28",
  "rating": 3,
  "status": "CheckedIn"
}' | json_pp
```

### Change Status
```
curl -X PUT "http://localhost:8080/book/{bookId}/status/{status}" \
     -H 'Accept: application/json' | json_pp
```

### Change Rating
```
curl -X PUT "http://localhost:8080/book/{bookId}/rating/{rating}" \
     -H 'Accept: application/json' | json_pp
```

### List Books
```
curl -X GET "http://localhost:8080/book" \
     -H 'Accept: application/json' | json_pp
```

### Get Book Detail
```
curl -X GET "http://localhost:8080/book/{bookId}" \
     -H 'Accept: application/json' | json_pp
```

### Delete Book
```
curl -X DELETE "http://localhost:8080/book/{bookId}"
```