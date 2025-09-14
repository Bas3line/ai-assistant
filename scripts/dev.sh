#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

check_prerequisites() {
    print_status "Checking prerequisites..."
    
    if ! command_exists docker; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    if ! command_exists docker-compose; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go 1.24+ first."
        exit 1
    fi
    
    print_success "All prerequisites are installed"
}

setup_environment() {
    print_status "Setting up environment..."
    
    if [ ! -f .env ]; then
        if [ -f .env.dev ]; then
            cp .env.dev .env
            print_success "Copied .env.dev to .env"
        else
            cp .env.example .env
            print_warning "Copied .env.example to .env. Please update with your API keys."
        fi
    else
        print_status ".env file already exists"
    fi
}

start_services() {
    print_status "Starting Docker services (PostgreSQL and Redis)..."
    
    docker-compose -f docker/docker-compose.dev.yml down 2>/dev/null || true
    
    docker-compose -f docker/docker-compose.dev.yml up -d
    
    print_status "Waiting for services to be ready..."
    
    print_status "Waiting for PostgreSQL..."
    timeout=30
    while [ $timeout -gt 0 ]; do
        if docker-compose -f docker/docker-compose.dev.yml exec postgres pg_isready -U postgres -d ai_assistant >/dev/null 2>&1; then
            print_success "PostgreSQL is ready"
            break
        fi
        sleep 1
        timeout=$((timeout - 1))
    done
    
    if [ $timeout -eq 0 ]; then
        print_error "PostgreSQL failed to start within 30 seconds"
        exit 1
    fi
    
    print_status "Waiting for Redis..."
    timeout=30
    while [ $timeout -gt 0 ]; do
        if docker-compose -f docker/docker-compose.dev.yml exec redis redis-cli ping >/dev/null 2>&1; then
            print_success "Redis is ready"
            break
        fi
        sleep 1
        timeout=$((timeout - 1))
    done
    
    if [ $timeout -eq 0 ]; then
        print_error "Redis failed to start within 30 seconds"
        exit 1
    fi
    
    print_success "All services are running"
}

build_app() {
    print_status "Building the application..."
    
    go mod download
    go mod tidy
    
    go build -o bin/ai-assistant ./cmd/api
    
    print_success "Application built successfully"
}

run_migrations() {
    print_status "Checking for database migrations..."
    
    if [ -d "migrations" ] && [ "$(ls -A migrations)" ]; then
        print_status "Running database migrations..."
        print_warning "Migration logic not implemented yet. Please run migrations manually if needed."
    else
        print_status "No migrations found"
    fi
}

display_info() {
    print_success "Development environment is ready!"
    echo
    echo "Services running:"
    echo "   - PostgreSQL: localhost:5432"
    echo "     Database: ai_assistant"
    echo "     User: postgres"
    echo "     Password: postgres"
    echo
    echo "   - Redis: localhost:6379"
    echo
    echo "To start the application:"
    echo "   ./bin/ai-assistant"
    echo
    echo "   Or for development with auto-reload:"
    echo "   go run ./cmd/api"
    echo
    echo "Useful commands:"
    echo "   - Stop services: docker-compose -f docker/docker-compose.dev.yml down"
    echo "   - View logs: docker-compose -f docker/docker-compose.dev.yml logs -f"
    echo "   - Connect to PostgreSQL: docker-compose -f docker/docker-compose.dev.yml exec postgres psql -U postgres -d ai_assistant"
    echo "   - Connect to Redis: docker-compose -f docker/docker-compose.dev.yml exec redis redis-cli"
    echo
    echo "API will be available at: http://localhost:8080"
    echo "Health check: http://localhost:8080/health"
    echo
}

cleanup() {
    if [ $? -ne 0 ]; then
        print_error "Setup failed. Cleaning up..."
        docker-compose -f docker/docker-compose.dev.yml down 2>/dev/null || true
    fi
}

main() {
    echo "AI Assistant Development Setup"
    echo "================================="
    echo
    
    trap cleanup EXIT
    
    check_prerequisites
    setup_environment
    start_services
    build_app
    run_migrations
    display_info
    
    trap - EXIT
}

case "${1:-}" in
    "stop")
        print_status "Stopping development services..."
        docker-compose -f docker/docker-compose.dev.yml down
        print_success "Services stopped"
        ;;
    "restart")
        print_status "Restarting development services..."
        docker-compose -f docker/docker-compose.dev.yml down
        main
        ;;
    "logs")
        docker-compose -f docker/docker-compose.dev.yml logs -f
        ;;
    "status")
        docker-compose -f docker/docker-compose.dev.yml ps
        ;;
    "clean")
        print_warning "This will remove all Docker containers and volumes. Are you sure? (y/N)"
        read -r response
        if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
            docker-compose -f docker/docker-compose.dev.yml down -v
            docker system prune -f
            print_success "Cleanup completed"
        else
            print_status "Cleanup cancelled"
        fi
        ;;
    "help"|"-h"|"--help")
        echo "AI Assistant Development Script"
        echo
        echo "Usage: $0 [command]"
        echo
        echo "Commands:"
        echo "  (no args)  Start development environment"
        echo "  stop       Stop all services"
        echo "  restart    Restart all services"
        echo "  logs       Show service logs"
        echo "  status     Show service status"
        echo "  clean      Remove all containers and volumes"
        echo "  help       Show this help"
        ;;
    "")
        main
        ;;
    *)
        print_error "Unknown command: $1"
        echo "Use '$0 help' for usage information"
        exit 1
        ;;
esac