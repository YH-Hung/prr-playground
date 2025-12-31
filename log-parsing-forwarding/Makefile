.PHONY: all build test lint coverage docker-build docker-up docker-down clean help

# Default target
all: lint test build

help:
	@echo "Available targets:"
	@echo "  make build          - Build all service binaries"
	@echo "  make build-server   - Build server binary"
	@echo "  make build-client   - Build client binary"
	@echo "  make test           - Run all tests"
	@echo "  make test-server    - Run server tests only"
	@echo "  make test-client    - Run client tests only"
	@echo "  make test-integration - Run integration tests"
	@echo "  make lint           - Run linter (requires golangci-lint)"
	@echo "  make coverage       - Generate test coverage report"
	@echo "  make docker-build   - Build Docker images"
	@echo "  make docker-up      - Start services with Docker Compose"
	@echo "  make docker-down    - Stop services with Docker Compose"
	@echo "  make clean          - Clean build artifacts and caches"

# Build targets
build: build-server build-client

build-server:
	@echo "Building server..."
	@cd services/server && go build -o ../../bin/server
	@echo "Server binary created at bin/server"

build-client:
	@echo "Building client..."
	@cd services/client && go build -o ../../bin/client
	@echo "Client binary created at bin/client"

# Test targets
test: test-internal test-integration
	@echo "All tests completed successfully!"

test-all: test-internal test-server test-client test-integration
	@echo "All tests completed successfully!"

test-server:
	@echo "Running server tests..."
	@cd services/server && go test ./...

test-client:
	@echo "Running client tests..."
	@cd services/client && go test ./...

test-integration:
	@echo "Running integration tests..."
	@cd test/integration && go test -v ./...

test-internal:
	@echo "Running internal package tests..."
	@go test github.com/yinghanhung/prr-playground/internal/config \
		github.com/yinghanhung/prr-playground/internal/trace \
		github.com/yinghanhung/prr-playground/internal/logger \
		github.com/yinghanhung/prr-playground/internal/retry

# Coverage target
coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./internal/... ./services/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total:

# Linting target
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Install it from https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

# Docker targets
docker-build:
	@echo "Building Docker images..."
	@docker compose build

docker-up:
	@echo "Starting services..."
	@docker compose up -d

docker-down:
	@echo "Stopping services..."
	@docker compose down

docker-logs:
	@docker compose logs -f

docker-restart: docker-down docker-build docker-up

# Cleanup target
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean -testcache
	@echo "Cleanup complete"

# Development helpers
run-server:
	@echo "Running server locally..."
	@cd services/server && go run main.go

run-client:
	@echo "Running client locally..."
	@cd services/client && go run main.go

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

tidy:
	@echo "Tidying modules..."
	@cd services/server && go mod tidy
	@cd services/client && go mod tidy
	@cd test/integration && go mod tidy
	@go work sync

# Check if dependencies are up to date
deps-check:
	@echo "Checking for outdated dependencies..."
	@cd services/server && go list -u -m all
	@cd services/client && go list -u -m all
