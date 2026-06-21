.PHONY: dev install run

dev:
	docker compose up -d
	go run ./cmd/linkd/

run:
	go run ./cmd/linkd/

install:
	go build -o $(shell go env GOPATH)/bin/lnk ./cmd/lnk/