BINARY_NAME=ai-assistant
BUILD_DIR=./bin

.PHONY: build run clean test deps dev-setup dev-start dev-stop dev-logs

# Build the application
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

# Run the application
run: build
	$(BUILD_DIR)/$(BINARY_NAME)

# Development mode with hot reload (requires air)
dev:
	air -c scripts/.air.toml

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Run tests
test:
	go test -v ./...

# Test with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	go mod tidy
	go mod download

# Development environment setup
dev-setup:
	@echo "Setting up development environment..."
	./scripts/dev.sh

# Start development services
dev-start:
	@echo "Starting development services..."
	docker-compose -f docker/docker-compose.dev.yml up -d

# Stop development services
dev-stop:
	@echo "Stopping development services..."
	docker-compose -f docker/docker-compose.dev.yml down

# Show development service logs
dev-logs:
	docker-compose -f docker/docker-compose.dev.yml logs -f

# Full development restart
dev-restart:
	@echo "Restarting development environment..."
	./scripts/dev.sh restart

# Development status
dev-status:
	docker-compose -f docker/docker-compose.dev.yml ps

# Clean development environment
dev-clean:
	@echo "Cleaning development environment..."
	./scripts/dev.sh clean

# Database operations (development)
db-connect:
	docker-compose -f docker/docker-compose.dev.yml exec postgres psql -U postgres -d ai_assistant

# Redis operations (development)
redis-connect:
	docker-compose -f docker/docker-compose.dev.yml exec redis redis-cli

# Generate Prisma client
prisma-generate:
	go run github.com/steebchen/prisma-client-go generate

# Database operations
db-push:
	prisma db push

db-migrate:
	prisma migrate dev

db-studio:
	prisma studio

# Docker operations (production)
docker-build:
	docker build -f docker/Dockerfile -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME)

docker-up:
	docker-compose -f docker/docker-compose.yml up -d

docker-down:
	docker-compose -f docker/docker-compose.yml down

# Linting and formatting
lint:
	golangci-lint run -c scripts/.golangci.yml

fmt:
	go fmt ./...

vet:
	go vet ./...

# Security check
sec:
	gosec ./...

# Install development tools
install-tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Help
help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "  run             - Build and run the application"
	@echo "  dev             - Run in development mode with hot reload"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  deps            - Download dependencies"
	@echo ""
	@echo "Development environment:"
	@echo "  dev-setup       - Set up development environment (PostgreSQL + Redis)"
	@echo "  dev-start       - Start development services"
	@echo "  dev-stop        - Stop development services"
	@echo "  dev-restart     - Restart development environment"
	@echo "  dev-logs        - Show development service logs"
	@echo "  dev-status      - Show development service status"
	@echo "  dev-clean       - Clean development environment"
	@echo ""
	@echo "Database:"
	@echo "  db-connect      - Connect to development PostgreSQL"
	@echo "  redis-connect   - Connect to development Redis"
	@echo ""
	@echo "Code quality:"
	@echo "  lint            - Run linter"
	@echo "  fmt             - Format code"
	@echo "  vet             - Run go vet"
	@echo "  sec             - Run security check"