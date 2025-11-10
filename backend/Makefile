.PHONY: run test build docs coverage

run:
	go run -tags=skip_coverage ./cmd/server

test:
	@go test -v -tags=skip_coverage ./... -coverprofile=coverage.out -coverpkg=./...

coverage:
	@go tool cover -html=coverage.out -o coverage.html
	@start coverage.html

build:
	go build -tags=skip_coverage -o bin/server ./cmd/server

docs:
	swag init -g ./cmd/server/main.go -o ./docs