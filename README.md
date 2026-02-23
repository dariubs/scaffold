
# Scaffold

A production-ready Go web application template with user authentication, Google OAuth, Cloudflare R2 file storage, and admin panel.

## Features

- User registration and login
- Multiple login methods: password, Google, GitHub, LinkedIn, X (Twitter) — each can be enabled/disabled via .env
- OAuth integration with CSRF protection
- Session-based authentication
- Admin panel with database-backed admin authentication
- Profile management with image uploads
- Cloudflare R2 file storage
- Email via Resend (welcome email on registration when configured)
- PostgreSQL database with GORM
- Structured logging (stdlib slog)
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
cp .env.example .env

# Edit .env with your database and API keys
# Required: DB_DSN and SESSION_SECRET
```

### 3. Environment Variables

See `.env.example` for all available configuration options.

**Required:**
- `DB_DSN` - PostgreSQL connection string
- `SESSION_SECRET` - Secret key for session encryption (use a strong random string in production)

**Optional (login methods):** Set to `true`, `1`, or `yes` to enable; omit or set to `false` to disable. Configure credentials for each provider you use.
- `LOGIN_PASSWORD_ENABLED` - Username/password login (default: true)
- `LOGIN_GOOGLE_ENABLED` - Google OAuth (default: true)
- `LOGIN_GITHUB_ENABLED` - GitHub OAuth (default: false)
- `LOGIN_LINKEDIN_ENABLED` - LinkedIn OAuth (default: false)
- `LOGIN_X_ENABLED` - X (Twitter) OAuth (default: false)

**Optional (OAuth credentials):** Create an app on each platform and set the callback URL to e.g. `http://localhost:3782/auth/github/callback` (or `/auth/linkedin/callback`, `/auth/x/callback`).
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL`
- `GITHUB_CLIENT_ID`, `GITHUB_CLIENT_SECRET`, `GITHUB_REDIRECT_URL`
- `LINKEDIN_CLIENT_ID`, `LINKEDIN_CLIENT_SECRET`, `LINKEDIN_REDIRECT_URL`
- `X_CLIENT_ID`, `X_CLIENT_SECRET`, `X_REDIRECT_URL`
- `CLOUDFLARE_ACCOUNT_ID` - Cloudflare R2 account ID
- `CLOUDFLARE_ACCESS_KEY_ID` - Cloudflare R2 access key
- `CLOUDFLARE_SECRET_ACCESS_KEY` - Cloudflare R2 secret key
- `CLOUDFLARE_R2_BUCKET` - Cloudflare R2 bucket name
- `RESEND_API_KEY` - Resend API key (optional; when set with `RESEND_FROM`, welcome emails are sent on registration)
- `RESEND_FROM` - Sender address for transactional email (e.g. `Scaffold <onboarding@resend.dev>`)
- `PORT` - Server port (default: 3782)
- `ADMIN_BASE_PATH` - Admin panel URL path (default: admin, e.g. /admin)
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
go run app/main/index/index.go
```

Or use the Makefile:
```bash
make dev
# or
make run
```

### 6. Access
- App: http://localhost:3782
- Admin panel: http://localhost:3782/admin (path configurable via `ADMIN_BASE_PATH`)
- Health check: http://localhost:3782/health
- Readiness check: http://localhost:3782/readiness

**Admin Login:**
1. Log in through the app (http://localhost:3782/login)
2. Use the default admin credentials:
   - Email: `admin@example.com`
   - Password: `admin123`
3. Navigate to the admin panel (http://localhost:3782/admin)

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
│   ├── index/    # Main server (serves app and admin)
│   └── migrate/  # Migration tool
├── middleware/   # HTTP middleware (auth, logging, etc.)
├── model/        # Data models
└── utils/        # Utilities (R2 service, logger, validator, errors)
views/            # HTML templates
```

## Development

### Build
```bash
make build

# Or manually:
go build -o bin/index app/main/index/index.go
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
make build     # Build application and migration tool
make run       # Run application
make dev       # Run in development mode
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
