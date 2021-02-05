project = $(shell basename $(shell pwd))
MODULE = $(shell go list -m)

build: test
	go build -o bin/bookserver cmd/bookserver/*.go

test: .env
	go test -coverprofile /tmp/$(project)-test-coverage ./...

run: build .env
	docker-compose up

clean:
	docker container rm books-db books-server && docker volume rm $(project)_books-data

migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)

.env:
	cp .env.example .env
