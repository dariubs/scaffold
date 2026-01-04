
# Scaffold

A production-ready Go web application template with user authentication, Google OAuth, Cloudflare R2 file storage, and admin panel.

## Features

- User registration and login
- Google OAuth integration with CSRF protection
- Session-based authentication
- Admin panel with database-backed admin authentication
- Profile management with image uploads
- Cloudflare R2 file storage
- PostgreSQL database with GORM
- Structured logging (logrus)
- Health and readiness check endpoints
- Graceful shutdown
- Connection pooling
- CORS support
- Rate limiting
- Tailwind CSS + Alpine.js frontend

## Quick Start

### 1. Prerequisites
- Go 1.21+ (tested with Go 1.24.2)
- PostgreSQL
- Cloudflare R2 account (optional, for file uploads)
- Google OAuth credentials (optional, for OAuth login)

### 2. Setup
```bash
# Clone and navigate
cd scaffold

# Install dependencies
go mod tidy

# Copy environment file
cp env.example .env

# Edit .env with your database and API keys
# Required: DB_DSN and SESSION_SECRET
```

### 3. Environment Variables

See `env.example` for all available configuration options.

**Required:**
- `DB_DSN` - PostgreSQL connection string
- `SESSION_SECRET` - Secret key for session encryption (use a strong random string in production)

**Optional:**
- `GOOGLE_CLIENT_ID` - Google OAuth client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth client secret
- `GOOGLE_REDIRECT_URL` - Google OAuth redirect URL
- `CLOUDFLARE_ACCOUNT_ID` - Cloudflare R2 account ID
- `CLOUDFLARE_ACCESS_KEY_ID` - Cloudflare R2 access key
- `CLOUDFLARE_SECRET_ACCESS_KEY` - Cloudflare R2 secret key
- `CLOUDFLARE_R2_BUCKET` - Cloudflare R2 bucket name
- `PORT` - Main server port (default: 3782)
- `ADMIN_PORT` - Admin server port (default: 3781)
- `LOG_LEVEL` - Log level (debug, info, warn, error) (default: info)

### 4. Database Setup
```bash
# Run migrations to create tables and seed admin user
go run app/main/migrate/migrate.go
```

This will:
- Create the database schema
- Create an admin user (email: `admin@example.com`, password: `admin123`)
- Set the admin user's `IsAdmin` flag to `true`

**Important:** Change the admin password after first login in production!

### 5. Run
```bash
# Start main app
go run app/main/index/index.go

# In another terminal, start admin panel
go run app/main/admin/admin.go
```

Or use the Makefile:
```bash
# Run both servers
make dev

# Or run individually
make run-index  # Main app
make run-admin  # Admin panel
```

### 6. Access
- Main app: http://localhost:3782
- Admin panel: http://localhost:3781/admin
- Health check: http://localhost:3782/health
- Readiness check: http://localhost:3782/readiness

**Admin Login:**
1. Log in through the main app (http://localhost:3782/login)
2. Use the default admin credentials:
   - Email: `admin@example.com`
   - Password: `admin123`
3. Navigate to the admin panel (http://localhost:3781/admin)

The admin panel uses session-based authentication and checks the `IsAdmin` flag in the database.

## Project Structure
```
app/
├── config/       # Configuration management
├── database/     # Database connection and pooling
├── handlers/     # HTTP handlers
│   ├── admin/    # Admin panel handlers
│   ├── health/   # Health check handlers
│   └── index/    # Main app handlers
├── main/         # Application entry points
│   ├── admin/    # Admin server
│   ├── index/    # Main server
│   └── migrate/  # Migration tool
├── middleware/   # HTTP middleware (auth, logging, etc.)
├── model/        # Data models
└── utils/        # Utilities (R2 service, logger, validator, errors)
views/            # HTML templates
```

## Development

### Build
```bash
# Build all binaries
make build

# Build main app
go build -o bin/index app/main/index/index.go

# Build admin panel
go build -o bin/admin app/main/admin/admin.go

# Build migration tool
go build -o bin/migrate app/main/migrate/migrate.go
```

### Run Migrations
```bash
make migrate
# or
go run app/main/migrate/migrate.go
```

### Makefile Commands
```bash
make help      # Show all available commands
make deps      # Install dependencies
make build     # Build all applications
make run-index # Run main application
make run-admin # Run admin panel
make dev       # Run both servers in development mode
make clean     # Clean build artifacts
make migrate   # Run database migrations
```

## Security Features

- Session-based authentication
- OAuth state validation (CSRF protection)
- Password hashing with bcrypt
- Database-backed admin authentication
- Rate limiting support
- CORS configuration
- Secure session management

## Production Considerations

1. **Environment Variables:** Never commit `.env` files. Use environment variables or secrets management in production.

2. **SESSION_SECRET:** Use a strong, random secret key in production (at least 32 characters).

3. **Database:** Configure proper connection pooling settings based on your load.

4. **Logging:** Set `LOG_LEVEL` appropriately (use `info` or `warn` in production).

5. **Admin User:** Change the default admin password immediately.

6. **HTTPS:** Always use HTTPS in production. Configure reverse proxy (nginx, Caddy, etc.).

7. **CORS:** Configure CORS origins in `app/middleware/cors.go` to restrict access.

8. **Rate Limiting:** Configure rate limits in `app/middleware/ratelimit.go` based on your needs.

## License

See LICENSE file for details.
