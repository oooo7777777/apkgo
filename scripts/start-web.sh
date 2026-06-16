#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADDR="${APKGO_WEB_ADDR:-127.0.0.1:8787}"
URL="http://${ADDR}"
CONFIG_PATH="${ROOT_DIR}/config/config.json"
CONFIG_EXAMPLE_PATH="${ROOT_DIR}/config/config.example.json"
LOG_PATH="${APKGO_WEB_LOG:-/tmp/apkgo-web.log}"

cd "${ROOT_DIR}"

detect_macos_arch() {
  local machine
  machine="$(uname -m)"
  case "${machine}" in
    arm64|aarch64)
      echo "arm64"
      ;;
    x86_64|amd64)
      echo "amd64"
      ;;
    *)
      echo ""
      ;;
  esac
}

resolve_binary_path() {
  local arch="$1"
  local -a candidates=(
    "${ROOT_DIR}/apkgo"
    "${ROOT_DIR}/bin/apkgo"
    "${ROOT_DIR}/bin/apkgo-darwin-${arch}"
    "${ROOT_DIR}/bin/apkgo-${arch}"
    "${ROOT_DIR}/dist/apkgo_Darwin_${arch}/apkgo"
  )
  local candidate

  for candidate in "${candidates[@]}"; do
    if [[ -x "${candidate}" ]]; then
      echo "${candidate}"
      return 0
    fi
  done

  while IFS= read -r candidate; do
    if [[ -n "${candidate}" && -x "${candidate}" ]]; then
      echo "${candidate}"
      return 0
    fi
  done < <(find "${ROOT_DIR}/dist" -type f -name apkgo 2>/dev/null | sort | rg "/(darwin|Darwin).*(x86_64|amd64|arm64)/apkgo$|/(apkgo_)?(darwin|Darwin)_(x86_64|amd64|arm64).*/apkgo$" || true)

  return 1
}

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

ARCH="$(detect_macos_arch)"
if [[ -z "${ARCH}" ]]; then
  echo "Unsupported macOS architecture: $(uname -m)"
  exit 1
fi

if ! BIN_PATH="$(resolve_binary_path "${ARCH}")"; then
  echo "Missing apkgo binary for macOS ${ARCH}"
  echo "Expected one of these layouts:"
  echo "  ${ROOT_DIR}/apkgo"
  echo "  ${ROOT_DIR}/bin/apkgo"
  echo "  ${ROOT_DIR}/bin/apkgo-darwin-${ARCH}"
  echo "  ${ROOT_DIR}/dist/apkgo_Darwin_${ARCH}/apkgo"
  echo "Download the matching release package for this Mac and place the binary in the project folder."
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
echo "Using binary: ${BIN_PATH}"
rm -f "${LOG_PATH}"
nohup "${BIN_PATH}" web --addr "${ADDR}" >"${LOG_PATH}" 2>&1 &
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
