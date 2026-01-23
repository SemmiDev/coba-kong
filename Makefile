# Makefile untuk mempermudah common tasks
# Makefile membuat development workflow lebih efisien dengan command yang simple
# Usage: make <target>

.PHONY: help build run test clean docker-build docker-up docker-down docker-logs kong-config

# Default target yang akan dijalankan jika kita run `make` tanpa argument
help:
	@echo "Available commands:"
	@echo "  make build         - Build the Go application"
	@echo "  make run           - Run the application locally"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start all services with Docker Compose"
	@echo "  make docker-down   - Stop all services"
	@echo "  make docker-logs   - View logs from all services"
	@echo "  make kong-config   - Apply Kong configuration"
	@echo "  make kong-health   - Check Kong health"

# Build aplikasi Go untuk local development
build:
	@echo "Building application..."
	go build -o bin/main cmd/api/main.go

# Run aplikasi secara lokal (tanpa Docker)
# Pastikan dependencies seperti Postgres sudah running
run:
	@echo "Running application..."
	go run cmd/api/main.go

# Run tests dengan coverage report
test:
	@echo "Running tests..."
	go test -v -cover ./...

# Run tests dengan coverage report detail
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts dan temporary files
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Build Docker image untuk aplikasi kita
docker-build:
	@echo "Building Docker image..."
	docker build -f docker/Dockerfile -t user-service:latest .

# Start semua services dengan Docker Compose
# -d flag untuk run di background (detached mode)
docker-up:
	@echo "Starting all services..."
	docker-compose up -d
	@echo "Services started. Kong Admin UI: http://localhost:8001"
	@echo "Konga Dashboard: http://localhost:1337"
	@echo "User Service (direct): http://localhost:8080"
	@echo "Kong Proxy: http://localhost:8000"

# Stop dan remove semua containers
docker-down:
	@echo "Stopping all services..."
	docker-compose down

# Stop, remove containers, dan hapus volumes
# HATI-HATI: ini akan menghapus semua data di database
docker-down-volumes:
	@echo "Stopping all services and removing volumes..."
	docker-compose down -v

# Restart semua services
docker-restart:
	@echo "Restarting all services..."
	docker-compose restart

# View logs dari semua services
# -f flag untuk follow logs (real-time)
docker-logs:
	docker-compose logs -f

# View logs dari specific service
docker-logs-api:
	docker-compose logs -f user-service

docker-logs-kong:
	docker-compose logs -f kong

# Apply Kong configuration menggunakan decK tool
# decK adalah tool untuk manage Kong configuration declaratively
# Install decK: https://docs.konghq.com/deck/latest/installation/
kong-config:
	@echo "Applying Kong configuration..."
	deck sync --kong-addr http://localhost:8001 --state kong/kong.yml

# Check Kong health
kong-health:
	@echo "Checking Kong health..."
	curl -i http://localhost:8001/status

# List all Kong services
kong-services:
	@echo "Listing Kong services..."
	curl -s http://localhost:8001/services | jq

# List all Kong routes
kong-routes:
	@echo "Listing Kong routes..."
	curl -s http://localhost:8001/routes | jq

# Test API through Kong
test-api:
	@echo "Testing API through Kong..."
	@echo "\nCreating user..."
	curl -X POST http://localhost:8000/users \
		-H "Content-Type: application/json" \
		-d '{"name":"John Doe","email":"john@example.com"}' | jq
	@echo "\nGetting all users..."
	curl -X GET http://localhost:8000/users | jq

# Format code menggunakan gofmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter untuk code quality check
# Install golangci-lint: https://golangci-lint.run/usage/install/
lint:
	@echo "Running linter..."
	golangci-lint run

# Development mode - rebuild dan restart on code changes
# Requires air: go install github.com/cosmtrek/air@latest
dev:
	@echo "Starting development mode..."
	air
