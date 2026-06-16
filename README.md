# apkgo

将 APK 发布到多个安卓应用商店的 CLI，同时提供本地 Web 控制台。

## 概览

- 适用于本地发布、CI/CD、自动化脚本
- 标准输出为 JSON，日志输出到 stderr
- 默认推荐使用 `config/config.json`
- 支持本地 Web 页面进行上传、审核查询和发布记录查看

## 支持的市场

`huawei`、`xiaomi`、`oppo`、`vivo`、`honor`、`tencent`、`googleplay`、`samsung`、`pgyer`、`fir`、`script`


# 本地 Web 控制台

推荐将本地 Web 作为日常发布入口。

```bash
apkgo web
apkgo web --addr 127.0.0.1:8787
```

macOS 用户可以直接双击：

```bash
启动APKGO网页.command
停止APKGO网页.command
```

`启动APKGO网页.command` 会自动识别当前 Mac 是 `Intel` 还是 `Apple Silicon`，并优先启动项目目录中的对应 `apkgo` 二进制。

Windows 用户请在 PowerShell 或终端中运行：

```powershell
apkgo web
```

如果二进制还没有加入 PATH，也可以直接运行下载或构建后的可执行文件。

开发源码场景下，如果需要自行构建本地二进制，可以在项目根目录运行：

```powershell
go build -o apkgo .
./apkgo web
```

`.command` 文件仅适用于 macOS，不适用于 Windows 或 Linux。

当前能力：

- 自动读取 `config/config.json`
- 支持上传 `1 个 zip` 或 `1..n 个 apk`
- zip 与 apk 不能混传
- zip 内会递归扫描所有层级目录中的 `.apk`
- 根据 `market_aliases` 文件名规则自动识别市场
- 单个 APK 未命中别名时，展示所有已配置市场供手动选择
- 仅展示已配置凭证的市场
- 支持发布任务、审核查询、发布记录查看与删除



# 命令行使用
## 安装

```bash
go install .
```

请在当前仓库根目录执行，这样安装的是本地源码版本。



## 快速开始

```bash
mkdir -p config
apkgo init -c config/config.json
vim config/config.json
apkgo web
```

默认情况下，程序优先读取 `config/config.json`。`apkgo.yaml` 和 `APKGO_*` 环境变量仅用于兼容旧流程。



## 常用命令

### 上传

```bash
apkgo upload -f app.apk
apkgo upload -f app.apk --store huawei,xiaomi
apkgo upload -f app.apk --notes "修复若干问题"
apkgo upload -f app.apk --notes-file CHANGELOG.md
apkgo upload -f app.apk --dry-run
apkgo upload -f app.apk --publish-mode auto
apkgo upload -f app.apk --publish-mode scheduled --publish-time "2026-06-20 10:00:00"
```

### 审核查询

```bash
apkgo audit -p com.example.app
apkgo audit -f app.apk -s tencent,huawei
apkgo audit -p com.example.app --watch --interval 1m -t 1h
```

### 配置检查

```bash
apkgo doctor
apkgo doctor -s huawei
apkgo doctor -s huawei -p com.example.app
apkgo doctor -s huawei -f app.apk
```

### 其他命令

```bash
apkgo init --store huawei,xiaomi -c config/config.json
apkgo stores
apkgo version
apkgo history
```

## 配置

推荐维护 `config/` 目录：

- `config/config.json`
- `config/config.example.json`
- `config/README.md`
- `config/huawei.json`
- `config/xiaomi.cer`

初始化：

```bash
cp config/config.example.json config/config.json
```

详细字段说明直接看：

- [config/README.md](/Users/wangwei/Documents/go/apkgo/config/README.md)

配置解析顺序：

1. 显式传入的 `--config`
2. `config/config.json`
3. `apkgo.yaml`
4. `APKGO_*` 环境变量覆盖同名字段

常见结构：

```json
{
  "hooks": {
    "after": "go run . notify feishu --webhook 'https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook'"
  },
  "market_aliases": {
    "tencent": ["tencent", "qq"],
    "pgyer": ["pgyer", "merit"],
    "xiaomi": ["xiaomi", "xm"]
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

如果你想逐字段看“这个 key 是干什么的、什么时候必填、应该填什么值”，请直接看 [config/README.md](/Users/wangwei/Documents/go/apkgo/config/README.md)。

`market_aliases` 用于按 APK 文件名识别市场，值是“市场 -> 别名数组”的映射。Web 上传会用正则做大小写不敏感匹配，只要文件名中包含别名就会命中，例如：

- `demo-apk-xiaomi-release.apk` -> `xiaomi`
- `demo_xm_release.apk` -> `xiaomi`
- `demo.qq.release.apk` -> `tencent`
- `demo-merit-release.apk` -> `pgyer`

Web 上传的补充规则：

- 命中别名时，按识别结果自动带出市场
- 单个 APK 未命中别名时，页面会展示所有已配置市场，交给用户手动选择
- 多个 APK 或 zip 中所有 APK 都未命中别名时，不会自动兜底到全部市场，避免误发

如果你在 `config/config.json` 里为某个市场配置了 `market_aliases`，它会覆盖该市场的默认别名数组。

## 隐私与本地文件

以下文件不应提交到 Git：

- `config/config.json`
- `config/huawei.json`
- `config/xiaomi.cer`
- `~/.apkgo/history.jsonl`

仓库中应保留的配置文件只有：

- `config/config.example.json`
- `config/README.md`

## 输出与退出码

`stdout` 为结构化 JSON，`stderr` 为日志。

退出码：

- `0`：全部成功
- `1`：部分失败
- `2`：全部失败
- `3`：输入或配置错误

## 项目结构

```text
cmd/           CLI 命令
pkg/apk/       APK 信息解析
pkg/config/    配置加载
pkg/history/   本地发布记录
pkg/store/     各市场实现
pkg/uploader/  上传编排与进度事件
```
### 凭证获取指南

| 商店 | 控制台地址 | 说明 |
|------|-----------|------|
| 华为 | [AppGallery Connect](https://developer.huawei.com/consumer/cn/console) | 用户与权限 > 服务账号（[详细步骤](#华为-appgallery-connect)） |
| 小米 | [小米开放平台](https://dev.mi.com) | 账号管理 > 接口密钥（[详细步骤](#小米开放平台)） |
| OPPO | [OPPO 开放平台](https://open.oppomobile.com) | 管理中心 > API 密钥管理（[详细步骤](#oppo-开放平台)） |
| vivo | [vivo 开放平台](https://dev.vivo.com.cn) | 账号管理 > API 接入（[详细步骤](#vivo-开放平台)） |
| 荣耀 | [荣耀开发者平台](https://developer.honor.com) | API 管理（[详细步骤](#荣耀开发者平台)） |
| 腾讯 | [腾讯开放平台](https://app.open.qq.com) | 应用 > 账户管理 > API 发布接口 > 申请开通（[详细步骤](#腾讯应用宝)） |
| 蒲公英 | [pgyer.com](https://www.pgyer.com/account/api) | 账户设置 > API 密钥（[详细步骤](#蒲公英-pgyer)） |
| fir.im | [betaqr.com.cn](https://www.betaqr.com.cn) | 账户 > API Token（[详细步骤](#firim)） |

每家的凭证申请流程都以官方文档为准（链接见下文），README 这边只描述 **apkgo 特有的事**：要哪几个字段、`doctor` 怎么验、需要注意的非显然行为。

#### 华为 AppGallery Connect

📖 官方文档：[Service Account 接入介绍](https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-Guides/agcapi-getstarted-0000001111845114#section1785535363715)

推荐用**开发者级**服务账号（PS256 JWT 鉴权），不要选项目级——访问发布 API 会被拒。下载到的 JSON 凭证文件直接交给 apkgo：

```yaml
stores:
  huawei:
    service_account_file: "/secure/path/huawei-sa.json"
    # 或 base64(JSON) 内联：service_account: "ewogICJrZXlfaWQiOiAi..."
```

旧版 `client_id` + `client_secret` 仍兼容，但华为已不推荐。

```bash
apkgo doctor -s huawei -p com.example.app
```

3 项探针：`token` / `appid-list`（包名 → appId）/ `release-permission`（应用发布权限）。

#### 小米开放平台

📖 官方文档：[API 上传应用](https://dev.mi.com/xiaomihyperos/documentation/detail?pId=1134)

要在小米后台「接口密钥」页面拿两样东西：**接口密钥**（SDK 里叫 password）和**公钥证书**（`.cer` 文件）。两个都是开发者账号绑定的。

```yaml
stores:
  xiaomi:
    email: "<开发者账号邮箱>"
    private_key: "<接口密钥>"
    cert_file: "/secure/path/xiaomi-pubkey.cer"
    # 也支持: cert: "-----BEGIN CERTIFICATE-----..." 或 base64(.cer)
```

```bash
apkgo doctor -s xiaomi -p com.example.app
```

> ⚠️ apkgo v3.0 之前内置了一份公钥证书，但那份 **2023-05 已过期**（且来源不明），从 v3.0 起必须自己提供。

#### OPPO 开放平台

📖 官方文档：[发布接口接入指引](https://open.oppomobile.com/new/developmentDoc/info?id=10998)

```yaml
stores:
  oppo:
    client_id: "<19 位数字>"
    client_secret: "<密钥>"
```

```bash
apkgo doctor -s oppo -p com.example.app
```

OPPO 的发布是异步任务，apkgo 会自动处理两个非显然的状态：撞 `911216 任务处理中` 时跳过 publish 直接等任务结束；撞 `911215 应用审核中` 视为成功（已进入审核队列）。

#### vivo 开放平台

📖 官方文档：[开放接口指引](https://dev.vivo.com.cn/documentCenter/doc/326)

```yaml
stores:
  vivo:
    access_key: "<...>"
    access_secret: "<...>"
```

```bash
apkgo doctor -s vivo -p com.example.app
```

vivo 的错误码分两层：网关 `code` + 业务 `subCode`。apkgo 同时识别两层，错误信息直接打印中文消息（比如 `[15042] 请上传与历史签名一致的APK包...`）。

#### 荣耀开发者平台

📖 官方文档：[发布接口指南](https://developer.honor.com/cn/doc/guides/101159)

```yaml
stores:
  honor:
    client_id: "<...>"
    client_secret: "<...>"
    # app_id: ""  # 可选，不填则按 APK 包名自动查
```

```bash
apkgo doctor -s honor -p com.example.app
```

doctor `app-detail` 探针会预检 *应用简介*（intro）—— 这个字段在荣耀后台必须填，否则 `update-language-info` 会以 `[20076] app introduction is empty` 拒绝。先在控制台填好再发版。

#### 腾讯应用宝

📖 官方文档：[API 接口传包-接入介绍](https://wikinew.open.qq.com/index.html#/iwiki/4015262492)

腾讯没有 list 或 pkg→id 反查接口，所以 `app_id` 必须手填。一份 yaml 服务多个应用用 `app_id_map`：

```yaml
stores:
  tencent:
    user_id: "<开发者 ID>"
    access_secret: "<接口密钥>"
    # 单 app:
    app_id: "<应用 ID>"
    # 多 app: 按 APK 包名命中
    # app_id_map: '{"com.example.foo":"111","com.example.bar":"222"}'
```

```bash
apkgo doctor -s tencent -p com.example.app
```

发布是异步任务，apkgo 会轮询 `query_app_update_status` 直到 `audit_status` 终态（最长 5 分钟）；超时视为成功（任务已交给腾讯）。

#### 蒲公英 (Pgyer)

📖 官方文档：[API 上传应用](https://www.pgyer.com/doc/view/app_upload)

```yaml
stores:
  pgyer:
    api_key: "<...>"
```

```bash
apkgo doctor -s pgyer -p com.example.app
```

#### fir.im

📖 官方文档：[betaqr.com.cn/docs](https://www.betaqr.com.cn/docs)

```yaml
stores:
  fir:
    api_token: "<...>"
```

```bash
apkgo doctor -s fir
```

> ⚠️ **fir 上传要求账号已完成实名认证**，否则 `/apps` 接口会以 `没有实名认证不能上传app` 拒绝。先去后台做实名再用。

## AI Agent 集成

apkgo 的输出格式专为 AI Agent 和自动化场景设计：

**结构化 JSON 输出** (stdout):
```json
{
  "apk": {"package": "com.example.app", "version_name": "1.0.0", "version_code": 1},
  "results": [
    {"store": "huawei", "success": true, "category": "success", "duration_ms": 12300},
    {"store": "oppo",   "success": true, "category": "already_done", "duration_ms": 3200},
    {"store": "xiaomi", "success": false, "category": "policy_block", "error": "签名不一致...", "duration_ms": 400}
  ]
}
```
