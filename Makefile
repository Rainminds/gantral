.PHONY: all test build clean docs docs-install docs-build build-ui lint

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
	@go test -v -count=1 ./tests/unit/... ./tests/statemachine/... ./tests/artifact/... ./tests/replay/... ./internal/... ./core/... ./pkg/... ./adapters/... ./cmd/... | grep -v "no test files" || true
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

docs-install:
	@echo "Installing specialized dependencies for docs site..."
	cd docs-site && npm install

docs-build:
	@echo "Building optimized production version of documentation..."
	cd docs-site && npm run build

LINT_BIN := $(shell command -v golangci-lint 2>/dev/null || echo "$(shell go env GOPATH)/bin/golangci-lint")

lint:
	@echo "Running linter ($(LINT_BIN))..."
	@$(LINT_BIN) run ./...

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
	@echo "  make docs            - Start documentation site (Dev Mode)"
	@echo "  make docs-install    - Install documentation dependencies"
	@echo "  make docs-build      - Build production documentation"
	@echo "  make lint            - Run golangci-lint (smart path resolution)"

