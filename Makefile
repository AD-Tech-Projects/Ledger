.PHONY: run tidy clean test test-coverage build migrate-create migrate-apply migrate-rollback

MAIN_PATH=./main.go
BINARY_NAME=out/app

-include .env

DB_URL=postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_DIR=file://migrations

run:
	@echo "Running the application..."
	go run $(MAIN_PATH)

build:
	@echo "Building the application..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)

tidy:
	@echo "Updating and cleaning Go module dependencies..."
	go mod tidy

clean:
	@echo "Cleaning build artifacts and module cache..."
	go clean -modcache
	rm -rf out/

test:
	@echo "Running all project tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

migrate-create:
	@echo "Creating a new migration: $(name)"
	@atlas migrate diff $(name) --env gorm
	
migrate-hash:
	@echo "Generating migration hash..."
	@atlas migrate hash --env gorm

migrate-apply:
	@echo "Applying migrations to the database..."
	@atlas migrate apply --url "$(DB_URL)" --dir "$(MIGRATIONS_DIR)"

migrate-rollback:
	@echo "Rolling back the last migration..."
	@atlas migrate rollback --url "$(DB_URL)" --dir "$(MIGRATIONS_DIR)"
