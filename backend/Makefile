.PHONY: dev test test-cover build lint migrate-up migrate-down

dev:
	go run ./cmd/api/...

test:
	go test ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

build:
	go build -o bin/server ./cmd/api/...

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path migrations/ -database "$$DATABASE_URL" up

migrate-down:
	migrate -path migrations/ -database "$$DATABASE_URL" down 1
