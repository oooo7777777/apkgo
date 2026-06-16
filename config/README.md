# config.json 字段说明

本文档说明 [config/config.json](/Users/wangwei/Documents/go/apkgo/config/config.json) 中已出现字段的用途与配置方式，适用于当前仓库默认配置文件。

## 配置文件位置

默认配置文件路径：

```text
config/config.json
```

以下命令在未显式传入 `--config` 时，默认读取该文件：

```bash
apkgo web
apkgo upload -f app.apk
```

## 顶层字段

### `hooks`

当前配置中包含以下字段：

```json
"hooks": {
  "after": "/usr/local/go/bin/go run . notify feishu --webhook 'https://open.feishu.cn/open-apis/bot/v2/hook/xxxx'"
}
```

字段说明：

- `hooks.after`
  所有市场上传结束后执行一次附加命令。

当前示例用途：

- 调用 `notify feishu`，将上传结果发送到飞书机器人。

执行语义：

- `hooks.after` 属于上传完成后的收尾动作。
- 该命令执行失败时，不会回滚已经完成的上传结果。

### `market_aliases`

用于 Web 页面按文件名识别市场。

```json
"market_aliases": {
  "fir": ["fir"],
  "googleplay": ["googleplay"],
  "honor": ["honor"],
  "huawei": ["huawei"],
  "oppo": ["oppo"],
  "pgyer": ["pgyer", "merit"],
  "samsung": ["samsung"],
  "script": ["script"],
  "tencent": ["tencent", "qq"],
  "vivo": ["vivo"],
  "xiaomi": ["xiaomi", "xm"]
}
```

字段说明：

- key：市场名
- value：该市场允许识别的文件名别名数组

匹配规则：

- Web 页面使用别名匹配 APK 文件名
- 文件名中包含某个别名时，即识别为对应市场

示例：

- 文件名包含 `xm` 时，识别为 `xiaomi`
- 文件名包含 `qq` 时，识别为 `tencent`
- 文件名包含 `merit` 时，识别为 `pgyer`

当前配置中的常用别名：

- `xiaomi: ["xiaomi", "xm"]`
- `tencent: ["tencent", "qq"]`
- `pgyer: ["pgyer", "merit"]`

### `ui`

用于本地 Web 页面的辅助配置。

#### `ui.default_audit_package`

```json
"default_audit_package": "uni.UNIE7FC6F0"
```

字段说明：

- Web 审核页默认带出的包名。

#### `ui.manual_urls`

```json
"manual_urls": {
  "huawei": "https://developer.huawei.com/consumer/cn/",
  "tencent": "https://open.tencent.com/",
  "oppo": "https://open.oppomobile.com/",
  "honor": "https://developer.honor.com/cn/?source=yingyongtuiguang0603",
  "vivo": "https://developer.vivo.com.cn/",
  "xiaomi": "https://dev.mi.com/xiaomihyperos",
  "pgyer": "https://www.pgyer.com/manager/dashboard/app/32925788dcc968b8ba3ec6e08c5b39b1"
}
```

字段说明：

- key：市场名
- value：对应市场的手动查看地址

用途：

- Web 页面查看审核状态时，可提供跳转到后台控制台的入口。

## 各市场字段

以下市场配置均已出现在当前 `config/config.json` 中。

### `pgyer`

```json
"pgyer": {
  "api_key": "..."
}
```

字段说明：

- `api_key`
  蒲公英上传接口的 API key。

### `fir`

```json
"fir": {
  "api_token": ""
}
```

字段说明：

- `api_token`
  fir.im 的 API Token。

补充说明：

- 空值通常表示未启用该市场配置。

### `tencent`

```json
"tencent": {
  "user_id": "...",
  "app_id": "...",
  "access_secret": "..."
}
```

字段说明：

- `user_id`
  应用宝开发者账号的用户 ID。
- `app_id`
  应用宝后台中对应应用的 app id。
- `access_secret`
  应用宝发布接口密钥。

### `huawei`

```json
"huawei": {
  "service_account_file": "./config/huawei.json"
}
```

字段说明：

- `service_account_file`
  华为服务账号 JSON 文件路径。

关联文件：

- [config/huawei.json](/Users/wangwei/Documents/go/apkgo/config/huawei.json)

### `vivo`

```json
"vivo": {
  "access_key": "...",
  "access_secret": "..."
}
```

字段说明：

- `access_key`
  vivo 开放平台 access key。
- `access_secret`
  vivo 开放平台 access secret。

### `honor`

```json
"honor": {
  "client_id": "...",
  "client_secret": "..."
}
```

字段说明：

- `client_id`
  荣耀开放平台 API 的 client id。
- `client_secret`
  荣耀开放平台 API 的 client secret。

### `oppo`

```json
"oppo": {
  "client_id": "...",
  "client_secret": "..."
}
```

字段说明：

- `client_id`
  OPPO 开放平台 API 的 client id。
- `client_secret`
  OPPO 开放平台 API 的 client secret。

### `xiaomi`

```json
"xiaomi": {
  "email": "product@merach.com",
  "private_key": "...",
  "cert_file": "./config/xiaomi.cer"
}
```

字段说明：

- `email`
  小米开发者账号邮箱。
- `private_key`
  小米上传接口使用的私钥。
- `cert_file`
  小米公钥证书文件路径。

关联文件：

- [config/xiaomi.cer](/Users/wangwei/Documents/go/apkgo/config/xiaomi.cer)

### `googleplay`

```json
"googleplay": {
  "json_key_file": "",
  "package_name": "",
  "track": ""
}
```

字段说明：

- `json_key_file`
  Google Play service account JSON 文件路径。
- `package_name`
  Android 包名。
- `track`
  发布轨道，例如 `production`、`beta`、`internal`。

补充说明：

- 空值通常表示未启用该市场配置。

### `samsung`

```json
"samsung": {
  "service_account_id": "",
  "private_key": "",
  "content_id": ""
}
```

字段说明：

- `service_account_id`
  Samsung Seller Portal 服务账号 ID。
- `private_key`
  Samsung 上传接口私钥。
- `content_id`
  Galaxy Store 中应用对应的 content id。

补充说明：

- 空值通常表示未启用该市场配置。

### `script`

```json
"script": {
  "command": ""
}
```

字段说明：

- `command`
  自定义脚本命令。

用途：

- 启用后，`apkgo` 会将上传参数作为 JSON 传递给该脚本。

补充说明：

- 空值通常表示未启用该市场配置。

## 当前配置中已启用的主要市场

按当前 `config/config.json` 内容，以下市场已配置有效值：

- `pgyer`
- `tencent`
- `huawei`
- `vivo`
- `honor`
- `oppo`
- `xiaomi`

以下市场当前为空，通常表示未启用：

- `fir`
- `googleplay`
- `samsung`
- `script`
