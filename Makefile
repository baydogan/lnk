.PHONY: dev up down reset install run

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
