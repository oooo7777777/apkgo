package cmd

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KevinGong2013/apkgo/v3/pkg/config"
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

func TestBuildAPKBundle_ManualSelectionFallbackForSingleFile(t *testing.T) {
	tmp := t.TempDir()
	apkPath := filepath.Join(tmp, "demo-release.apk")
	if err := os.WriteFile(apkPath, []byte("apk"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg := &config.Config{
		Stores: map[string]map[string]string{
			"huawei": {"client_id": "x"},
			"xiaomi": {"email": "x"},
		},
	}
	bundle, err := buildAPKBundle(cfg, []webUploadedFile{{Name: "demo-release.apk", Path: apkPath}})
	if err != nil {
		t.Fatalf("buildAPKBundle: %v", err)
	}
	defer bundle.Cleanup()

	if bundle.AutoDetected {
		t.Fatalf("expected manual selection fallback")
	}
	if len(bundle.Summary) != 2 {
		t.Fatalf("summary len = %d, want 2", len(bundle.Summary))
	}
	for _, item := range bundle.Summary {
		if item.Channel != "" {
			t.Fatalf("channel = %q, want empty for manual selection", item.Channel)
		}
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
	var out map[string]map[string]string
	if err := json.Unmarshal(saved, &out); err != nil {
		t.Fatalf("Unmarshal saved: %v", err)
	}
	if got := out["vivo"]["access_secret"]; got != "" {
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
