.PHONY: all test build clean

all: build

build:
	@echo "Building Gantral Core..."
	# go build ./core/...

test:
	@echo "Running Tests..."
	# go test ./...
