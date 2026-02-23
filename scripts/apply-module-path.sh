#!/usr/bin/env bash
# Apply GO_MODULE_PATH from .env to go.mod and all Go files. Safe to run multiple times.

set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

# Load GO_MODULE_PATH from .env if present (parse one var to avoid .env syntax issues)
if [ -z "$GO_MODULE_PATH" ] && [ -f .env ]; then
  GO_MODULE_PATH=$(grep '^GO_MODULE_PATH=' .env 2>/dev/null | cut -d= -f2- | sed "s/^['\"]//;s/['\"]$//")
fi

[ -n "$GO_MODULE_PATH" ] || { echo "GO_MODULE_PATH not set (add to .env or set in environment)."; exit 0; }

GO_MOD_FILE="$ROOT/go.mod"
[ -f "$GO_MOD_FILE" ] || { echo "go.mod not found."; exit 1; }

CURRENT_MODULE=$(sed -n '1s/^module //p' "$GO_MOD_FILE")
[ -n "$CURRENT_MODULE" ] || { echo "Could not read module from go.mod."; exit 1; }
[ "$CURRENT_MODULE" != "$GO_MODULE_PATH" ] || { echo "Module path already set to $GO_MODULE_PATH."; exit 0; }

echo "Applying Go module path: $GO_MODULE_PATH (replacing $CURRENT_MODULE)"
sed -i.bak "s|${CURRENT_MODULE}|${GO_MODULE_PATH}|g" "$GO_MOD_FILE" && rm -f "${GO_MOD_FILE}.bak"
while IFS= read -r -d '' f; do
  sed -i.bak "s|${CURRENT_MODULE}|${GO_MODULE_PATH}|g" "$f" && rm -f "${f}.bak"
done < <(find "$ROOT" -name "*.go" -not -path "*/vendor/*" -print0 2>/dev/null)
go mod tidy
echo "Done. Module path applied."
