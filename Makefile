# Binary name
BINARY_NAME=learn.exe
BINARY_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

# Linker flags to strip debug info and reduce binary size
LDFLAGS=-ldflags "-w -s"

.PHONY: all build clean test coverage run migrate lint help

all: clean build test

build:
	@echo "Building binary..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) main.go
	@echo "Build complete: $(BINARY_DIR)/$(BINARY_NAME)"

run:
	@echo "Starting server..."
	$(GORUN) main.go serve

migrate:
	@echo "Running database migrations..."
	$(GORUN) main.go migrate

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BINARY_DIR) coverage.out coverage.html
	@echo "Clean complete"

lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

help:
	@echo "Available commands:"
	@echo "  make build     - Build the application binary"
	@echo "  make run       - Run the API server"
	@echo "  make migrate   - Run database migrations"
	@echo "  make test      - Run unit tests"
	@echo "  make coverage  - Run tests and generate coverage report"
	@echo "  make clean     - Remove binary and build artifacts"
	@echo "  make lint      - Check code style and errors"
