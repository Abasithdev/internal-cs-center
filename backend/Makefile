.PHONY: run test build docs

run:
	go run ./cmd/server

test:
	go test ./... -v

build:
	go build -o bin/server ./cmd/server

docs:
	swag init -g ./cmd/server/main.go -o ./cmd/server/docs