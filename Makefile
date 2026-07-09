.PHONY: dev up down reset install run test test-race cover test-integration test-all

dev:
	docker compose up -d mongodb redis
	go run ./cmd/lnkd/

up:
	docker compose up -d --build

down:
	docker compose down

reset:
	docker compose down -v

run:
	go run ./cmd/lnk/

install:
	go build -o $(shell go env GOPATH)/bin/lnk ./cmd/lnk/

test:
	go test ./...

test-race:
	go test ./... -race

cover:
	go test ./... -cover

test-integration:
	go test -tags integration ./internal/repository/...

test-all:
	go test -tags integration ./...
