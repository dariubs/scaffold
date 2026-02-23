.PHONY: help build run clean deps test migrate setup apply-module

# Default target
help:
	@echo "Available commands:"
	@echo "  setup         - Interactive .env setup (step-by-step)"
	@echo "  apply-module  - Apply GO_MODULE_PATH from .env to go.mod and all Go files"
	@echo "  deps      - Install dependencies"
	@echo "  build     - Build application and migration tool"
	@echo "  run       - Run the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  migrate   - Run database migrations"

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build application and migration tool
build:
	@echo "Building application..."
	go build -o bin/index app/main/index/index.go
	@echo "Building migration tool..."
	go build -o bin/migrate app/main/migrate/migrate.go
	@echo "Build complete! Binaries are in the bin/ directory"

# Run the application
run:
	@echo "Starting application on port 3782..."
	go run app/main/index/index.go

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

# Interactive .env setup (step-by-step, orange/terminal styled)
setup:
	@bash scripts/setup-env.sh

# Apply Go module path from .env to go.mod and all .go files
apply-module:
	@bash scripts/apply-module-path.sh

# Create directories only (used when .env already exists)
setup-dirs:
	@mkdir -p bin log
	@echo "Created bin/ and log/ directories."

# Development mode - run single server
dev: setup-dirs
	@echo "Starting development server..."
	@echo "App: http://localhost:3782"
	@echo "Admin: http://localhost:3782/admin (or set ADMIN_BASE_PATH in .env)"
	@echo "Press Ctrl+C to stop"
	go run app/main/index/index.go
