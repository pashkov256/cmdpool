.PHONY: build run test clean install deps

# Binary name
BINARY_NAME=cmdpool

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty)"

# Default target
all: build

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build the application
build: deps
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/cmdpool

# Build for different platforms
build-linux: deps
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 ./cmd/cmdpool

build-windows: deps
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe ./cmd/cmdpool

build-darwin: deps
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 ./cmd/cmdpool

# Run the application
run: build
	./$(BINARY_NAME)

# Run in TUI mode
run-tui: build
	./$(BINARY_NAME)

# Run in CLI mode with example commands
run-cli: build
	./$(BINARY_NAME) "echo 'Hello World'" "sleep 5" "ping -c 3 google.com"

# Run tests
test: deps
	go test -v ./...

# Run tests with coverage
test-coverage: deps
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out

# Install to system
install: build
	go install ./cmd/cmdpool

# Development mode
dev: deps
	go run ./cmd/cmdpool

# Development mode with CLI
dev-cli: deps
	go run ./cmd/cmdpool "echo 'Dev mode'" "sleep 2"

# Format code
fmt:
	go fmt ./...

# Lint code
lint: deps
	golangci-lint run

# Generate documentation
docs:
	godoc -http=:6060

# Create release builds
release: clean build-linux build-windows build-darwin
	tar -czf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	tar -czf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	zip $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Build and run the application"
	@echo "  run-tui      - Run in TUI mode"
	@echo "  run-cli      - Run in CLI mode with example commands"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  install      - Install to system"
	@echo "  dev          - Run in development mode"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  release      - Create release builds for all platforms" 