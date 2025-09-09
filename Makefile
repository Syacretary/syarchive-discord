# Makefile for Discord Bot

# Variables
BINARY_NAME=discord-bot
MAIN_FILE=cmd/bot/main.go

# Default target
all: build

# Build the application
build:
	go build -o ${BINARY_NAME} ${MAIN_FILE}

# Run the application
run:
	go run ${MAIN_FILE}

# Install dependencies
deps:
	go mod tidy

# Clean build files
clean:
	rm -f ${BINARY_NAME}

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run tests with race detector
test-race:
	go test -race ./...

# Run benchmarks
test-bench:
	go test -bench=. ./...

# Run tests with verbose output and coverage
test-verbose-cover:
	go test -v -cover ./...

# Build Docker image
docker-build:
	docker build -t ${BINARY_NAME} .

# Run Docker container
docker-run:
	docker run --env-file .env ${BINARY_NAME}

# Run with docker-compose
docker-compose-up:
	docker-compose up -d

# Stop docker-compose
docker-compose-down:
	docker-compose down

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Run all quality checks
check: fmt vet test

.PHONY: all build run deps clean test test-cover test-race test-bench test-verbose-cover docker-build docker-run docker-compose-up docker-compose-down fmt vet check