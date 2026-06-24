package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KevinGong2013/apkgo/v3/pkg/apk"
	"github.com/KevinGong2013/apkgo/v3/pkg/config"
	"github.com/KevinGong2013/apkgo/v3/pkg/history"
)

func TestDetectChannelFromName_UsesConfigAliases(t *testing.T) {
	cfg := &config.Config{}
	got, ok := detectChannelFromName(cfg, "demo_xm_release.apk")
	if !ok {
		t.Fatalf("expected xiaomi alias to match")
	}
	if got.Store != "xiaomi" || got.Channel != "xm" {
		t.Fatalf("got %#v, want xiaomi/xm", got)
	}
}

func TestDetectChannelFromName_UsesOverrideAliases(t *testing.T) {
	cfg := &config.Config{
		MarketAliases: map[string][]string{
			"xiaomi": {"mi"},
		},
	}
	if _, ok := detectChannelFromName(cfg, "demo_xm_release.apk"); ok {
		t.Fatalf("xm should not match after override")
	}
	got, ok := detectChannelFromName(cfg, "demo.mi.release.apk")
	if !ok {
		t.Fatalf("expected mi alias to match")
	}
	if got.Store != "xiaomi" || got.Channel != "mi" {
		t.Fatalf("got %#v, want xiaomi/mi", got)
	}
}

func TestBuildAPKBundle_ErrorsForMultipleUnmatchedFiles(t *testing.T) {
	tmp := t.TempDir()
	apk1 := filepath.Join(tmp, "demo-release.apk")
	apk2 := filepath.Join(tmp, "other-release.apk")
	if err := os.WriteFile(apk1, []byte("apk1"), 0644); err != nil {
		t.Fatalf("WriteFile apk1: %v", err)
	}
	if err := os.WriteFile(apk2, []byte("apk2"), 0644); err != nil {
		t.Fatalf("WriteFile apk2: %v", err)
	}

	cfg := &config.Config{
		Stores: map[string]map[string]string{
			"huawei": {"client_id": "x"},
		},
	}
	_, err := buildAPKBundle(cfg, []webUploadedFile{
		{Name: "demo-release.apk", Path: apk1},
		{Name: "other-release.apk", Path: apk2},
	})
	if err == nil {
		t.Fatalf("expected error for multiple unmatched files")
	}
}

func TestBuildAPKBundle_RejectsInvalidSingleAPK(t *testing.T) {
	tmp := t.TempDir()
	apkPath := filepath.Join(tmp, "demo.apk")
	if err := os.WriteFile(apkPath, []byte("not-a-real-apk"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg := &config.Config{
		Stores: map[string]map[string]string{
			"pgyer": {"api_key": "demo"},
		},
	}
	_, err := buildAPKBundle(cfg, []webUploadedFile{{Name: "demo.apk", Path: apkPath}})
	if err == nil {
		t.Fatalf("expected invalid apk error")
	}
	if !strings.Contains(err.Error(), "上传的 APK 文件无效") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSaveUploadedFiles_PreservesSingleAPKBytes(t *testing.T) {
	tmp := t.TempDir()
	apkPath := filepath.Join(tmp, "demo.apk")
	apkBytes := []byte("PK\x03\x04fake-apk-content")
	if err := os.WriteFile(apkPath, apkBytes, 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("archive", filepath.Base(apkPath))
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	src, err := os.Open(apkPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer src.Close()
	if _, err := io.Copy(part, src); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(webUploadFormMemory); err != nil {
		t.Fatalf("ParseMultipartForm: %v", err)
	}

	files, cleanup, err := saveUploadedFiles(req, "archive")
	if err != nil {
		t.Fatalf("saveUploadedFiles: %v", err)
	}
	defer cleanup()

	if len(files) != 1 {
		t.Fatalf("files len = %d, want 1", len(files))
	}
	got, err := os.ReadFile(files[0].Path)
	if err != nil {
		t.Fatalf("ReadFile saved temp: %v", err)
	}
	if !bytes.Equal(got, apkBytes) {
		t.Fatalf("saved bytes changed: got %q want %q", string(got), string(apkBytes))
	}
}

func TestBuildManualSelectionBundle_PreservesSingleAPKBytes(t *testing.T) {
	tmp := t.TempDir()
	apkPath := filepath.Join(tmp, "demo.apk")
	apkBytes := []byte("PK\x03\x04fake-apk-content")
	if err := os.WriteFile(apkPath, apkBytes, 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg := &config.Config{
		Stores: map[string]map[string]string{
			"pgyer": {"api_key": "demo"},
		},
	}
	dir := filepath.Join(tmp, "bundle")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	bundle, err := buildManualSelectionBundle(dir, func() {}, webUploadedFile{
		Name: "demo.apk",
		Path: apkPath,
	}, cfg)
	if err != nil {
		t.Fatalf("buildManualSelectionBundle: %v", err)
	}

	ref, ok := bundle.Artifacts["pgyer"]
	if !ok {
		t.Fatalf("expected pgyer artifact")
	}
	got, err := os.ReadFile(ref.Path)
	if err != nil {
		t.Fatalf("ReadFile copied apk: %v", err)
	}
	if !bytes.Equal(got, apkBytes) {
		t.Fatalf("copied bytes changed: got %q want %q", string(got), string(apkBytes))
	}
}

func TestHandleWebConfigSave_PreservesExistingTopLevelFields(t *testing.T) {
	tmp := t.TempDir()
	oldConfig := flagConfig
	flagConfig = filepath.Join(tmp, "config.json")
	defer func() { flagConfig = oldConfig }()

	initial := map[string]any{
		"market_aliases": map[string][]string{
			"xiaomi": {"xiaomi", "xm"},
		},
		"ui": map[string]any{
			"default_audit_package": "com.old.app",
			"manual_urls": map[string]string{
				"huawei": "https://developer.huawei.com/consumer/cn/",
			},
		},
		"hooks": map[string]string{
			"after": "go run . notify feishu --webhook 'https://old.example/hook'",
		},
		"huawei": map[string]string{
			"service_account_file": "./config/huawei.json",
		},
		"custom_store": map[string]string{
			"token": "keep-me",
		},
	}
	data, err := json.Marshal(initial)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if err := os.WriteFile(flagConfig, data, 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	payload := webConfigPayload{
		UI: webUIConfig{
			DefaultAuditPackage: "com.new.app",
		},
		Stores: map[string]map[string]string{
			"huawei": {
				"service_account_file": "./config/new-huawei.json",
			},
		},
	}
	payload.Hooks.FeishuWebhook = "https://open.feishu.cn/open-apis/bot/v2/hook/new"
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/config/save", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleWebConfigSave(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	saved, err := os.ReadFile(flagConfig)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var out map[string]json.RawMessage
	if err := json.Unmarshal(saved, &out); err != nil {
		t.Fatalf("Unmarshal saved: %v", err)
	}

	var marketAliases map[string][]string
	if err := json.Unmarshal(out["market_aliases"], &marketAliases); err != nil {
		t.Fatalf("Unmarshal market_aliases: %v", err)
	}
	if len(marketAliases["xiaomi"]) != 2 {
		t.Fatalf("market_aliases lost: %#v", marketAliases)
	}

	var custom map[string]string
	if err := json.Unmarshal(out["custom_store"], &custom); err != nil {
		t.Fatalf("Unmarshal custom_store: %v", err)
	}
	if custom["token"] != "keep-me" {
		t.Fatalf("custom_store token = %q, want keep-me", custom["token"])
	}

	var hooks map[string]string
	if err := json.Unmarshal(out["hooks"], &hooks); err != nil {
		t.Fatalf("Unmarshal hooks: %v", err)
	}
	if hooks["after"] == "" || hooks["after"] == initial["hooks"].(map[string]string)["after"] {
		t.Fatalf("hooks.after not updated: %#v", hooks)
	}

	var ui map[string]any
	if err := json.Unmarshal(out["ui"], &ui); err != nil {
		t.Fatalf("Unmarshal ui: %v", err)
	}
	if ui["default_audit_package"] != "com.new.app" {
		t.Fatalf("ui.default_audit_package = %#v, want com.new.app", ui["default_audit_package"])
	}
}

func TestHandleWebConfig_UsesConfigFieldDefinitions(t *testing.T) {
	tmp := t.TempDir()
	oldConfig := flagConfig
	flagConfig = filepath.Join(tmp, "config.json")
	defer func() { flagConfig = oldConfig }()

	if err := os.WriteFile(flagConfig, []byte(`{"xiaomi":{"email":"demo@example.com","private_key":"abc","cert":"pem-cert"}}`), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	w := httptest.NewRecorder()
	handleWebConfig(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	var resp struct {
		Sections []webConfigSection `json:"sections"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unmarshal response: %v", err)
	}

	var xiaomi *webConfigSection
	for i := range resp.Sections {
		if resp.Sections[i].Key == "xiaomi" {
			xiaomi = &resp.Sections[i]
			break
		}
	}
	if xiaomi == nil {
		t.Fatalf("missing xiaomi section")
	}

	fields := map[string]webConfigField{}
	for _, field := range xiaomi.Fields {
		fields[field.Key] = field
	}
	if _, ok := fields["cert_file"]; !ok {
		t.Fatalf("xiaomi cert_file field missing: %#v", xiaomi.Fields)
	}
	if _, ok := fields["cert"]; ok {
		t.Fatalf("xiaomi cert should not be exposed in config modal: %#v", xiaomi.Fields)
	}
	if !fields["private_key"].Secret {
		t.Fatalf("xiaomi private_key should be secret")
	}
}

func TestBuildWebConfigListItems_IncludesDefaultAuditPackage(t *testing.T) {
	doc := &webConfigDocument{
		Hooks: map[string]map[string]string{},
		UI: webUIConfig{
			DefaultAuditPackage: "uni.UNIE7FC6F0",
		},
		Stores:        map[string]map[string]string{},
		MarketAliases: map[string][]string{},
	}

	items := buildWebConfigListItems(doc)
	for _, item := range items {
		if item.GroupKey == "ui" && item.Key == "default_audit_package" {
			if item.Summary != "uni.UNIE7FC6F0" {
				t.Fatalf("summary = %q, want uni.UNIE7FC6F0", item.Summary)
			}
			if !item.Configured {
				t.Fatalf("expected default_audit_package to be configured")
			}
			return
		}
	}
	t.Fatalf("default_audit_package item not found")
}

func TestHandleWebConfigSave_PersistsDefaultAuditPackageOnly(t *testing.T) {
	tmp := t.TempDir()
	oldConfig := flagConfig
	flagConfig = filepath.Join(tmp, "config.json")
	defer func() { flagConfig = oldConfig }()

	payload := webConfigPayload{
		UI: webUIConfig{
			DefaultAuditPackage: "uni.UNIE7FC6F0",
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/config/save", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleWebConfigSave(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	saved, err := os.ReadFile(flagConfig)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var out map[string]json.RawMessage
	if err := json.Unmarshal(saved, &out); err != nil {
		t.Fatalf("Unmarshal saved: %v", err)
	}

	var ui map[string]any
	if err := json.Unmarshal(out["ui"], &ui); err != nil {
		t.Fatalf("Unmarshal ui: %v", err)
	}
	if ui["default_audit_package"] != "uni.UNIE7FC6F0" {
		t.Fatalf("ui.default_audit_package = %#v, want uni.UNIE7FC6F0", ui["default_audit_package"])
	}
}

func TestBuildWebConfigListItems_TreatsTemplateValuesAsUnconfigured(t *testing.T) {
	doc := &webConfigDocument{
		Hooks: map[string]map[string]string{
			"after": {
				"command": "go run . notify feishu --webhook 'https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook'",
			},
		},
		UI: webUIConfig{
			DefaultAuditPackage: "com.example.app",
		},
		Stores: map[string]map[string]string{
			"huawei": {
				"service_account_file": "./config/huawei.json",
			},
			"xiaomi": {
				"cert_file": "./config/xiaomi.cer",
			},
		},
		MarketAliases: map[string][]string{},
	}

	items := buildWebConfigListItems(doc)
	for _, item := range items {
		switch {
		case item.GroupKey == "stores" && item.Key == "huawei" && item.Configured:
			t.Fatalf("template huawei should not be configured: %#v", item)
		case item.GroupKey == "stores" && item.Key == "xiaomi" && item.Configured:
			t.Fatalf("template xiaomi should not be configured: %#v", item)
		case item.GroupKey == "ui" && item.Key == "default_audit_package" && item.Configured:
			t.Fatalf("template default_audit_package should not be configured: %#v", item)
		case item.GroupKey == "hooks" && item.Key == "feishu" && item.Configured:
			t.Fatalf("template feishu webhook should not be configured: %#v", item)
		}
	}
}

func TestHasConfiguredValues_FilePathRequiresActualFile(t *testing.T) {
	tmp := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	if err := os.MkdirAll("config", 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	if hasConfiguredValues(map[string]string{"service_account_file": "./config/huawei.json"}) {
		t.Fatalf("missing uploaded file should not count as configured")
	}

	if err := os.WriteFile(filepath.Join("config", "huawei.json"), []byte(`{"key_id":"demo"}`), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if !hasConfiguredValues(map[string]string{"service_account_file": "./config/huawei.json"}) {
		t.Fatalf("existing uploaded file should count as configured")
	}
}

func TestHandleWebConfigSave_DoesNotPreserveClearedSecretFields(t *testing.T) {
	tmp := t.TempDir()
	oldConfig := flagConfig
	flagConfig = filepath.Join(tmp, "config.json")
	defer func() { flagConfig = oldConfig }()

	initial := map[string]any{
		"vivo": map[string]string{
			"access_key":    "demo-key",
			"access_secret": "old-secret",
		},
	}
	data, err := json.Marshal(initial)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if err := os.WriteFile(flagConfig, data, 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	payload := webConfigPayload{
		Stores: map[string]map[string]string{
			"vivo": {
				"access_key":    "demo-key",
				"access_secret": "",
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/config/save", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleWebConfigSave(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	saved, err := os.ReadFile(flagConfig)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var out map[string]json.RawMessage
	if err := json.Unmarshal(saved, &out); err != nil {
		t.Fatalf("Unmarshal saved: %v", err)
	}
	var vivo map[string]string
	if err := json.Unmarshal(out["vivo"], &vivo); err != nil {
		t.Fatalf("Unmarshal vivo: %v", err)
	}
	if got := vivo["access_secret"]; got != "" {
		t.Fatalf("access_secret = %q, want empty", got)
	}
}

func TestHandleWebConfigSave_ReturnsValidationErrorForModal(t *testing.T) {
	tmp := t.TempDir()
	oldConfig := flagConfig
	oldCredsFrom := flagCredsFrom
	flagConfig = filepath.Join(tmp, "config.json")
	flagCredsFrom = ""
	defer func() {
		flagConfig = oldConfig
		flagCredsFrom = oldCredsFrom
	}()

	payload := webConfigPayload{
		Stores: map[string]map[string]string{
			"vivo": {
				"access_key":    "",
				"access_secret": "",
			},
		},
		TargetGroup:   "stores",
		TargetSection: "vivo",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/config/save", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleWebConfigSave(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unmarshal response: %v", err)
	}
	if resp["error"] == "" {
		t.Fatalf("expected validation error, got %#v", resp)
	}
}

func TestBuildWebConfigListItems_SortsConfiguredFirst(t *testing.T) {
	doc := &webConfigDocument{
		Hooks: map[string]map[string]string{},
		Stores: map[string]map[string]string{
			"vivo":   {"access_key": "demo", "access_secret": "demo"},
			"xiaomi": {},
		},
		MarketAliases: map[string][]string{},
	}

	items := buildWebConfigListItems(doc)
	var storeItems []webConfigListItem
	for _, item := range items {
		if item.GroupKey == "stores" {
			storeItems = append(storeItems, item)
		}
	}
	if len(storeItems) < 2 {
		t.Fatalf("expected at least 2 store items, got %d", len(storeItems))
	}
	foundUnconfigured := false
	for _, item := range storeItems {
		if !item.Configured {
			foundUnconfigured = true
			continue
		}
		if foundUnconfigured {
			t.Fatalf("configured item appears after unconfigured item: %#v", item)
		}
	}
}

func TestHandleWebConfigSave_UploadsStoreFileIntoConfigDir(t *testing.T) {
	tmp := t.TempDir()
	configDir := filepath.Join(tmp, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	oldConfig := flagConfig
	oldCredsFrom := flagCredsFrom
	flagConfig = filepath.Join(configDir, "config.json")
	flagCredsFrom = ""
	defer func() {
		flagConfig = oldConfig
		flagCredsFrom = oldCredsFrom
	}()

	if err := os.WriteFile(flagConfig, []byte(`{"huawei":{"service_account_file":"./config/old.json"}}`), 0644); err != nil {
		t.Fatalf("WriteFile config: %v", err)
	}

	payload := `{"ui":{"default_audit_package":""},"stores":{"huawei":{"service_account_file":"./config/huawei.json"}},"market_aliases":{},"target_group":"","target_section":"","hooks":{"feishu_webhook":""}}`
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("payload", payload); err != nil {
		t.Fatalf("WriteField: %v", err)
	}
	part, err := writer.CreateFormFile("store_file_huawei__service_account_file", "sa.json")
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := part.Write([]byte(`{"key_id":"demo"}`)); err != nil {
		t.Fatalf("part.Write: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/config/save", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	handleWebConfigSave(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	savedConfig, err := os.ReadFile(flagConfig)
	if err != nil {
		t.Fatalf("ReadFile config: %v", err)
	}
	if !strings.Contains(string(savedConfig), `"service_account_file": "./config/huawei.json"`) {
		t.Fatalf("config not updated: %s", string(savedConfig))
	}

	savedFile, err := os.ReadFile(filepath.Join(configDir, "huawei.json"))
	if err != nil {
		t.Fatalf("ReadFile uploaded file: %v", err)
	}
	if string(savedFile) != `{"key_id":"demo"}` {
		t.Fatalf("uploaded file content mismatch: %s", string(savedFile))
	}
}

func TestHandleWebConfigSave_ReloadsAndUsesSavedStoreConfig(t *testing.T) {
	tmp := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(config.DefaultJSONKeysPath), 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	oldConfig := flagConfig
	flagConfig = config.DefaultJSONKeysPath
	defer func() { flagConfig = oldConfig }()

	payload := webConfigPayload{
		Stores: map[string]map[string]string{
			"vivo": {
				"access_key":    "demo-key",
				"access_secret": "demo-secret",
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/config/save", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleWebConfigSave(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	cfg, err := loadWebConfig()
	if err != nil {
		t.Fatalf("loadWebConfig: %v", err)
	}
	if got := cfg.Stores["vivo"]["access_key"]; got != "demo-key" {
		t.Fatalf("loadWebConfig vivo.access_key = %q, want demo-key", got)
	}
	if got := cfg.Stores["vivo"]["access_secret"]; got != "demo-secret" {
		t.Fatalf("loadWebConfig vivo.access_secret = %q, want demo-secret", got)
	}

	runtimeCfg, err := loadWebRuntimeConfig()
	if err != nil {
		t.Fatalf("loadWebRuntimeConfig: %v", err)
	}
	if got := runtimeCfg.Stores["vivo"]["access_key"]; got != "demo-key" {
		t.Fatalf("loadWebRuntimeConfig vivo.access_key = %q, want demo-key", got)
	}

	visible := visibleWebStores(runtimeCfg)
	found := false
	for _, item := range visible {
		if item.Key == "vivo" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("vivo should be visible after config save, got %#v", visible)
	}
}

func TestHandleWebConfigSave_UploadedFileIsUsedByRuntimeConfig(t *testing.T) {
	tmp := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(config.DefaultJSONKeysPath), 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	oldConfig := flagConfig
	oldCredsFrom := flagCredsFrom
	flagConfig = config.DefaultJSONKeysPath
	flagCredsFrom = ""
	defer func() {
		flagConfig = oldConfig
		flagCredsFrom = oldCredsFrom
	}()

	payload := `{"ui":{"default_audit_package":""},"stores":{"huawei":{"service_account_file":"./config/huawei.json"}},"market_aliases":{},"target_group":"","target_section":"","hooks":{"feishu_webhook":""}}`
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("payload", payload); err != nil {
		t.Fatalf("WriteField: %v", err)
	}
	part, err := writer.CreateFormFile("store_file_huawei__service_account_file", "sa.json")
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := part.Write([]byte(`{"key_id":"demo"}`)); err != nil {
		t.Fatalf("part.Write: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/config/save", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	handleWebConfigSave(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	runtimeCfg, err := loadWebRuntimeConfig()
	if err != nil {
		t.Fatalf("loadWebRuntimeConfig: %v", err)
	}
	if got := runtimeCfg.Stores["huawei"]["service_account_file"]; got != "./config/huawei.json" {
		t.Fatalf("runtime huawei.service_account_file = %q, want ./config/huawei.json", got)
	}
	if _, err := os.Stat(filepath.Join(tmp, "config", "huawei.json")); err != nil {
		t.Fatalf("uploaded file not present for runtime use: %v", err)
	}
}

func TestListWebApps_MigratesExistingMainConfig(t *testing.T) {
	tmp := t.TempDir()
	oldConfig := flagConfig
	oldCredsFrom := flagCredsFrom
	flagConfig = filepath.Join(tmp, "config", "config.json")
	flagCredsFrom = ""
	defer func() {
		flagConfig = oldConfig
		flagCredsFrom = oldCredsFrom
	}()

	if err := os.MkdirAll(filepath.Dir(flagConfig), 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(flagConfig, []byte(`{"ui":{"default_audit_package":"com.demo.app"},"vivo":{"access_key":"k","access_secret":"s"}}`), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	items, current, err := listWebApps()
	if err != nil {
		t.Fatalf("listWebApps: %v", err)
	}
	if current == "" {
		t.Fatalf("current app should be initialized")
	}
	if len(items) != 1 {
		t.Fatalf("items len = %d, want 1", len(items))
	}
	if !items[0].Selected {
		t.Fatalf("expected migrated app to be selected")
	}
	if items[0].Name == "" {
		t.Fatalf("expected migrated app name")
	}
	if _, err := os.Stat(webAppConfigPath(current)); err != nil {
		t.Fatalf("migrated app config missing: %v", err)
	}
}

func TestSelectWebApp_OverridesMainConfig(t *testing.T) {
	tmp := t.TempDir()
	oldConfig := flagConfig
	oldCredsFrom := flagCredsFrom
	flagConfig = filepath.Join(tmp, "config", "config.json")
	flagCredsFrom = ""
	defer func() {
		flagConfig = oldConfig
		flagCredsFrom = oldCredsFrom
	}()

	if err := os.MkdirAll(filepath.Dir(flagConfig), 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	mainDoc := `{"ui":{"app_name":"App One","default_audit_package":"com.one.app"},"vivo":{"access_key":"one","access_secret":"one-secret"}}`
	if err := os.WriteFile(flagConfig, []byte(mainDoc), 0644); err != nil {
		t.Fatalf("WriteFile main: %v", err)
	}
	if _, _, err := listWebApps(); err != nil {
		t.Fatalf("listWebApps init: %v", err)
	}

	app2 := &webConfigDocument{
		UI: webUIConfig{
			AppName:             "App Two",
			DefaultAuditPackage: "com.two.app",
			ManualURLs:          map[string]string{},
		},
		Hooks:         map[string]map[string]string{},
		Stores:        map[string]map[string]string{"xiaomi": {"email": "two@example.com", "private_key": "key"}},
		MarketAliases: map[string][]string{},
	}
	if err := saveWebEditableConfigAt(webAppConfigPath("app-two"), app2); err != nil {
		t.Fatalf("saveWebEditableConfigAt: %v", err)
	}

	item, err := selectWebApp("app-two")
	if err != nil {
		t.Fatalf("selectWebApp: %v", err)
	}
	if item.Name != "App Two" {
		t.Fatalf("selected item name = %q, want App Two", item.Name)
	}

	doc, err := loadWebEditableConfig()
	if err != nil {
		t.Fatalf("loadWebEditableConfig: %v", err)
	}
	if doc.UI.AppName != "App Two" {
		t.Fatalf("main config app_name = %q, want App Two", doc.UI.AppName)
	}
	if doc.UI.DefaultAuditPackage != "com.two.app" {
		t.Fatalf("main config package = %q, want com.two.app", doc.UI.DefaultAuditPackage)
	}
	if got := doc.Stores["xiaomi"]["email"]; got != "two@example.com" {
		t.Fatalf("main config xiaomi.email = %q, want two@example.com", got)
	}
}

func TestHandleWebConfigSave_ForNonSelectedApp_DoesNotOverrideMainConfig(t *testing.T) {
	tmp := t.TempDir()
	oldConfig := flagConfig
	oldCredsFrom := flagCredsFrom
	flagConfig = filepath.Join(tmp, "config", "config.json")
	flagCredsFrom = ""
	defer func() {
		flagConfig = oldConfig
		flagCredsFrom = oldCredsFrom
	}()

	if err := os.MkdirAll(filepath.Dir(flagConfig), 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(flagConfig, []byte(`{"ui":{"app_name":"App One","default_audit_package":"com.one.app"},"vivo":{"access_key":"one","access_secret":"one-secret"}}`), 0644); err != nil {
		t.Fatalf("WriteFile main: %v", err)
	}
	if _, _, err := listWebApps(); err != nil {
		t.Fatalf("listWebApps init: %v", err)
	}

	app2 := &webConfigDocument{
		UI: webUIConfig{
			AppName:             "App Two",
			DefaultAuditPackage: "com.two.app",
			ManualURLs:          map[string]string{},
		},
		Hooks:         map[string]map[string]string{},
		Stores:        map[string]map[string]string{"xiaomi": {"email": "old@example.com", "private_key": "old-key"}},
		MarketAliases: map[string][]string{},
	}
	if err := saveWebEditableConfigAt(webAppConfigPath("app-two"), app2); err != nil {
		t.Fatalf("saveWebEditableConfigAt: %v", err)
	}

	payload := webConfigPayload{
		AppID: "app-two",
		UI: webUIConfig{
			AppName:             "App Two",
			DefaultAuditPackage: "com.two.updated",
		},
		Stores: map[string]map[string]string{
			"xiaomi": {
				"email":       "new@example.com",
				"private_key": "new-key",
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/config/save", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleWebConfigSave(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	mainDoc, err := loadWebEditableConfig()
	if err != nil {
		t.Fatalf("loadWebEditableConfig main: %v", err)
	}
	if mainDoc.UI.AppName != "App One" {
		t.Fatalf("main app_name = %q, want App One", mainDoc.UI.AppName)
	}
	if mainDoc.UI.DefaultAuditPackage != "com.one.app" {
		t.Fatalf("main package = %q, want com.one.app", mainDoc.UI.DefaultAuditPackage)
	}

	editedDoc, err := loadWebEditableConfigForApp("app-two")
	if err != nil {
		t.Fatalf("loadWebEditableConfigForApp: %v", err)
	}
	if editedDoc.UI.DefaultAuditPackage != "com.two.updated" {
		t.Fatalf("edited package = %q, want com.two.updated", editedDoc.UI.DefaultAuditPackage)
	}
	if got := editedDoc.Stores["xiaomi"]["email"]; got != "new@example.com" {
		t.Fatalf("edited xiaomi.email = %q, want new@example.com", got)
	}
}

func TestHandleWebHistory_UsesSelectedAppHistory(t *testing.T) {
	tmp := t.TempDir()
	oldConfig := flagConfig
	oldCredsFrom := flagCredsFrom
	flagConfig = filepath.Join(tmp, "config", "config.json")
	flagCredsFrom = ""
	defer func() {
		flagConfig = oldConfig
		flagCredsFrom = oldCredsFrom
	}()

	if err := os.MkdirAll(filepath.Dir(flagConfig), 0755); err != nil {
		t.Fatalf("MkdirAll config: %v", err)
	}
	if err := os.WriteFile(flagConfig, []byte(`{"ui":{"app_name":"App One","default_audit_package":"com.one.app"},"vivo":{"access_key":"one","access_secret":"one-secret"}}`), 0644); err != nil {
		t.Fatalf("WriteFile config: %v", err)
	}
	if _, _, err := listWebApps(); err != nil {
		t.Fatalf("listWebApps init: %v", err)
	}

	home := filepath.Join(tmp, "home")
	if err := os.MkdirAll(home, 0755); err != nil {
		t.Fatalf("MkdirAll home: %v", err)
	}
	oldHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", home); err != nil {
		t.Fatalf("Setenv HOME: %v", err)
	}
	defer func() {
		_ = os.Setenv("HOME", oldHome)
	}()

	if err := history.AppendRecord(webAppHistoryPath("app-one"), history.Record{
		Timestamp: "2026-06-24T10:00:00Z",
		APK: &apk.Info{
			PackageName: "com.one.app",
			VersionName: "1.0.0",
			VersionCode: 1,
			AppName:     "App One",
		},
	}); err != nil {
		t.Fatalf("AppendRecord one: %v", err)
	}
	app2 := &webConfigDocument{
		UI: webUIConfig{
			AppName:             "App Two",
			DefaultAuditPackage: "com.two.app",
			ManualURLs:          map[string]string{},
		},
		Hooks:         map[string]map[string]string{},
		Stores:        map[string]map[string]string{},
		MarketAliases: map[string][]string{},
	}
	if err := saveWebEditableConfigAt(webAppConfigPath("app-two"), app2); err != nil {
		t.Fatalf("saveWebEditableConfigAt: %v", err)
	}
	if err := history.AppendRecord(webAppHistoryPath("app-two"), history.Record{
		Timestamp: "2026-06-24T11:00:00Z",
		APK: &apk.Info{
			PackageName: "com.two.app",
			VersionName: "2.0.0",
			VersionCode: 2,
			AppName:     "App Two",
		},
	}); err != nil {
		t.Fatalf("AppendRecord two: %v", err)
	}
	if _, err := selectWebApp("app-one"); err != nil {
		t.Fatalf("selectWebApp app-one: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/history", nil)
	w := httptest.NewRecorder()
	handleWebHistory(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	var resp struct {
		Records []webHistoryItem `json:"records"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unmarshal response: %v", err)
	}
	if len(resp.Records) != 1 {
		t.Fatalf("records len = %d, want 1", len(resp.Records))
	}
	if resp.Records[0].PackageName != "com.one.app" {
		t.Fatalf("record package = %q, want com.one.app", resp.Records[0].PackageName)
	}

	if _, err := selectWebApp("app-two"); err != nil {
		t.Fatalf("selectWebApp app-two: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/history", nil)
	w = httptest.NewRecorder()
	handleWebHistory(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unmarshal response after switch: %v", err)
	}
	if len(resp.Records) != 1 {
		t.Fatalf("records len after switch = %d, want 1", len(resp.Records))
	}
	if resp.Records[0].PackageName != "com.two.app" {
		t.Fatalf("record package after switch = %q, want com.two.app", resp.Records[0].PackageName)
	}
}

func TestHandleWebHistoryDelete_SyncsSelectedAppHistory(t *testing.T) {
	tmp := t.TempDir()
	oldConfig := flagConfig
	oldCredsFrom := flagCredsFrom
	flagConfig = filepath.Join(tmp, "config", "config.json")
	flagCredsFrom = ""
	defer func() {
		flagConfig = oldConfig
		flagCredsFrom = oldCredsFrom
	}()

	if err := os.MkdirAll(filepath.Dir(flagConfig), 0755); err != nil {
		t.Fatalf("MkdirAll config: %v", err)
	}
	if err := os.WriteFile(flagConfig, []byte(`{"ui":{"app_name":"App One","default_audit_package":"com.one.app"}}`), 0644); err != nil {
		t.Fatalf("WriteFile config: %v", err)
	}
	if _, _, err := listWebApps(); err != nil {
		t.Fatalf("listWebApps init: %v", err)
	}

	if err := history.AppendRecord(webAppHistoryPath("app-one"), history.Record{
		Timestamp: "2026-06-24T10:00:00Z",
		APK: &apk.Info{
			PackageName: "com.one.app",
			VersionName: "1.0.0",
		},
	}); err != nil {
		t.Fatalf("AppendRecord app-one #1: %v", err)
	}
	if err := history.AppendRecord(webAppHistoryPath("app-one"), history.Record{
		Timestamp: "2026-06-24T11:00:00Z",
		APK: &apk.Info{
			PackageName: "com.one.app",
			VersionName: "1.0.1",
		},
	}); err != nil {
		t.Fatalf("AppendRecord app-one #2: %v", err)
	}
	if _, err := selectWebApp("app-one"); err != nil {
		t.Fatalf("selectWebApp app-one: %v", err)
	}

	body := bytes.NewReader([]byte(`{"timestamp":"2026-06-24T10:00:00Z"}`))
	req := httptest.NewRequest(http.MethodPost, "/api/history/delete", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleWebHistoryDelete(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	mainRecords, err := history.Read(mainWebHistoryPath())
	if err != nil {
		t.Fatalf("Read main history: %v", err)
	}
	if len(mainRecords) != 1 || mainRecords[0].Timestamp != "2026-06-24T11:00:00Z" {
		t.Fatalf("main records = %#v, want one remaining record", mainRecords)
	}

	appRecords, err := history.Read(webAppHistoryPath("app-one"))
	if err != nil {
		t.Fatalf("Read app history: %v", err)
	}
	if len(appRecords) != 1 || appRecords[0].Timestamp != "2026-06-24T11:00:00Z" {
		t.Fatalf("app records = %#v, want one remaining record", appRecords)
	}
}
