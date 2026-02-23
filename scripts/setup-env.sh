#!/usr/bin/env bash
# Interactive .env setup for Scaffold - step-by-step with orange/terminal styling.

set -e

# Ensure we read from the terminal (fixes "make setup" when Make does not attach stdin)
if [ ! -t 0 ]; then
  if [ -e /dev/tty ]; then
    exec 0</dev/tty
  else
    echo "Error: This script must be run from an interactive terminal (e.g. run: make setup)" >&2
    echo "       No TTY available. Run from your terminal, not from an IDE run button." >&2
    exit 1
  fi
fi

# Colors (orange theme + terminal accents)
ORANGE='\033[38;5;208m'
BOLD='\033[1m'
DIM='\033[2m'
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RESET='\033[0m'

# Project root (parent of scripts/)
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
ENV_FILE="$ROOT/.env"

prompt() {
  local label="$2"
  local default="$3"
  local required="${4:-0}"
  local hint="$5"
  local val=""
  local prompt_text=""
  if [ -n "$hint" ]; then
    echo -e "${DIM}  $hint${RESET}" >&2
  fi
  if [ -n "$default" ]; then
    prompt_text="${ORANGE}${BOLD}  $label ${DIM}[$default]${RESET}: "
    read -r -p "$(echo -e "$prompt_text")" val
    val="${val:-$default}"
  else
    prompt_text="${ORANGE}${BOLD}  $label${RESET}: "
    read -r -p "$(echo -e "$prompt_text")" val
    while [ "$required" -eq 1 ] && [ -z "$val" ]; do
      echo -e "  ${YELLOW}This value is required.${RESET}" >&2
      read -r -p "$(echo -e "$prompt_text")" val
    done
  fi
  printf '%s' "$val"
}

section() {
  echo ""
  echo -e "${ORANGE}${BOLD}=== $1 ===${RESET}"
  echo ""
}

note() {
  echo -e "${CYAN}  $1${RESET}"
}

is_enabled() {
  case "$(echo "$1" | tr '[:upper:]' '[:lower:]')" in
    true|1|yes) return 0 ;;
    *) return 1 ;;
  esac
}

# Banner
echo ""
echo -e "${ORANGE}${BOLD}"
echo "  +-------------------------------------------+"
echo "  |     Scaffold - Environment Setup          |"
echo "  |     Create .env step by step              |"
echo "  +-------------------------------------------+"
echo -e "${RESET}"
echo -e "${DIM}  Answer the prompts below. Press Enter to use [default] when shown.${RESET}"

if [ -f "$ENV_FILE" ]; then
  echo -e "${YELLOW}  Warning: .env already exists. This will overwrite it.${RESET}"
  confirm_prompt="${BOLD}  Continue? (y/N): ${RESET}"
  read -r -p "$(echo -e "$confirm_prompt")" confirm
  case "$confirm" in
    [yY][eE][sS]|[yY]) ;;
    *) echo "  Aborted."; exit 0 ;;
  esac
fi

# Project (Go module path)
section "0. Project"
note "Go module path used in go.mod and all imports (e.g. github.com/username/myapp)."
GO_MODULE_PATH=$(prompt GO_MODULE_PATH "GO_MODULE_PATH" "github.com/dariubs/scaffold" 0 "")

# Required
section "1. Required"
note "Database connection string (PostgreSQL)."
DB_DSN=$(prompt DB_DSN "DB_DSN" "postgres://username:password@localhost:5432/scaffold_db?sslmode=disable" 0 "")
note "Secret for session encryption (use a long random string in production)."
SESSION_SECRET=$(prompt SESSION_SECRET "SESSION_SECRET" "your-super-secret-session-key-here" 0 "")

# Server
section "2. Server"
PORT=$(prompt PORT "PORT" "3782" 0 "")
ADMIN_BASE_PATH=$(prompt ADMIN_BASE_PATH "ADMIN_BASE_PATH" "admin" 0 "URL path for admin panel (e.g. /admin or /master).")

# Login toggles
section "3. Login methods (enable/disable)"
note "Set to true/1/yes to enable, false/0/no to disable."
LOGIN_PASSWORD_ENABLED=$(prompt LOGIN_PASSWORD_ENABLED "LOGIN_PASSWORD_ENABLED" "true" 0 "")
LOGIN_GOOGLE_ENABLED=$(prompt LOGIN_GOOGLE_ENABLED "LOGIN_GOOGLE_ENABLED" "true" 0 "")
LOGIN_GITHUB_ENABLED=$(prompt LOGIN_GITHUB_ENABLED "LOGIN_GITHUB_ENABLED" "false" 0 "")
LOGIN_LINKEDIN_ENABLED=$(prompt LOGIN_LINKEDIN_ENABLED "LOGIN_LINKEDIN_ENABLED" "false" 0 "")
LOGIN_X_ENABLED=$(prompt LOGIN_X_ENABLED "LOGIN_X_ENABLED" "false" 0 "")

# Google OAuth (only if enabled)
if is_enabled "$LOGIN_GOOGLE_ENABLED"; then
  section "4. Google OAuth"
  note "Set all three to enable Google login."
  GOOGLE_CLIENT_ID=$(prompt GOOGLE_CLIENT_ID "GOOGLE_CLIENT_ID" "" 0 "")
  GOOGLE_CLIENT_SECRET=$(prompt GOOGLE_CLIENT_SECRET "GOOGLE_CLIENT_SECRET" "" 0 "")
  GOOGLE_REDIRECT_URL=$(prompt GOOGLE_REDIRECT_URL "GOOGLE_REDIRECT_URL" "http://localhost:${PORT}/auth/google/callback" 0 "")
else
  GOOGLE_CLIENT_ID=""
  GOOGLE_CLIENT_SECRET=""
  GOOGLE_REDIRECT_URL=""
fi

# GitHub OAuth (only if enabled)
if is_enabled "$LOGIN_GITHUB_ENABLED"; then
  section "5. GitHub OAuth"
  GITHUB_CLIENT_ID=$(prompt GITHUB_CLIENT_ID "GITHUB_CLIENT_ID" "" 0 "")
  GITHUB_CLIENT_SECRET=$(prompt GITHUB_CLIENT_SECRET "GITHUB_CLIENT_SECRET" "" 0 "")
  GITHUB_REDIRECT_URL=$(prompt GITHUB_REDIRECT_URL "GITHUB_REDIRECT_URL" "http://localhost:${PORT}/auth/github/callback" 0 "")
else
  GITHUB_CLIENT_ID=""
  GITHUB_CLIENT_SECRET=""
  GITHUB_REDIRECT_URL=""
fi

# LinkedIn OAuth (only if enabled)
if is_enabled "$LOGIN_LINKEDIN_ENABLED"; then
  section "6. LinkedIn OAuth"
  LINKEDIN_CLIENT_ID=$(prompt LINKEDIN_CLIENT_ID "LINKEDIN_CLIENT_ID" "" 0 "")
  LINKEDIN_CLIENT_SECRET=$(prompt LINKEDIN_CLIENT_SECRET "LINKEDIN_CLIENT_SECRET" "" 0 "")
  LINKEDIN_REDIRECT_URL=$(prompt LINKEDIN_REDIRECT_URL "LINKEDIN_REDIRECT_URL" "http://localhost:${PORT}/auth/linkedin/callback" 0 "")
else
  LINKEDIN_CLIENT_ID=""
  LINKEDIN_CLIENT_SECRET=""
  LINKEDIN_REDIRECT_URL=""
fi

# X (Twitter) OAuth (only if enabled)
if is_enabled "$LOGIN_X_ENABLED"; then
  section "7. X (Twitter) OAuth"
  X_CLIENT_ID=$(prompt X_CLIENT_ID "X_CLIENT_ID" "" 0 "")
  X_CLIENT_SECRET=$(prompt X_CLIENT_SECRET "X_CLIENT_SECRET" "" 0 "")
  X_REDIRECT_URL=$(prompt X_REDIRECT_URL "X_REDIRECT_URL" "http://localhost:${PORT}/auth/x/callback" 0 "")
else
  X_CLIENT_ID=""
  X_CLIENT_SECRET=""
  X_REDIRECT_URL=""
fi

# R2
section "8. Cloudflare R2 (optional)"
note "Leave blank to skip file uploads."
CLOUDFLARE_ACCOUNT_ID=$(prompt CLOUDFLARE_ACCOUNT_ID "CLOUDFLARE_ACCOUNT_ID" "" 0 "")
CLOUDFLARE_ACCESS_KEY_ID=$(prompt CLOUDFLARE_ACCESS_KEY_ID "CLOUDFLARE_ACCESS_KEY_ID" "" 0 "")
CLOUDFLARE_SECRET_ACCESS_KEY=$(prompt CLOUDFLARE_SECRET_ACCESS_KEY "CLOUDFLARE_SECRET_ACCESS_KEY" "" 0 "")
CLOUDFLARE_R2_BUCKET=$(prompt CLOUDFLARE_R2_BUCKET "CLOUDFLARE_R2_BUCKET" "" 0 "")
CLOUDFLARE_R2_REGION=$(prompt CLOUDFLARE_R2_REGION "CLOUDFLARE_R2_REGION" "auto" 0 "")

# Resend
section "9. Resend Email (optional)"
note "Leave blank to skip welcome emails on registration."
RESEND_API_KEY=$(prompt RESEND_API_KEY "RESEND_API_KEY" "" 0 "")
RESEND_FROM=$(prompt RESEND_FROM "RESEND_FROM" 'Scaffold <onboarding@resend.dev>' 0 "")

# Write .env
section "Writing .env"
cat > "$ENV_FILE" << ENVEOF
# Generated by: make setup
# Go module path (used in go.mod and imports)
GO_MODULE_PATH=$GO_MODULE_PATH

# Database Configuration
DB_DSN=$DB_DSN

# Session Configuration
SESSION_SECRET=$SESSION_SECRET

# Server
PORT=$PORT
ADMIN_BASE_PATH=$ADMIN_BASE_PATH

# Login method toggles
LOGIN_PASSWORD_ENABLED=$LOGIN_PASSWORD_ENABLED
LOGIN_GOOGLE_ENABLED=$LOGIN_GOOGLE_ENABLED
LOGIN_GITHUB_ENABLED=$LOGIN_GITHUB_ENABLED
LOGIN_LINKEDIN_ENABLED=$LOGIN_LINKEDIN_ENABLED
LOGIN_X_ENABLED=$LOGIN_X_ENABLED

# Google OAuth Configuration
GOOGLE_CLIENT_ID=$GOOGLE_CLIENT_ID
GOOGLE_CLIENT_SECRET=$GOOGLE_CLIENT_SECRET
GOOGLE_REDIRECT_URL=$GOOGLE_REDIRECT_URL

# GitHub OAuth Configuration
GITHUB_CLIENT_ID=$GITHUB_CLIENT_ID
GITHUB_CLIENT_SECRET=$GITHUB_CLIENT_SECRET
GITHUB_REDIRECT_URL=$GITHUB_REDIRECT_URL

# LinkedIn OAuth Configuration
LINKEDIN_CLIENT_ID=$LINKEDIN_CLIENT_ID
LINKEDIN_CLIENT_SECRET=$LINKEDIN_CLIENT_SECRET
LINKEDIN_REDIRECT_URL=$LINKEDIN_REDIRECT_URL

# X (Twitter) OAuth Configuration
X_CLIENT_ID=$X_CLIENT_ID
X_CLIENT_SECRET=$X_CLIENT_SECRET
X_REDIRECT_URL=$X_REDIRECT_URL

# Cloudflare R2 Configuration
CLOUDFLARE_ACCOUNT_ID=$CLOUDFLARE_ACCOUNT_ID
CLOUDFLARE_ACCESS_KEY_ID=$CLOUDFLARE_ACCESS_KEY_ID
CLOUDFLARE_SECRET_ACCESS_KEY=$CLOUDFLARE_SECRET_ACCESS_KEY
CLOUDFLARE_R2_BUCKET=$CLOUDFLARE_R2_BUCKET
CLOUDFLARE_R2_REGION=$CLOUDFLARE_R2_REGION

# Resend Email
RESEND_API_KEY=$RESEND_API_KEY
RESEND_FROM=$RESEND_FROM
ENVEOF

# Create dirs
mkdir -p "$ROOT/bin" "$ROOT/log"

# Apply Go module path in go.mod and all .go files
if [ -n "$GO_MODULE_PATH" ]; then
  GO_MODULE_PATH="$GO_MODULE_PATH" bash "$ROOT/scripts/apply-module-path.sh" 2>/dev/null && echo -e "${GREEN}  ✓ Go module path applied${RESET}" || true
fi

echo ""
echo -e "${GREEN}${BOLD}  ✓ .env created at $ENV_FILE${RESET}"
echo -e "${GREEN}  ✓ bin/ and log/ directories ready${RESET}"
echo ""
echo -e "${DIM}  Next steps:${RESET}"
echo -e "    ${CYAN}make migrate${RESET}  - run database migrations"
echo -e "    ${CYAN}make dev${RESET}     - start the application"
echo ""
