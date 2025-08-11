
# Scaffold

A simple Go web application template with user authentication, Google OAuth, and Cloudflare R2 file storage.

## Features

- User registration and login
- Google OAuth integration
- Profile management with image uploads
- Cloudflare R2 file storage
- PostgreSQL database with GORM
- Tailwind CSS + Alpine.js frontend

## Quick Start

### 1. Prerequisites
- Go 1.21+
- PostgreSQL
- Cloudflare R2 account (optional, for file uploads)

### 2. Setup
```bash
# Clone and navigate
cd scaffold

# Install dependencies
go mod tidy

# Copy environment file
cp env.example .env

# Edit .env with your database and API keys
```

### 3. Environment Variables
```bash
# Required
DB_DSN=postgres://user:pass@localhost:5432/scaffold_db?sslmode=disable
SESSION_SECRET=your-secret-key

# Optional (for Google OAuth)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

# Optional (for R2 file storage)
CLOUDFLARE_ACCOUNT_ID=your-account-id
CLOUDFLARE_ACCESS_KEY_ID=your-access-key
CLOUDFLARE_SECRET_ACCESS_KEY=your-secret-key
CLOUDFLARE_R2_BUCKET=your-bucket-name
```

### 4. Run
```bash
# Start main app
go run app/main/index/index.go

# Start admin panel (optional)
go run app/main/admin/admin.go
```

### 5. Access
- Main app: http://localhost:3782
- Admin: http://localhost:3781/admin (admin/admin123)

## Project Structure
```
app/
├── database/     # Database connection
├── handlers/     # HTTP handlers
├── main/         # Application entry points
├── model/        # Data models
└── utils/        # Utilities (R2 service)
views/            # HTML templates
```

## Development
```bash
# Build
go build -o scaffold app/main/index/index.go

# Test
go test ./...

# Run with Makefile
make dev
```
