.PHONY: all test build clean docs

all: build

build:
	@echo "Building Gantral Core..."
	# go build ./core/...

test:
	@echo "Running Tests..."
	go test ./... -v
docs:
	@echo "Starting Docusaurus..."
	cd docs-site && npm start

dev:
	@echo "Starting Dev Environment..."
	docker-compose up -d postgres
	@echo "Waiting for Postgres..."
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
