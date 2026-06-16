#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADDR="${APKGO_WEB_ADDR:-127.0.0.1:8787}"
URL="http://${ADDR}"
CONFIG_PATH="${ROOT_DIR}/config/config.json"
CONFIG_EXAMPLE_PATH="${ROOT_DIR}/config/config.example.json"
LOG_PATH="${APKGO_WEB_LOG:-/tmp/apkgo-web.log}"

cd "${ROOT_DIR}"

if [[ ! -f "${CONFIG_PATH}" ]]; then
  if [[ -f "${CONFIG_EXAMPLE_PATH}" ]]; then
    cp "${CONFIG_EXAMPLE_PATH}" "${CONFIG_PATH}"
    echo "Created config from template: ${CONFIG_PATH}"
  else
    echo "Missing config: ${CONFIG_PATH}"
    echo "Template not found: ${CONFIG_EXAMPLE_PATH}"
    echo "Read config instructions in: ${ROOT_DIR}/config/README.md"
    exit 1
  fi
  echo "Starting with template config. Fill in credentials in: ${CONFIG_PATH}"
  echo "Read config instructions in: ${ROOT_DIR}/config/README.md"
fi

if ! command -v go >/dev/null 2>&1; then
  echo "Missing dependency: go"
  echo "Install Go first, then retry:"
  echo "  brew install go"
  echo "  or download from https://go.dev/dl/"
  exit 1
fi

if command -v lsof >/dev/null 2>&1; then
  pids="$(lsof -nP -iTCP:"${ADDR##*:}" -sTCP:LISTEN -t 2>/dev/null || true)"
  if [[ -n "${pids}" ]]; then
    kill ${pids} 2>/dev/null || true
    sleep 0.3
  fi
fi

echo "Starting apkgo web at ${URL}"
rm -f "${LOG_PATH}"
nohup go run . web --addr "${ADDR}" >"${LOG_PATH}" 2>&1 &
pid=$!

for _ in {1..30}; do
  if curl -sf "${URL}" >/dev/null 2>&1; then
    open "${URL}"
    echo "apkgo web started (pid: ${pid})"
    echo "log: ${LOG_PATH}"
    exit 0
  fi
  if ! kill -0 "${pid}" 2>/dev/null; then
    break
  fi
  sleep 0.5
done

echo "Failed to start apkgo web at ${URL}"
echo "log: ${LOG_PATH}"
if [[ -f "${LOG_PATH}" ]]; then
  echo "--- recent log ---"
  tail -n 40 "${LOG_PATH}"
fi
exit 1
