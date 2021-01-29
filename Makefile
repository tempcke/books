project = $(shell basename $(shell pwd))
MODULE = $(shell go list -m)

build:
	go build -o bin/bookserver cmd/bookserver/*.go

test:
	go test -coverprofile /tmp/$(project)-test-coverage ./...

run: build
	bin/rpmserver

clean:
	docker container rm books-db books-server && docker volume rm $(project)_books-data

migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)