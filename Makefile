.PHONY: help proto sqlc migrate-up migrate-down test test-integration build run clean migrate-status migrate-new lint

DB_PATH=db/expenses.db
MIGRATIONS_DIR=db/migrations

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

proto: ## Generate connect-rpc Go code and OpenAPI spec from proto files
	@echo "Generating Protobuf/Connect code..."
	buf generate proto

sqlc: ## Generate Go code from SQL queries using sqlc
	@echo "Generating sqlc code..."
	sqlc generate

generate-all: proto sqlc ## Generate all code (protobuf, connect, sqlc)

migrate: ## Apply all database migrations
	@echo "Applying migrations..."
	atlas migrate apply --url "sqlite://$(DB_PATH)" --dir "file://$(MIGRATIONS_DIR)"

migrate-status: ## Check migration status
	@echo "Checking migration status..."
	atlas migrate status --url "sqlite://$(DB_PATH)" --dir "file://$(MIGRATIONS_DIR)"

migrate-new: ## Create a new migration file. Usage: make migrate-new name=migration_name
	@echo "Creating new migration..."
	atlas migrate diff $(name) --dir "file://$(MIGRATIONS_DIR)" --to "file://db/schema.sql" --dev-url "sqlite://file?mode=memory"

test: ## Run tests with real database
	@echo "Running tests..."
	richgo test -v -race ./...

build: ## Build the server binary
	@echo "Building server..."
	mkdir -p bin
	go build -o ./bin/expense-manager

run: build ## Build and run the server
	@echo "Running server..."
	./bin/expense-manager serve

dev: ## Build and run the server in development mode
	@echo "Running server in development mode..."
	air

seed: build ## Seed the database with master data
	@echo "Seeding database with master data..."
	./bin/expense-manager seed

clean: ## Clean generated files and build artifacts
	@echo "Cleaning generated files and build artifacts..."
	rm -rf bin/*
	rm -rf internal/repo/gen/*
	rm -rf internal/rpc/gen/*
	# Keep .gitkeep files if any
	find internal/repo/gen/ -type f ! -name '.gitkeep' -delete
	find internal/rpc/gen/ -type f ! -name '.gitkeep' -delete
	find bin/ -type f ! -name '.gitkeep' -delete

# Combine setup steps
setup-tools:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install ariga.io/atlas/cmd/atlas@latest
	go install github.com/kyoh86/richgo@latest
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

lint: ## Run golangci-lint on the project
	@echo "Running linter..."
	golangci-lint run
