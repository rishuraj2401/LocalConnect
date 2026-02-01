# LocalConnect Backend Makefile

.PHONY: help run build test test-setup test-clean clean migrate-up migrate-down dev

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	@echo "Starting LocalConnect API..."
	go run ./cmd/api

build: ## Build the application
	@echo "Building LocalConnect API..."
	go build -o bin/api ./cmd/api
	@echo "Binary created at bin/api"

dev: ## Run in development mode with auto-reload (requires air)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Running without hot reload..."; \
		go run ./cmd/api; \
	fi

test-setup: ## Setup test database
	@echo "Setting up test database..."
	@docker exec -i localconnect-postgres psql -U postgres < test_setup.sql 2>/dev/null || true
	@echo "Test database ready"

test: ## Run all tests
	@chmod +x run_tests.sh
	@./run_tests.sh

test-quick: ## Run tests without setup
	@echo "Running tests..."
	@go test -v ./internal/handlers/...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./internal/handlers/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"
	@go tool cover -func=coverage.out | tail -n 1

test-clean: ## Clean test database
	@echo "Cleaning test database..."
	@docker exec -i localconnect-postgres psql -U postgres -c "DROP DATABASE IF EXISTS localconnect_test;" 2>/dev/null || true
	@echo "Test database cleaned"

lint: ## Run linter
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install from https://golangci-lint.run/usage/install/"; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "No issues found"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@docker exec -i localconnect-postgres psql -U postgres -d localconnect < migrations/001_init.sql
	@docker exec -i localconnect-postgres psql -U postgres -d localconnect < migrations/002_messaging.sql
	@echo "Migrations complete"

db-shell: ## Open PostgreSQL shell
	@docker exec -it localconnect-postgres psql -U postgres -d localconnect

db-shell-test: ## Open test PostgreSQL shell
	@docker exec -it localconnect-postgres psql -U postgres -d localconnect_test

docker-start: ## Start PostgreSQL container
	@cd .. && docker compose up -d

docker-stop: ## Stop PostgreSQL container
	@cd .. && docker compose down

docker-logs: ## View PostgreSQL logs
	@docker logs -f localconnect-postgres

docker-reset: ## Reset PostgreSQL (WARNING: deletes all data)
	@echo "⚠️  This will delete all data. Press Ctrl+C to cancel, or Enter to continue..."
	@read confirm
	@cd .. && docker compose down -v
	@cd .. && docker compose up -d
	@sleep 3
	@$(MAKE) migrate-up
	@echo "Database reset complete"

all: fmt vet lint test build ## Run all checks and build

.DEFAULT_GOAL := help
