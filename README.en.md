# apkgo

CLI for publishing APKs to multiple Android app stores, with a local Web console.

## Overview

- Suitable for local release workflows, CI/CD, and automation
- Structured JSON on stdout, logs on stderr
- `config/config.json` is the recommended default config
- Includes a local Web UI for upload, review lookup, and release history

## Supported stores

`huawei`, `xiaomi`, `oppo`, `vivo`, `honor`, `tencent`, `googleplay`, `samsung`, `pgyer`, `fir`, `script`

## Install

```bash
go install .
```

Run it from the repository root to install the local version.

If you need the remote repository version instead, use:

```bash
go install github.com/KevinGong2013/apkgo@latest
```

## Quick start

```bash
mkdir -p config
apkgo init -c config/config.json
vim config/config.json
apkgo web
```

By default, apkgo prefers `config/config.json`. `apkgo.yaml` and `APKGO_*` environment variables are kept only for legacy compatibility.

## Local Web console

The local Web UI is the recommended daily entry point.

```bash
apkgo web
apkgo web --addr 127.0.0.1:8787
```

On macOS, you can also double-click:

```bash
启动APKGO网页.command
停止APKGO网页.command
```

On Windows, run in PowerShell or Terminal:

```powershell
apkgo web
```

If the binary is not in PATH yet, run from the project root:

```powershell
go run . web
```

`.command` files are macOS-only and do not apply to Windows or Linux.

Current capabilities:

- auto-loads `config/config.json`
- accepts either `1 zip` or `1..n apk` files
- does not allow mixing zip and apk in one upload
- recursively scans nested directories inside zip files
- detects target stores from file-name rules
- only shows stores with configured credentials
- supports publish tasks, review lookup, and release history view/delete

## Common commands

### Upload

```bash
apkgo upload -f app.apk
apkgo upload -f app.apk --store huawei,xiaomi
apkgo upload -f app.apk --notes "Bug fixes"
apkgo upload -f app.apk --notes-file CHANGELOG.md
apkgo upload -f app.apk --dry-run
apkgo upload -f app.apk --publish-mode auto
apkgo upload -f app.apk --publish-mode scheduled --publish-time "2026-06-20 10:00:00"
```

### Review lookup

```bash
apkgo audit -p com.example.app
apkgo audit -f app.apk -s tencent,huawei
apkgo audit -p com.example.app --watch --interval 1m -t 1h
```

### Credential checks

```bash
apkgo doctor
apkgo doctor -s huawei
apkgo doctor -s huawei -p com.example.app
apkgo doctor -s huawei -f app.apk
```

### Other commands

```bash
apkgo init --store huawei,xiaomi -c config/config.json
apkgo stores
apkgo version
apkgo history
```

## Configuration

Recommended files in `config/`:

- `config/config.json`
- `config/config.example.json`
- `config/huawei.json`
- `config/xiaomi.cer`

Initialize from the template:

```bash
cp config/config.example.json config/config.json
```

Config resolution order:

1. explicit `--config`
2. `config/config.json`
3. `apkgo.yaml`
4. `APKGO_*` environment variables override matching fields

Common structure:

```json
{
  "hooks": {
    "after": "go run . notify feishu --webhook 'https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook'"
  },
  "ui": {
    "default_audit_package": "com.example.app",
    "manual_urls": {
      "huawei": "https://developer.huawei.com/consumer/cn/"
    }
  },
  "huawei": {
    "service_account_file": "./config/huawei.json"
  },
  "xiaomi": {
    "email": "product@example.com",
    "private_key": "your-private-key",
    "cert_file": "./config/xiaomi.cer"
  }
}
```

See the full guide in [config/README.md](/Users/wangwei/Documents/go/apkgo/config/README.md).

## Private files

Do not commit these files:

- `config/config.json`
- `config/huawei.json`
- `config/xiaomi.cer`
- `~/.apkgo/history.jsonl`

Only these config files should remain in the repository:

- `config/config.example.json`
- `config/README.md`

## Output and exit codes

`stdout` is structured JSON. `stderr` is used for logs.

Exit codes:

- `0`: all succeeded
- `1`: partial failure
- `2`: all failed
- `3`: input or config error

## Project structure

```text
cmd/           CLI commands
pkg/apk/       APK metadata parsing
pkg/config/    Config loading
pkg/history/   Local release history
pkg/store/     Store implementations
pkg/uploader/  Upload orchestration and progress events
```
