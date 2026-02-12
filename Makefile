.PHONY: all test build clean docs build-ui build-ui

all: build

build: build-ui
	@go build -o bin/server ./cmd/server

build-ui:
	@echo "Building UI... (Skipping: Vanilla JS)" -v

testbuild: build-ui
	@echo "Building Gantral Core..."
	@go build -o bin/server ./cmd/server

test-tier1:
	@echo "Running Tier 1 Tests (Scope: Unit, StateMachine, Artifact, Replay, Core, Pkg)..."
	@go test -v -count=1 ./tests/unit/... ./tests/statemachine/... ./tests/artifact/... ./tests/replay/... ./core/... ./pkg/... ./adapters/... ./cmd/... | grep -v "no test files" || true
	@echo "Tier 1 Tests Completed."

test-tier2:
	@echo "Running Tier 2 Integration Tests..."
	@go test -tags=integration -v ./tests/integration/...

test-integration: test-tier2

test: test-tier1 test-tier2
	@echo "All Tests Completed."

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

