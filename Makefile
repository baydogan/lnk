.PHONY: dev install run

dev:
	docker compose up -d
	go run ./cmd/lnkd/

run:
	go run ./cmd/lnk/

install:
	go build -o $(shell go env GOPATH)/bin/lnk ./cmd/lnk/
	go build -o $(shell go env GOPATH)/bin/lnkd ./cmd/lnkd/