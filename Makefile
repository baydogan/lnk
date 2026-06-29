.PHONY: dev up down install run

# Local code against containerized deps (fast iteration loop).
dev:
	docker compose up -d mongodb redis
	go run ./cmd/lnkd/

# Full stack in containers, rebuilding the lnkd image.
up:
	docker compose up -d --build

down:
	docker compose down

run:
	go run ./cmd/lnk/

install:
	go build -o $(shell go env GOPATH)/bin/lnk ./cmd/lnk/
