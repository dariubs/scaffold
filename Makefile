.PHONY: help build run-admin run-index clean deps test migrate

# Default target
help:
	@echo "Available commands:"
	@echo "  deps      - Install dependencies"
	@echo "  build     - Build both applications"
	@echo "  run-index - Run the main application"
	@echo "  run-admin - Run the admin panel"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  migrate   - Run database migrations"

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build both applications
build:
	@echo "Building main application..."
	go build -o bin/index app/main/index/index.go
	@echo "Building admin panel..."
	go build -o bin/admin app/main/admin/admin.go
	@echo "Building migration tool..."
	go build -o bin/migrate app/main/migrate/migrate.go
	@echo "Build complete! Binaries are in the bin/ directory"

# Run the main application
run-index:
	@echo "Starting main application on port 3782..."
	go run app/main/index/index.go

# Run the admin panel
run-admin:
	@echo "Starting admin panel on port 3781..."
	go run app/main/admin/admin.go

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run tests
test:
	go test ./...

# Run database migrations
migrate:
	@echo "Running database migrations..."
	go run app/main/migrate/migrate.go

# Create necessary directories
setup:
	mkdir -p bin
	mkdir -p log
	@echo "Setup complete! Created bin/ and log/ directories"

# Development mode - run both servers
dev: setup
	@echo "Starting development environment..."
	@echo "Main app: http://localhost:3782"
	@echo "Admin panel: http://localhost:3781/admin"
	@echo "Press Ctrl+C to stop both servers"
	@trap 'kill %1 %2' SIGINT; \
	go run app/main/index/index.go & \
	go run app/main/admin/admin.go & \
	wait
