#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${1:-$ROOT_DIR/.env}"
DAEMON_BIN="$ROOT_DIR/bin/daemon"

if [ ! -f "$ENV_FILE" ]; then
  echo "[ERROR] Env file not found: $ENV_FILE"
  echo "Usage: $0 [path-to-env-file]"
  exit 1
fi

if [ ! -x "$DAEMON_BIN" ]; then
  echo "[ERROR] daemon binary not found: $DAEMON_BIN"
  echo "Run: ./scripts/build.sh build"
  exit 1
fi

# shellcheck disable=SC1090
set -a
source "$ENV_FILE"
set +a

# Authorization and Region are request-level configuration: the caller sends
# them in every HTTP request header, so UAPI_AUTHORIZATION is no longer
# required at startup. UAPI_ENDPOINT only serves as the default endpoint
# prefix (used when no region is specified).
if [ -n "${UAPI_AUTHORIZATION:-}" ]; then
  echo "[WARN] UAPI_AUTHORIZATION is now supplied per request by callers; the startup-time env var is ignored."
fi

UAPI_ENV="${UAPI_ENV:-test}"
PORT="${PORT:-8080}"

echo "[INFO] Starting daemon with env file: $ENV_FILE"

exec "$DAEMON_BIN" \
  -env "$UAPI_ENV" \
  -port "$PORT"
