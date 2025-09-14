#!/bin/bash

# Development setup script
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 AI Assistant Development Setup${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.24+ first.${NC}"
    exit 1
fi

GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | cut -c3-)
echo -e "${GREEN}✓ Go version: $GO_VERSION${NC}"

# Check if required tools are installed
echo -e "${YELLOW}📦 Installing development tools...${NC}"

# Install Air for hot reloading
if ! command -v air &> /dev/null; then
    echo "Installing Air (hot reloading)..."
    go install github.com/cosmtrek/air@latest
fi

# Install golangci-lint
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

# Install gosec
if ! command -v gosec &> /dev/null; then
    echo "Installing gosec..."
    go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
fi

# Copy environment file
if [ ! -f .env ]; then
    echo -e "${YELLOW}📝 Creating .env file from template...${NC}"
    cp .env.example .env
    echo -e "${YELLOW}⚠️  Please update .env with your actual values${NC}"
fi

# Install Go dependencies
echo -e "${YELLOW}📦 Installing Go dependencies...${NC}"
go mod tidy
go mod download

# Create necessary directories
echo -e "${YELLOW}📁 Creating directories...${NC}"
mkdir -p logs tmp

# Check if Docker is available
if command -v docker &> /dev/null; then
    echo -e "${GREEN}✓ Docker is available${NC}"
    if command -v docker-compose &> /dev/null; then
        echo -e "${GREEN}✓ Docker Compose is available${NC}"
        echo -e "${BLUE}💡 You can use 'docker-compose up -d' to start PostgreSQL and Redis${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Docker not found. You'll need to set up PostgreSQL and Redis manually${NC}"
fi

# Build the application
echo -e "${YELLOW}🔨 Building application...${NC}"
make build

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Setup completed successfully!${NC}"
    echo ""
    echo -e "${BLUE}🎯 Next steps:${NC}"
    echo "1. Update your .env file with actual API keys and database URLs"
    echo "2. Start PostgreSQL and Redis (docker-compose up -d)"
    echo "3. Run database migrations (make db-push)"
    echo "4. Start the development server (make dev)"
    echo ""
    echo -e "${BLUE}📚 Available commands:${NC}"
    echo "  make dev       - Start development server with hot reload"
    echo "  make build     - Build the application"
    echo "  make test      - Run tests"
    echo "  make lint      - Run linter"
    echo "  make db-push   - Push database schema"
else
    echo -e "${RED}❌ Setup failed during build${NC}"
    exit 1
fi