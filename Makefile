.PHONY: all test build clean docs build-ui build-ui

all: build

build: build-ui
	@go build -o bin/server ./cmd/server

build-ui:
	@echo "Building UI... (Skipping: Vanilla JS)" -v

testbuild: build-ui
	@echo "Building Gantral Core..."
	@go build -o bin/server ./cmd/server

test:
	@echo "Running Tests..."
	@go test ./... -v

test-e2e:
	@echo "Running E2E Verification..."
	@# Ensure DB and Temporal are up (assumes docker-compose is running or started)
	@# We can try to start them if not, but being non-safe, we just assume or check.
	@# For CI/CD, this would block on health checks.
	@go test ./tests/e2e_workflow_test.go -v
docs:
	@echo "Starting Docusaurus..."
	cd docs-site && npm start

dev:
	@echo "Starting Dev Environment..."
	docker-compose up -d postgres
	@echo "Waiting for Postgres..."
	@sleep 3
	@sleep 3
	go run cmd/server/main.go

dev-down:
	docker-compose down

dev-reset:
	docker-compose down -v

up:
	docker-compose up -d

down:
	docker-compose down

test-integration:
	@echo "Running Integration Tests..."
	@go test ./tests -v

clean:
	@echo "Cleaning up..."
	@rm -rf bin
	@rm -f coverage.out coverage.html

help:
	@echo "Available commands:"
	@echo "  make build           - Build the server binary"
	@echo "  make run             - Run the server (requires DB)"
	@echo "  make dev             - Start Postgres + Run Server"
	@echo "  make test            - Run unit tests"
	@echo "  make test-integration - Run integration tests (requires DB)"
	@echo "  make clean           - Remove artifacts"
	@echo "  make docs            - Start documentation site"

