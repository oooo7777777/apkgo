#!/usr/bin/env bash
set -euo pipefail

ADDR="${APKGO_WEB_ADDR:-127.0.0.1:8787}"
PORT="${ADDR##*:}"

if ! command -v lsof >/dev/null 2>&1; then
  echo "lsof not found, cannot stop apkgo web automatically."
  exit 1
fi

pids="$(lsof -nP -iTCP:"${PORT}" -sTCP:LISTEN -t 2>/dev/null || true)"
if [[ -z "${pids}" ]]; then
  echo "apkgo web is not running on ${ADDR}"
  exit 0
fi

kill ${pids} 2>/dev/null || true
echo "Stopped apkgo web on ${ADDR}"
