.PHONY: all test build clean docs

all: build

build:
	@echo "Building Gantral Core..."
	# go build ./core/...

test:
	@echo "Running Tests..."
	# go test ./...
docs:
	@echo "Starting Docusaurus..."
	cd docs-site && npm start
