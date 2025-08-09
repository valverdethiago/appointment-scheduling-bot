.PHONY: run test migrate cli build clean

# Default target
all: run

# Run the API server
run:
	go run ./cmd/api

# Run tests
test:
	go test ./...

# Build the application
build:
	go build -o bin/api ./cmd/api
	go build -o bin/cli ./cmd/cli

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run with environment file
run-env:
	@if [ -f .env ]; then \
		export $$(cat .env | xargs) && go run ./cmd/api; \
	else \
		echo "No .env file found. Running with default environment..."; \
		go run ./cmd/api; \
	fi 