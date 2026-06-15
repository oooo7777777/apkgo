# apkgo

CLI tool for uploading APK files to Chinese Android app stores. All output is structured JSON on stdout; logs go to stderr.

## Install

```bash
go install github.com/KevinGong2013/apkgo@latest
# or download binary from https://github.com/KevinGong2013/apkgo/releases
```

## Commands

```bash
apkgo init [-s store1,store2] [-c config/config.json]   # Generate config file
apkgo upload -f <apk> [flags]                     # Upload APK to stores
apkgo doctor [-s stores] [-f apk | -p package]    # Diagnose store credentials/permissions
apkgo stores                                      # List stores and config schema (JSON)
apkgo version                                     # Version info (JSON)
```

## Upload flags

```
-f, --file         APK or AAB file path (required; .aab is googleplay-only)
    --file64       64-bit APK for split-arch uploads
-s, --store        Comma-separated store names (default: all configured)
-n, --notes        Release notes text
    --notes-file   Read release notes from file (overrides --notes)
    --dry-run      Validate without uploading
-t, --timeout      Global timeout (default: 10m)
-c, --config       Config file path (loader prefers config/config.json when omitted)
-o, --output       Output format: json or text (default: json)
```

## Supported stores

huawei, xiaomi, oppo, vivo, honor, tencent, googleplay, samsung, pgyer, fir, script

## Configuration

Preferred local config is `config/config.json`. `apkgo.yaml` and environment variables (`APKGO_<STORE>_<KEY>`) are still supported for legacy or CI-only flows.

```json
{
  "huawei": {
    "service_account": "",
    "service_account_file": "./config/huawei.json",
    "client_id": "",
    "client_secret": "",
    "app_id": ""
  },
  "xiaomi": {
    "email": "",
    "private_key": "",
    "cert": "",
    "cert_file": "./config/xiaomi.cer"
  },
  "oppo": {
    "client_id": "",
    "client_secret": ""
  },
  "vivo": {
    "access_key": "",
    "access_secret": ""
  },
  "honor": {
    "client_id": "",
    "client_secret": "",
    "app_id": ""
  },
  "tencent": {
    "user_id": "",
    "access_secret": "",
    "app_id": "",
    "app_id_map": "",
    "package_name": ""
  },
  "script": {
    "command": "./deploy.sh"
  },
  "script.cdn-upload": {
    "command": "./upload-cdn.sh"
  },
  "script.dingtalk": {
    "command": "./notify-dingtalk.sh"
  }
}
```

Env var example: `APKGO_HUAWEI_SERVICE_ACCOUNT=$(base64 -w0 huawei-sa.json) apkgo upload -f app.apk --store huawei`

## Hooks

Shell commands executed before/after uploads. Receive context as JSON on stdin.

### Configuration

```json
{
  "hooks": {
    "before": "./scripts/before-all.sh",
    "after": "./scripts/after-all.sh"
  },
  "huawei": {
    "client_id": "...",
    "before": "./scripts/before-huawei.sh",
    "after": "./scripts/after-huawei.sh"
  }
}
```

### Protocol

**Exit codes:**
- `0` — success (continue)
- non-zero — failure (`before` hooks abort the upload; `after` hooks log warning only)

**Environment variables** (set automatically):
- `APKGO_STORE` — store name (empty for global hooks)
- `APKGO_PACKAGE` — package name (e.g. `com.example.app`)
- `APKGO_VERSION` — version name (e.g. `1.2.0`)

**Errors:** stderr is captured as the error message.

### Stdin JSON schemas

**Global before** (`hooks.before`):
```json
{
  "file_path": "/path/to/app.apk",
  "apk": {"package": "com.example.app", "version_name": "1.0.0", "version_code": 1, "app_name": "MyApp"},
  "stores": ["huawei", "xiaomi"]
}
```

**Global after** (`hooks.after`):
```json
{
  "file_path": "/path/to/app.apk",
  "apk": {"package": "com.example.app", "version_name": "1.0.0", "version_code": 1, "app_name": "MyApp"},
  "results": [
    {"store": "huawei", "success": true, "duration_ms": 12300},
    {"store": "xiaomi", "success": false, "error": "auth failed", "duration_ms": 400}
  ]
}
```

**Per-store before** (`stores.<name>.before`):
```json
{
  "file_path": "/path/to/app.apk",
  "apk": {"package": "com.example.app", "version_name": "1.0.0", "version_code": 1, "app_name": "MyApp"},
  "store": "huawei"
}
```

**Per-store after** (`stores.<name>.after`):
```json
{
  "file_path": "/path/to/app.apk",
  "apk": {"package": "com.example.app", "version_name": "1.0.0", "version_code": 1, "app_name": "MyApp"},
  "store": "huawei",
  "result": {"store": "huawei", "success": true, "duration_ms": 12300}
}
```

## Output format

stdout is always parseable JSON:

```json
{
  "apk": {"package": "com.example", "version_name": "1.0.0", "version_code": 1, "app_name": "MyApp"},
  "results": [
    {"store": "huawei", "success": true, "duration_ms": 12300},
    {"store": "xiaomi", "success": false, "error": "auth: invalid private key", "duration_ms": 400}
  ]
}
```

## Exit codes

- **0**: All uploads succeeded
- **1**: Some uploads failed (partial success)
- **2**: All uploads failed
- **3**: Input/config error

## Typical agent workflow

```bash
# 1. Check if apkgo is installed
which apkgo

# 2. Generate config for needed stores
mkdir -p config
apkgo init --store huawei,xiaomi -c config/config.json

# 3. Discover required config fields
apkgo stores

# 4. Dry-run to validate
apkgo upload -f app.apk --dry-run

# 5. Upload
apkgo upload -f app.apk --notes "v1.0.0 release" --timeout 15m

# 6. Parse JSON result from stdout, check exit code
```

## Project structure

```
cmd/           CLI commands (cobra)
pkg/store/     Store interface + implementations (self-registering via init())
pkg/config/    JSON config + YAML/env fallback loading
pkg/apk/       APK metadata parser
pkg/uploader/  Concurrent upload orchestrator
```

Adding a new store: create `pkg/store/<name>/<name>.go`, implement `store.Store` interface, call `store.Register()` in `init()`. Zero changes to existing code.
