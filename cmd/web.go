package cmd

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/KevinGong2013/apkgo/v3/pkg/apkgo"
	"github.com/KevinGong2013/apkgo/v3/pkg/config"
	"github.com/KevinGong2013/apkgo/v3/pkg/history"
	"github.com/KevinGong2013/apkgo/v3/pkg/store"
	"github.com/KevinGong2013/apkgo/v3/pkg/uploader"
)

const webUploadFormMemory = 2 << 30 // 2 GiB, large enough for local Jenkins web uploads

var (
	flagWebAddr string
)

type webArtifact struct {
	Store       string `json:"store"`
	DisplayName string `json:"display_name"`
	Channel     string `json:"channel"`
	FileName    string `json:"file_name"`
	Configured  bool   `json:"configured"`
}

type webArchiveBundle struct {
	Dir          string                    `json:"-"`
	Cleanup      func()                    `json:"-"`
	Artifacts    map[string]webArtifactRef `json:"-"`
	Summary      []webArtifact             `json:"artifacts"`
	AutoDetected bool                      `json:"auto_detected"`
}

type webArtifactRef struct {
	webArtifact
	Path string
}

type webUploadedFile struct {
	Name string
	Path string
}

type webStoreRunResult struct {
	Store       string        `json:"store"`
	DisplayName string        `json:"display_name"`
	Channel     string        `json:"channel"`
	FileName    string        `json:"file_name"`
	Result      *apkgo.Result `json:"result,omitempty"`
	Error       string        `json:"error,omitempty"`
}

type webConfigStore struct {
	Key         string `json:"key"`
	DisplayName string `json:"display_name"`
}

type webUIConfig struct {
	DefaultAuditPackage string            `json:"default_audit_package"`
	ManualURLs          map[string]string `json:"manual_urls"`
}

type webConfigField struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
	Accept      string `json:"accept,omitempty"`
	Secret      bool   `json:"secret,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Multiline   bool   `json:"multiline,omitempty"`
	Advanced    bool   `json:"advanced,omitempty"`
	File        bool   `json:"file,omitempty"`
}

type webConfigSection struct {
	Key         string           `json:"key"`
	DisplayName string           `json:"display_name"`
	Description string           `json:"description,omitempty"`
	DocURL      string           `json:"doc_url,omitempty"`
	Fields      []webConfigField `json:"fields"`
}

type webConfigListItem struct {
	GroupKey    string `json:"group_key"`
	Key         string `json:"key"`
	DisplayName string `json:"display_name"`
	Subtitle    string `json:"subtitle,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Configured  bool   `json:"configured"`
	SectionKey  string `json:"section_key"`
	EditLabel   string `json:"edit_label"`
}

type webConfigDocument struct {
	Hooks         map[string]map[string]string `json:"hooks,omitempty"`
	MarketAliases map[string][]string          `json:"market_aliases,omitempty"`
	UI            webUIConfig                  `json:"ui,omitempty"`
	UpdateCheck   string                       `json:"update_check,omitempty"`
	Stores        map[string]map[string]string `json:"stores,omitempty"`
}

type webConfigPayload struct {
	UI            webUIConfig                  `json:"ui"`
	Stores        map[string]map[string]string `json:"stores"`
	MarketAliases map[string][]string          `json:"market_aliases"`
	TargetGroup   string                       `json:"target_group,omitempty"`
	TargetSection string                       `json:"target_section,omitempty"`
	Hooks         struct {
		FeishuWebhook string `json:"feishu_webhook"`
	} `json:"hooks"`
}

type webConfigFileUpload struct {
	Store string
	Field string
	File  webUploadedFile
}

var webConfigSections = []webConfigSection{
	{
		Key:         "hooks",
		DisplayName: "通知",
		Description: "上传完成后的飞书机器人通知配置。",
		Fields: []webConfigField{
			{Key: "feishu_webhook", Label: "飞书 Webhook", Placeholder: "https://open.feishu.cn/open-apis/bot/v2/hook/...", Secret: true},
		},
	},
	{
		Key:         "ui",
		DisplayName: "Web 默认值",
		Description: "审核页使用的默认包名。",
		Fields: []webConfigField{
			{Key: "default_audit_package", Label: "默认包名", Placeholder: "com.example.app"},
		},
	},
}

type webStoreSectionMeta struct {
	Description string
	DocURL      string
}

var webStoreSectionMetaMap = map[string]webStoreSectionMeta{
	"huawei":     {DocURL: "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-Guides/agcapi-getstarted-0000001111845114#section1785535363715"},
	"xiaomi":     {DocURL: "https://dev.mi.com/xiaomihyperos/documentation/detail?pId=1134"},
	"oppo":       {DocURL: "https://open.oppomobile.com/new/developmentDoc/info?id=10998"},
	"vivo":       {DocURL: "https://dev.vivo.com.cn/documentCenter/doc/326"},
	"honor":      {DocURL: "https://developer.honor.com/cn/doc/guides/101360"},
	"tencent":    {DocURL: "https://wikinew.open.qq.com/index.html#/iwiki/4015262492"},
	"pgyer":      {DocURL: "https://www.pgyer.com/doc/view/app_upload"},
	"fir":        {DocURL: "https://www.betaqr.com.cn/docs"},
	"googleplay": {DocURL: "https://play.google.com/console"},
	"samsung":    {DocURL: "https://seller.samsungapps.com"},
}

type webConfigFieldMeta struct {
	Placeholder string
	Secret      bool
	Multiline   bool
	Advanced    bool
}

var webStoreFieldMetaMap = map[string][]webConfigField{
	"huawei": {
		{Key: "service_account_file", Label: "SERVICE_ACCOUNT_FILE", Placeholder: "./config/huawei.json", Accept: ".json,application/json", File: true},
	},
	"xiaomi": {
		{Key: "email", Label: "EMAIL"},
		{Key: "private_key", Label: "PRIVATE_KEY", Secret: true, Multiline: true},
		{Key: "cert_file", Label: "CERT_FILE", Placeholder: "./config/xiaomi.cer", Accept: ".cer,.crt,.pem", File: true},
	},
	"oppo": {
		{Key: "client_id", Label: "CLIENT_ID"},
		{Key: "client_secret", Label: "CLIENT_SECRET", Secret: true},
	},
	"vivo": {
		{Key: "access_key", Label: "ACCESS_KEY"},
		{Key: "access_secret", Label: "ACCESS_SECRET", Secret: true},
	},
	"honor": {
		{Key: "client_id", Label: "CLIENT_ID"},
		{Key: "client_secret", Label: "CLIENT_SECRET", Secret: true},
	},
	"tencent": {
		{Key: "user_id", Label: "USER_ID"},
		{Key: "app_id", Label: "APP_ID"},
		{Key: "access_secret", Label: "ACCESS_SECRET", Secret: true},
	},
	"pgyer": {
		{Key: "api_key", Label: "API_KEY", Secret: true},
	},
	"fir": {
		{Key: "api_token", Label: "API_TOKEN", Secret: true},
	},
	"googleplay": {
		{Key: "json_key_file", Label: "JSON_KEY_FILE"},
		{Key: "package_name", Label: "PACKAGE_NAME"},
		{Key: "track", Label: "TRACK"},
	},
	"samsung": {
		{Key: "service_account_id", Label: "SERVICE_ACCOUNT_ID"},
		{Key: "private_key", Label: "PRIVATE_KEY", Secret: true, Multiline: true},
		{Key: "content_id", Label: "CONTENT_ID"},
	},
	"script": {
		{Key: "command", Label: "COMMAND", Placeholder: "./deploy.sh"},
	},
}

var webStoreMeta = map[string]webConfigStore{
	"fir":        {Key: "fir", DisplayName: "fir.im"},
	"huawei":     {Key: "huawei", DisplayName: "华为"},
	"xiaomi":     {Key: "xiaomi", DisplayName: "小米"},
	"googleplay": {Key: "googleplay", DisplayName: "Google Play"},
	"oppo":       {Key: "oppo", DisplayName: "OPPO"},
	"samsung":    {Key: "samsung", DisplayName: "Samsung"},
	"script":     {Key: "script", DisplayName: "Script"},
	"vivo":       {Key: "vivo", DisplayName: "vivo"},
	"honor":      {Key: "honor", DisplayName: "荣耀"},
	"tencent":    {Key: "tencent", DisplayName: "应用宝"},
	"pgyer":      {Key: "pgyer", DisplayName: "蒲公英"},
}

func buildWebConfigSections() []webConfigSection {
	sections := make([]webConfigSection, 0, len(webConfigSections)+len(webStoreFieldMetaMap))
	sections = append(sections, webConfigSections...)

	var storeKeys []string
	for key := range webStoreFieldMetaMap {
		storeKeys = append(storeKeys, key)
	}
	sort.Strings(storeKeys)

	for _, storeKey := range storeKeys {
		meta := webStoreSectionMetaMap[storeKey]
		fields := append([]webConfigField(nil), webStoreFieldMetaMap[storeKey]...)
		sections = append(sections, webConfigSection{
			Key:         storeKey,
			DisplayName: storeDisplayName(storeKey),
			Description: meta.Description,
			DocURL:      meta.DocURL,
			Fields:      fields,
		})
	}

	return sections
}

func init() {
	webCmd.Flags().StringVar(&flagWebAddr, "addr", "127.0.0.1:8787", "web server listen address")
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the local web UI",
	RunE: func(cmd *cobra.Command, args []string) error {
		mux := http.NewServeMux()
		mux.HandleFunc("/", handleWebIndex)
		mux.HandleFunc("/audit", handleWebAuditPage)
		mux.HandleFunc("/config", handleWebConfigPage)
		mux.HandleFunc("/history", handleWebHistoryPage)
		mux.HandleFunc("/history/detail", handleWebHistoryDetailPage)
		mux.HandleFunc("/api/audit", handleWebAudit)
		mux.HandleFunc("/api/audit/sync-feishu", handleWebAuditSyncFeishu)
		mux.HandleFunc("/api/config", handleWebConfig)
		mux.HandleFunc("/api/config/save", handleWebConfigSave)
		mux.HandleFunc("/api/history", handleWebHistory)
		mux.HandleFunc("/api/history/delete", handleWebHistoryDelete)
		mux.HandleFunc("/api/stores", handleWebStores)
		mux.HandleFunc("/api/inspect", handleWebInspect)
		mux.HandleFunc("/api/upload", handleWebUpload)

		srv := &http.Server{
			Addr:              flagWebAddr,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		}

		slog.Info("web ui listening", "addr", flagWebAddr, "config", flagConfig)
		fmt.Fprintf(os.Stderr, "Open http://%s\n", flagWebAddr)
		return srv.ListenAndServe()
	},
}

func handleWebIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, webIndexHTML)
}

func handleWebAuditPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/audit" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, webAuditHTML)
}

func handleWebConfigPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/config" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, webConfigHTML)
}

func handleWebHistoryPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/history" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, webHistoryHTML)
}

func handleWebHistoryDetailPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/history/detail" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, webHistoryDetailHTML)
}

func handleWebStores(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeWebError(w, http.StatusMethodNotAllowed, "请求方法不支持")
		return
	}
	cfg, err := loadWebConfig()
	if err != nil {
		writeWebError(w, http.StatusInternalServerError, err.Error())
		return
	}
	stores := visibleWebStores(cfg)
	writeWebJSON(w, http.StatusOK, map[string]any{"stores": stores})
}

func handleWebConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeWebError(w, http.StatusMethodNotAllowed, "请求方法不支持")
		return
	}
	doc, err := loadWebEditableConfig()
	if err != nil {
		writeWebError(w, http.StatusInternalServerError, err.Error())
		return
	}
	cfg, err := loadWebConfig()
	if err != nil {
		cfg = &config.Config{Stores: map[string]map[string]string{}}
	}
	stores := visibleWebStores(cfg)
	sections := buildWebConfigSections()
	writeWebJSON(w, http.StatusOK, map[string]any{
		"path":              flagConfig,
		"configured_stores": stores,
		"ui":                doc.UI,
		"hooks": map[string]any{
			"feishu_webhook": extractFeishuWebhook(doc.Hooks["after"]["command"]),
		},
		"stores_config":  doc.Stores,
		"market_aliases": doc.MarketAliases,
		"items":          buildWebConfigListItems(doc),
		"sections":       sections,
	})
}

func handleWebConfigSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeWebError(w, http.StatusMethodNotAllowed, "请求方法不支持")
		return
	}
	if flagCredsFrom != "" {
		writeWebError(w, http.StatusBadRequest, "当前运行模式不支持在 Web 中保存配置")
		return
	}

	payload, uploads, cleanupUploads, err := parseWebConfigSaveRequest(r)
	if err != nil {
		writeWebError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer cleanupUploads()

	doc, err := loadWebEditableConfig()
	if err != nil {
		writeWebError(w, http.StatusInternalServerError, err.Error())
		return
	}

	doc.UI = payload.UI
	if doc.UI.ManualURLs == nil {
		doc.UI.ManualURLs = map[string]string{}
	}
	if payload.MarketAliases != nil {
		doc.MarketAliases = payload.MarketAliases
	}
	if doc.Stores == nil {
		doc.Stores = map[string]map[string]string{}
	}

	for _, section := range buildWebConfigSections() {
		if section.Key == "hooks" || section.Key == "ui" {
			continue
		}
		values := map[string]string{}
		for _, field := range section.Fields {
			if payload.Stores[section.Key] == nil {
				continue
			}
			v := strings.TrimSpace(payload.Stores[section.Key][field.Key])
			if v != "" {
				values[field.Key] = v
			}
		}
		if len(values) == 0 {
			delete(doc.Stores, section.Key)
			continue
		}
		doc.Stores[section.Key] = values
	}

	if err := applyWebConfigFileUploads(doc, uploads); err != nil {
		writeWebError(w, http.StatusBadRequest, err.Error())
		return
	}

	webhook := strings.TrimSpace(payload.Hooks.FeishuWebhook)
	if doc.Hooks == nil {
		doc.Hooks = map[string]map[string]string{}
	}
	if webhook == "" {
		delete(doc.Hooks, "after")
	} else {
		doc.Hooks["after"] = map[string]string{
			"command": fmt.Sprintf("go run . notify feishu --webhook '%s'", webhook),
		}
	}

	if strings.TrimSpace(payload.TargetGroup) == "stores" && strings.TrimSpace(payload.TargetSection) != "" {
		sectionKey := strings.TrimSpace(payload.TargetSection)
		if err := validateWebStoreSection(r.Context(), sectionKey, doc.Stores[sectionKey], doc.UI); err != nil {
			writeWebError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	if err := persistWebConfigUploads(doc, uploads); err != nil {
		writeWebError(w, http.StatusInternalServerError, fmt.Sprintf("save uploaded file: %v", err))
		return
	}

	if err := saveWebEditableConfig(doc); err != nil {
		writeWebError(w, http.StatusInternalServerError, fmt.Sprintf("save config: %v", err))
		return
	}

	cfg, err := loadWebConfig()
	if err != nil {
		cfg = &config.Config{Stores: map[string]map[string]string{}}
	}
	writeWebJSON(w, http.StatusOK, map[string]any{
		"ok":                true,
		"path":              flagConfig,
		"configured_stores": visibleWebStores(cfg),
	})
}

func handleWebHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeWebError(w, http.StatusMethodNotAllowed, "请求方法不支持")
		return
	}

	records, err := history.Read(history.DefaultPath())
	if err != nil {
		writeWebError(w, http.StatusInternalServerError, fmt.Sprintf("read history: %v", err))
		return
	}

	items := make([]webHistoryItem, 0, len(records))
	for i := len(records) - 1; i >= 0; i-- {
		items = append(items, newWebHistoryItem(records[i]))
	}

	writeWebJSON(w, http.StatusOK, map[string]any{
		"path":    history.DefaultPath(),
		"records": items,
	})
}

func handleWebHistoryDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeWebError(w, http.StatusMethodNotAllowed, "请求方法不支持")
		return
	}

	var payload struct {
		Timestamp string `json:"timestamp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeWebError(w, http.StatusBadRequest, fmt.Sprintf("解析删除参数失败: %v", err))
		return
	}
	if strings.TrimSpace(payload.Timestamp) == "" {
		writeWebError(w, http.StatusBadRequest, "缺少记录时间戳")
		return
	}

	err := history.DeleteByTimestamp(history.DefaultPath(), strings.TrimSpace(payload.Timestamp))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeWebError(w, http.StatusNotFound, "没有找到要删除的记录")
			return
		}
		writeWebError(w, http.StatusInternalServerError, fmt.Sprintf("删除记录失败: %v", err))
		return
	}

	writeWebJSON(w, http.StatusOK, map[string]any{
		"ok":        true,
		"timestamp": strings.TrimSpace(payload.Timestamp),
	})
}

func handleWebUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeWebError(w, http.StatusMethodNotAllowed, "请求方法不支持")
		return
	}
	if err := r.ParseMultipartForm(webUploadFormMemory); err != nil {
		writeWebError(w, http.StatusBadRequest, fmt.Sprintf("解析上传表单失败，可能是文件过大或上传中断: %v", err))
		return
	}

	cfg, err := loadWebRuntimeConfig()
	if err != nil {
		writeWebError(w, http.StatusBadRequest, fmt.Sprintf("读取 Web 配置失败：%v", err))
		return
	}

	files, cleanupFiles, err := saveUploadedFiles(r, "archive")
	if err != nil {
		writeWebError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer cleanupFiles()

	bundle, err := buildUploadBundle(cfg, files)
	if err != nil {
		writeWebError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer bundle.Cleanup()

	selectedStores := csvToSlice(r.FormValue("stores"))
	filtered := filterBundleForUpload(bundle, configuredStoreSet(cfg), selectedStores)
	if len(filtered.Artifacts) == 0 {
		writeWebError(w, http.StatusBadRequest, "没有可发布的市场，请先选择已配置 key 的市场")
		return
	}

	publishMode, publishTime, err := validatePublishFlags(r.FormValue("publish_mode"), r.FormValue("publish_time"))
	if err != nil {
		writeWebError(w, http.StatusBadRequest, err.Error())
		return
	}

	notes := r.FormValue("notes")
	dryRun := strings.EqualFold(r.FormValue("dry_run"), "true")
	streamWebHeaders(w)
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeWebStreamEvent(w, webStreamEvent{Type: "error", Message: "当前环境不支持流式返回"})
		return
	}
	result, err := runWebBundleStream(context.Background(), cfg, filtered, notes, publishMode, publishTime, dryRun, func(ev webStreamEvent) {
		writeWebStreamEvent(w, ev)
		flusher.Flush()
	})
	if err != nil {
		writeWebStreamEvent(w, webStreamEvent{Type: "error", Message: err.Error()})
		flusher.Flush()
		return
	}
	_ = result
}

func handleWebAudit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeWebError(w, http.StatusMethodNotAllowed, "请求方法不支持")
		return
	}
	if err := r.ParseMultipartForm(webUploadFormMemory); err != nil {
		writeWebError(w, http.StatusBadRequest, fmt.Sprintf("解析查询表单失败: %v", err))
		return
	}

	cfg, err := loadWebRuntimeConfig()
	if err != nil {
		writeWebError(w, http.StatusBadRequest, fmt.Sprintf("读取 Web 配置失败：%v", err))
		return
	}

	packageName := strings.TrimSpace(r.FormValue("package"))
	if packageName == "" {
		packageName = strings.TrimSpace(loadWebUIConfig().DefaultAuditPackage)
	}
	apkPath, cleanupAPK, err := saveOptionalUploadedFile(r, "apk")
	if err != nil {
		writeWebError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer cleanupAPK()

	if packageName == "" && apkPath == "" {
		writeWebError(w, http.StatusBadRequest, "请填写包名，或上传一个 APK 文件")
		return
	}

	report, err := apkgo.QueryAudit(r.Context(), apkgo.AuditJob{
		Config:  cfg,
		Package: packageName,
		APKFile: apkPath,
	})
	if err != nil {
		writeWebError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeWebJSON(w, http.StatusOK, map[string]any{
		"package": report.Package,
		"stores":  report.Stores,
	})
}

func handleWebAuditSyncFeishu(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeWebError(w, http.StatusMethodNotAllowed, "请求方法不支持")
		return
	}

	var report apkgo.AuditReport
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		writeWebError(w, http.StatusBadRequest, fmt.Sprintf("解析审核结果失败: %v", err))
		return
	}
	if strings.TrimSpace(report.Package) == "" || len(report.Stores) == 0 {
		writeWebError(w, http.StatusBadRequest, "请先查询到审核结果，再同步到飞书")
		return
	}

	webhook, err := loadWebFeishuWebhook()
	if err != nil {
		writeWebError(w, http.StatusBadRequest, err.Error())
		return
	}

	card := buildFeishuAuditCard(&report)
	if err := postFeishuCard(webhook, card); err != nil {
		writeWebError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeWebJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"package": report.Package,
		"stores":  report.Stores,
	})
}

func handleWebInspect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeWebError(w, http.StatusMethodNotAllowed, "请求方法不支持")
		return
	}
	if err := r.ParseMultipartForm(webUploadFormMemory); err != nil {
		writeWebError(w, http.StatusBadRequest, fmt.Sprintf("解析上传表单失败，可能是文件过大或上传中断: %v", err))
		return
	}

	cfg, err := loadWebRuntimeConfig()
	if err != nil {
		writeWebError(w, http.StatusBadRequest, fmt.Sprintf("读取 Web 配置失败：%v", err))
		return
	}

	files, cleanupFiles, err := saveUploadedFiles(r, "archive")
	if err != nil {
		writeWebError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer cleanupFiles()

	bundle, err := buildUploadBundle(cfg, files)
	if err != nil {
		writeWebError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer bundle.Cleanup()

	annotated := annotateBundleConfigured(bundle, configuredStoreSet(cfg))
	writeWebJSON(w, http.StatusOK, map[string]any{
		"upload": map[string]any{
			"artifacts":      annotated.Summary,
			"auto_detected":  annotated.AutoDetected,
			"selection_mode": webSelectionMode(annotated),
		},
	})
}

func loadWebConfig() (*config.Config, error) {
	return loadConfigForCmd()
}

func parseWebConfigSaveRequest(r *http.Request) (webConfigPayload, []webConfigFileUpload, func(), error) {
	var payload webConfigPayload
	if strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "multipart/form-data") {
		if err := r.ParseMultipartForm(webUploadFormMemory); err != nil {
			return payload, nil, func() {}, fmt.Errorf("解析配置参数失败：%v", err)
		}
		raw := strings.TrimSpace(r.FormValue("payload"))
		if raw == "" {
			return payload, nil, func() {}, fmt.Errorf("解析配置参数失败：缺少 payload")
		}
		if err := json.Unmarshal([]byte(raw), &payload); err != nil {
			return payload, nil, func() {}, fmt.Errorf("解析配置参数失败：%v", err)
		}
		uploads, cleanup, err := saveWebConfigUploadedFiles(r)
		if err != nil {
			return payload, nil, func() {}, err
		}
		return payload, uploads, cleanup, nil
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return payload, nil, func() {}, fmt.Errorf("解析配置参数失败：%v", err)
	}
	return payload, nil, func() {}, nil
}

func saveWebConfigUploadedFiles(r *http.Request) ([]webConfigFileUpload, func(), error) {
	if r.MultipartForm == nil {
		return nil, func() {}, nil
	}
	var uploads []webConfigFileUpload
	var paths []string
	cleanup := func() {
		for _, path := range paths {
			_ = os.Remove(path)
		}
	}

	for fieldName, headers := range r.MultipartForm.File {
		if !strings.HasPrefix(fieldName, "store_file_") {
			continue
		}
		storeKey, configField, ok := parseWebConfigUploadFieldName(fieldName)
		if !ok {
			continue
		}
		for _, header := range headers {
			f, err := header.Open()
			if err != nil {
				cleanup()
				return nil, func() {}, fmt.Errorf("读取 %s 失败: %w", header.Filename, err)
			}
			path, remove, err := copyMultipartToTemp(f, header)
			_ = f.Close()
			if err != nil {
				cleanup()
				return nil, func() {}, err
			}
			uploads = append(uploads, webConfigFileUpload{
				Store: storeKey,
				Field: configField,
				File: webUploadedFile{
					Name: header.Filename,
					Path: path,
				},
			})
			paths = append(paths, path)
			_ = remove
		}
	}

	return uploads, cleanup, nil
}

func parseWebConfigUploadFieldName(name string) (string, string, bool) {
	trimmed := strings.TrimPrefix(name, "store_file_")
	parts := strings.SplitN(trimmed, "__", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	storeKey := strings.TrimSpace(parts[0])
	fieldKey := strings.TrimSpace(parts[1])
	if storeKey == "" || fieldKey == "" {
		return "", "", false
	}
	return storeKey, fieldKey, true
}

func applyWebConfigFileUploads(doc *webConfigDocument, uploads []webConfigFileUpload) error {
	for _, upload := range uploads {
		if !supportsWebConfigFileUpload(upload.Store, upload.Field) {
			return fmt.Errorf("%s 暂不支持通过上传设置 %s", storeDisplayName(upload.Store), strings.ToUpper(upload.Field))
		}
		if doc.Stores[upload.Store] == nil {
			doc.Stores[upload.Store] = map[string]string{}
		}
		doc.Stores[upload.Store][upload.Field] = webConfigUploadedFileTargetPath(upload.Store, upload.File.Name)
	}
	return nil
}

func persistWebConfigUploads(doc *webConfigDocument, uploads []webConfigFileUpload) error {
	for _, upload := range uploads {
		dest := filepath.Clean(filepath.Join(filepath.Dir(flagConfig), filepath.Base(webConfigUploadedFileTargetPath(upload.Store, upload.File.Name))))
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		if err := copyLocalFile(upload.File.Path, dest); err != nil {
			return err
		}
		if doc.Stores[upload.Store] == nil {
			doc.Stores[upload.Store] = map[string]string{}
		}
		doc.Stores[upload.Store][upload.Field] = "./config/" + filepath.Base(dest)
	}
	return nil
}

func supportsWebConfigFileUpload(storeKey, fieldKey string) bool {
	for _, field := range webStoreFieldMetaMap[storeKey] {
		if field.Key == fieldKey && field.File {
			return true
		}
	}
	return false
}

func webConfigUploadedFileTargetPath(storeKey, originalName string) string {
	name := filepath.Base(strings.TrimSpace(originalName))
	ext := strings.ToLower(filepath.Ext(name))
	switch storeKey {
	case "huawei":
		if ext == "" {
			ext = ".json"
		}
		return "./config/huawei" + ext
	case "xiaomi":
		if ext == "" {
			ext = ".cer"
		}
		return "./config/xiaomi" + ext
	default:
		return "./config/" + name
	}
}

func loadWebEditableConfig() (*webConfigDocument, error) {
	doc := &webConfigDocument{
		Hooks:         map[string]map[string]string{},
		MarketAliases: map[string][]string{},
		UI: webUIConfig{
			ManualURLs: map[string]string{},
		},
		Stores: map[string]map[string]string{},
	}
	if flagCredsFrom != "" {
		return doc, nil
	}
	data, err := os.ReadFile(flagConfig)
	if err != nil {
		if os.IsNotExist(err) {
			return doc, nil
		}
		return nil, fmt.Errorf("read web config: %w", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse web config: %w", err)
	}

	if v, ok := raw["hooks"]; ok {
		var hooks map[string]string
		if err := json.Unmarshal(v, &hooks); err == nil {
			for k, value := range hooks {
				doc.Hooks[k] = map[string]string{"command": value}
			}
		}
	}
	if v, ok := raw["market_aliases"]; ok {
		_ = json.Unmarshal(v, &doc.MarketAliases)
	}
	if v, ok := raw["ui"]; ok {
		_ = json.Unmarshal(v, &doc.UI)
		if doc.UI.ManualURLs == nil {
			doc.UI.ManualURLs = map[string]string{}
		}
	}
	if v, ok := raw["update_check"]; ok {
		_ = json.Unmarshal(v, &doc.UpdateCheck)
	}
	for key, blob := range raw {
		if _, reserved := map[string]bool{
			"hooks": true, "market_aliases": true, "ui": true, "update_check": true, "stores": true,
		}[key]; reserved {
			continue
		}
		var values map[string]string
		if err := json.Unmarshal(blob, &values); err == nil {
			doc.Stores[key] = values
		}
	}
	return doc, nil
}

func saveWebEditableConfig(doc *webConfigDocument) error {
	out := map[string]any{}
	if len(doc.Hooks) > 0 {
		hooks := map[string]string{}
		for key, values := range doc.Hooks {
			if command := strings.TrimSpace(values["command"]); command != "" {
				hooks[key] = command
			}
		}
		if len(hooks) > 0 {
			out["hooks"] = hooks
		}
	}
	if len(doc.MarketAliases) > 0 {
		out["market_aliases"] = doc.MarketAliases
	}
	if strings.TrimSpace(doc.UI.DefaultAuditPackage) != "" || len(doc.UI.ManualURLs) > 0 {
		out["ui"] = doc.UI
	}
	if strings.TrimSpace(doc.UpdateCheck) != "" {
		out["update_check"] = doc.UpdateCheck
	}
	for storeKey, rawValues := range doc.Stores {
		values := map[string]string{}
		for key, value := range rawValues {
			if trimmed := strings.TrimSpace(value); trimmed != "" {
				values[key] = trimmed
			}
		}
		if len(values) > 0 {
			out[storeKey] = values
		}
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(flagConfig, data, 0644)
}

func loadWebUIConfig() webUIConfig {
	doc, err := loadWebEditableConfig()
	if err != nil {
		return webUIConfig{ManualURLs: map[string]string{}}
	}
	return doc.UI
}

func extractFeishuWebhook(command string) string {
	command = strings.TrimSpace(command)
	if command == "" {
		return ""
	}
	const marker = "--webhook"
	idx := strings.Index(command, marker)
	if idx < 0 {
		return ""
	}
	raw := strings.TrimSpace(command[idx+len(marker):])
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "'") {
		if end := strings.Index(raw[1:], "'"); end >= 0 {
			return raw[1 : end+1]
		}
	}
	if strings.HasPrefix(raw, "\"") {
		if end := strings.Index(raw[1:], "\""); end >= 0 {
			return raw[1 : end+1]
		}
	}
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func isPlaceholderValue(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return true
	}
	lower := strings.ToLower(trimmed)
	switch lower {
	case "com.example.app",
		"https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook":
		return true
	}
	return false
}

func hasConfiguredStoreValue(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return false
	}
	switch strings.ToLower(trimmed) {
	case "./config/huawei.json", "./config/xiaomi.cer":
		_, err := os.Stat(trimmed)
		return err == nil
	default:
		return true
	}
}

func loadWebFeishuWebhook() (string, error) {
	cfg, err := loadWebConfig()
	if err != nil {
		return "", err
	}
	command := strings.TrimSpace(cfg.Hooks.After)
	if command == "" {
		return "", fmt.Errorf("未找到飞书机器人配置")
	}
	const marker = "--webhook"
	idx := strings.Index(command, marker)
	if idx < 0 {
		return "", fmt.Errorf("未找到飞书 webhook 参数")
	}
	raw := strings.TrimSpace(command[idx+len(marker):])
	if raw == "" {
		return "", fmt.Errorf("飞书 webhook 为空")
	}
	if strings.HasPrefix(raw, "'") {
		if end := strings.Index(raw[1:], "'"); end >= 0 {
			return raw[1 : end+1], nil
		}
	}
	if strings.HasPrefix(raw, "\"") {
		if end := strings.Index(raw[1:], "\""); end >= 0 {
			return raw[1 : end+1], nil
		}
	}
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return "", fmt.Errorf("飞书 webhook 为空")
	}
	return fields[0], nil
}

func buildWebConfigListItems(doc *webConfigDocument) []webConfigListItem {
	items := make([]webConfigListItem, 0, len(webConfigSections)+len(doc.MarketAliases))
	storeItems := make([]webConfigListItem, 0, len(webStoreFieldMetaMap))
	for _, section := range buildWebConfigSections() {
		switch section.Key {
		case "hooks":
			webhook := extractFeishuWebhook(doc.Hooks["after"]["command"])
			configured := !isPlaceholderValue(webhook)
			items = append(items, webConfigListItem{
				GroupKey:    "hooks",
				Key:         "feishu",
				DisplayName: "飞书",
				Summary:     ternarySummary(configured, "已配置", "未配置"),
				Configured:  configured,
				SectionKey:  section.Key,
				EditLabel:   "编辑",
			})
		case "ui":
			configured := !isPlaceholderValue(doc.UI.DefaultAuditPackage)
			items = append(items, webConfigListItem{
				GroupKey:    "ui",
				Key:         "default_audit_package",
				DisplayName: "包名",
				Summary:     ternarySummary(configured, doc.UI.DefaultAuditPackage, "未配置"),
				Configured:  configured,
				SectionKey:  section.Key,
				EditLabel:   "编辑",
			})
		default:
			configured := hasConfiguredValues(doc.Stores[section.Key])
			storeItems = append(storeItems, webConfigListItem{
				GroupKey:    "stores",
				Key:         section.Key,
				DisplayName: storeDisplayName(section.Key),
				Subtitle:    section.Key,
				Summary:     ternarySummary(configured, "已配置", "未配置"),
				Configured:  configured,
				SectionKey:  section.Key,
				EditLabel:   "编辑",
			})
		}
	}
	slices.SortFunc(storeItems, func(a, b webConfigListItem) int {
		if a.Configured != b.Configured {
			if a.Configured {
				return -1
			}
			return 1
		}
		return strings.Compare(a.DisplayName, b.DisplayName)
	})
	items = append(items, storeItems...)

	aliases := doc.MarketAliases
	if len(aliases) == 0 {
		aliases = config.DefaultMarketAliases()
	}
	var aliasKeys []string
	for key := range aliases {
		aliasKeys = append(aliasKeys, key)
	}
	slices.Sort(aliasKeys)
	for _, key := range aliasKeys {
		parts := aliases[key]
		if len(parts) == 0 {
			parts = []string{key}
		}
		items = append(items, webConfigListItem{
			GroupKey:    "aliases",
			Key:         key,
			DisplayName: storeDisplayName(key),
			Subtitle:    key,
			Summary:     strings.Join(parts, "/"),
			Configured:  len(aliases[key]) > 0,
			SectionKey:  key,
			EditLabel:   "编辑",
		})
	}
	return items
}

func ternarySummary(ok bool, yes, no string) string {
	if ok {
		return yes
	}
	return no
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func validateWebStoreSection(ctx context.Context, sectionKey string, values map[string]string, ui webUIConfig) error {
	sectionKey = strings.TrimSpace(sectionKey)
	if sectionKey == "" {
		return fmt.Errorf("缺少要校验的市场")
	}
	clean := map[string]string{}
	for k, v := range values {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			clean[k] = trimmed
		}
	}
	if len(clean) == 0 {
		return fmt.Errorf("%s 配置为空", storeDisplayName(sectionKey))
	}

	cfg := &config.Config{
		Stores: map[string]map[string]string{
			sectionKey: clean,
		},
	}

	result, err := apkgo.Diagnose(ctx, apkgo.DiagnoseJob{
		Config:  cfg,
		Stores:  []string{sectionKey},
		Package: strings.TrimSpace(ui.DefaultAuditPackage),
	})
	if err == nil && result != nil && len(result.Stores) > 0 {
		report := result.Stores[0]
		if !report.Supported {
			_, createErr := store.Create(sectionKey, cloneStringMap(clean))
			if createErr != nil {
				return fmt.Errorf("%s 校验失败：%v", storeDisplayName(sectionKey), createErr)
			}
			return nil
		}
		var failed []string
		for _, probe := range report.Probes {
			if probe.Status == "fail" {
				msg := strings.TrimSpace(probe.Error)
				if msg == "" {
					msg = strings.TrimSpace(probe.Detail)
				}
				if msg == "" {
					msg = probe.Name
				}
				failed = append(failed, msg)
			}
		}
		if len(failed) > 0 {
			return fmt.Errorf("%s 校验失败：%s", storeDisplayName(sectionKey), strings.Join(failed, "；"))
		}
		return nil
	}

	_, createErr := store.Create(sectionKey, cloneStringMap(clean))
	if createErr != nil {
		return fmt.Errorf("%s 校验失败：%v", storeDisplayName(sectionKey), createErr)
	}
	return nil
}

func cloneStringMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func buildFeishuAuditCard(report *apkgo.AuditReport) map[string]any {
	successCount := 0
	reviewingCount := 0
	unsupportedCount := 0
	errorCount := 0
	elements := make([]map[string]any, 0, len(report.Stores)+3)
	stores := append([]apkgo.AuditStoreResult(nil), report.Stores...)
	slices.SortFunc(stores, func(a, b apkgo.AuditStoreResult) int {
		if d := webAuditSortRank(a) - webAuditSortRank(b); d != 0 {
			return d
		}
		return strings.Compare(storeDisplayName(a.Store), storeDisplayName(b.Store))
	})
	for _, item := range stores {
		statusText := "暂不支持"
		statusColor := "grey"
		versionLine := "版本号："
		if item.VersionCode != 0 {
			versionLine = fmt.Sprintf("版本号：%d", item.VersionCode)
		} else if item.VersionCodeRaw != 0 {
			versionLine = fmt.Sprintf("版本号：%d", item.VersionCodeRaw)
		}
		detailLines := []string{versionLine}
		if !item.Supported {
			unsupportedCount++
		} else if item.Error != "" {
			statusText = "查询失败"
			statusColor = "red"
			errorCount++
			detailLines = append(detailLines, item.Error)
		} else {
			switch item.State {
			case "approved":
				statusText = "审核通过"
				statusColor = "green"
				successCount++
			case "reviewing":
				statusText = "审核中"
				statusColor = "orange"
				reviewingCount++
			case "rejected":
				statusText = "审核驳回"
				statusColor = "red"
			case "withdrawn":
				statusText = "已撤回"
				statusColor = "grey"
			default:
				statusText = "状态未知"
				statusColor = "grey"
			}
		}
		if manualURL := manualViewURL(item.Store); manualURL != "" {
			detailLines = append(detailLines, fmt.Sprintf("[手动查看](%s)", manualURL))
		}
		detail := strings.Join(detailLines, "\n")
		elements = append(elements, map[string]any{
			"tag": "div",
			"text": map[string]any{
				"tag": "lark_md",
				"content": fmt.Sprintf("**%s**  `<font color='%s'>%s</font>`\n%s",
					escapeFeishu(storeDisplayName(item.Store)),
					statusColor,
					statusText,
					escapeFeishu(strings.TrimSpace(detail))),
			},
		})
	}

	summary := fmt.Sprintf("审核通过 %d 个，审核中 %d 个，暂不支持 %d 个，失败 %d 个", successCount, reviewingCount, unsupportedCount, errorCount)
	return map[string]any{
		"msg_type": "interactive",
		"card": map[string]any{
			"config": map[string]any{
				"wide_screen_mode": true,
				"enable_forward":   true,
			},
			"header": map[string]any{
				"title": map[string]any{
					"tag":     "plain_text",
					"content": "apkgo 审核状态同步",
				},
				"template": "blue",
			},
			"elements": append([]map[string]any{
				{
					"tag": "div",
					"text": map[string]any{
						"tag":     "lark_md",
						"content": fmt.Sprintf("**包名**\n%s", escapeFeishu(report.Package)),
					},
				},
				{
					"tag": "div",
					"text": map[string]any{
						"tag":     "lark_md",
						"content": escapeFeishu(summary),
					},
				},
				{"tag": "hr"},
			}, elements...),
		},
	}
}

func storeDisplayName(store string) string {
	if meta, ok := webStoreMeta[store]; ok && strings.TrimSpace(meta.DisplayName) != "" {
		return meta.DisplayName
	}
	return store
}

func webAuditSortRank(item apkgo.AuditStoreResult) int {
	if !item.Supported {
		return 90
	}
	if item.Error != "" {
		return 80
	}
	switch item.State {
	case "approved":
		return 10
	case "reviewing":
		return 20
	case "rejected":
		return 30
	case "withdrawn":
		return 40
	case "unknown":
		return 50
	default:
		return 60
	}
}

func manualViewURL(store string) string {
	ui := loadWebUIConfig()
	if url := strings.TrimSpace(ui.ManualURLs[store]); url != "" {
		return url
	}
	switch store {
	case "huawei":
		return "https://developer.huawei.com/consumer/cn/"
	case "tencent":
		return "https://open.tencent.com/"
	case "oppo":
		return "https://open.oppomobile.com/"
	case "honor":
		return "https://developer.honor.com/cn/"
	case "vivo":
		return "https://developer.vivo.com.cn/"
	case "xiaomi":
		return "https://dev.mi.com/xiaomihyperos"
	case "pgyer":
		return "https://www.pgyer.com/"
	default:
		return ""
	}
}

func loadWebRuntimeConfig() (*config.Config, error) {
	cfg, err := loadWebConfig()
	if err != nil {
		return nil, err
	}
	filtered := &config.Config{
		Hooks:         cfg.Hooks,
		MarketAliases: cfg.MarketAliases,
		UpdateCheck:   cfg.UpdateCheck,
		Stores:        map[string]map[string]string{},
	}
	for name, values := range cfg.Stores {
		if !hasConfiguredValues(values) {
			continue
		}
		clean := map[string]string{}
		for k, v := range values {
			if strings.TrimSpace(v) != "" {
				clean[k] = v
			}
		}
		filtered.Stores[name] = clean
	}
	if len(filtered.Stores) == 0 {
		return nil, fmt.Errorf("未配置任何可用市场，请先在 %s 中完成配置", flagConfig)
	}
	return filtered, nil
}

func configuredStoreSet(cfg *config.Config) map[string]bool {
	out := map[string]bool{}
	for name, values := range cfg.Stores {
		if hasConfiguredValues(values) {
			out[name] = true
		}
	}
	return out
}

func visibleWebStores(cfg *config.Config) []webConfigStore {
	var stores []webConfigStore
	for name, meta := range webStoreMeta {
		if hasConfiguredValues(cfg.Stores[name]) {
			stores = append(stores, meta)
		}
	}
	slices.SortFunc(stores, func(a, b webConfigStore) int {
		return strings.Compare(a.Key, b.Key)
	})
	return stores
}

func saveUploadedFiles(r *http.Request, field string) ([]webUploadedFile, func(), error) {
	if r.MultipartForm == nil {
		return nil, func() {}, fmt.Errorf("missing %s", field)
	}
	headers := r.MultipartForm.File[field]
	if len(headers) == 0 {
		return nil, func() {}, fmt.Errorf("missing %s", field)
	}

	files := make([]webUploadedFile, 0, len(headers))
	paths := make([]string, 0, len(headers))
	cleanup := func() {
		for _, path := range paths {
			_ = os.Remove(path)
		}
	}

	for _, header := range headers {
		f, err := header.Open()
		if err != nil {
			cleanup()
			return nil, func() {}, fmt.Errorf("读取 %s 失败: %w", header.Filename, err)
		}
		path, remove, err := copyMultipartToTemp(f, header)
		_ = f.Close()
		if err != nil {
			cleanup()
			return nil, func() {}, err
		}
		files = append(files, webUploadedFile{
			Name: header.Filename,
			Path: path,
		})
		paths = append(paths, path)
		_ = remove
	}
	return files, cleanup, nil
}

func saveOptionalUploadedFile(r *http.Request, field string) (string, func(), error) {
	f, header, err := r.FormFile(field)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", func() {}, nil
		}
		return "", func() {}, fmt.Errorf("读取 %s 失败: %w", field, err)
	}
	defer f.Close()
	return copyMultipartToTemp(f, header)
}

func copyMultipartToTemp(src multipart.File, header *multipart.FileHeader) (string, func(), error) {
	pattern := "apkgo-web-*-" + filepath.Base(header.Filename)
	tmp, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", func() {}, err
	}
	if _, err := io.Copy(tmp, src); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", func() {}, err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return "", func() {}, err
	}
	return tmp.Name(), func() { _ = os.Remove(tmp.Name()) }, nil
}

func buildUploadBundle(cfg *config.Config, files []webUploadedFile) (*webArchiveBundle, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("请先上传 zip 或 APK 文件")
	}

	zipFiles := make([]webUploadedFile, 0, len(files))
	apkFiles := make([]webUploadedFile, 0, len(files))
	for _, file := range files {
		switch strings.ToLower(filepath.Ext(file.Name)) {
		case ".zip":
			zipFiles = append(zipFiles, file)
		case ".apk":
			apkFiles = append(apkFiles, file)
		default:
			return nil, fmt.Errorf("暂不支持的文件类型：%s", filepath.Base(file.Name))
		}
	}

	if len(zipFiles) > 0 && len(apkFiles) > 0 {
		return nil, fmt.Errorf("不能同时上传 zip 和 APK，请二选一")
	}
	if len(zipFiles) > 1 {
		return nil, fmt.Errorf("暂时只支持上传 1 个 zip，或上传多个 APK")
	}
	if len(zipFiles) == 1 {
		return extractArchiveBundle(cfg, zipFiles[0].Path)
	}
	return buildAPKBundle(cfg, files)
}

func buildAPKBundle(cfg *config.Config, files []webUploadedFile) (*webArchiveBundle, error) {
	dir, err := os.MkdirTemp("", "apkgo-web-apk-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(dir) }

	bundle := &webArchiveBundle{
		Dir:          dir,
		Cleanup:      cleanup,
		Artifacts:    map[string]webArtifactRef{},
		AutoDetected: true,
	}

	for _, file := range files {
		base := filepath.Base(file.Name)
		meta, ok := detectChannelFromName(cfg, base)
		if !ok {
			continue
		}
		dest := uniqueBundlePath(dir, base)
		if err := copyLocalFile(file.Path, dest); err != nil {
			cleanup()
			return nil, fmt.Errorf("处理 %s 失败: %w", base, err)
		}
		ref := webArtifactRef{
			webArtifact: webArtifact{
				Store:       meta.Store,
				DisplayName: meta.DisplayName,
				Channel:     meta.Channel,
				FileName:    base,
			},
			Path: dest,
		}
		bundle.Artifacts[meta.Store] = ref
	}

	if len(bundle.Artifacts) == 0 {
		if len(files) == 1 {
			return buildManualSelectionBundle(dir, cleanup, files[0], cfg)
		}
		cleanup()
		return nil, fmt.Errorf("上传的 APK 里没有识别到支持的渠道包；多个文件未命中别名时，请修改文件名后重试，或改为单个 APK 手动选择市场")
	}

	for _, ref := range bundle.Artifacts {
		bundle.Summary = append(bundle.Summary, ref.webArtifact)
	}
	slices.SortFunc(bundle.Summary, func(a, b webArtifact) int {
		return strings.Compare(a.Store, b.Store)
	})
	return bundle, nil
}

func extractArchiveBundle(cfg *config.Config, zipPath string) (*webArchiveBundle, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("打开压缩包失败: %w", err)
	}
	defer reader.Close()

	dir, err := os.MkdirTemp("", "apkgo-web-archive-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(dir) }

	bundle := &webArchiveBundle{
		Dir:          dir,
		Cleanup:      cleanup,
		Artifacts:    map[string]webArtifactRef{},
		AutoDetected: true,
	}

	for _, f := range reader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		base := filepath.Base(f.Name)
		if !strings.HasSuffix(strings.ToLower(base), ".apk") {
			continue
		}
		meta, ok := detectChannelFromName(cfg, base)
		if !ok {
			continue
		}
		dest := uniqueBundlePath(dir, base)
		if err := extractZipFile(f, dest); err != nil {
			cleanup()
			return nil, fmt.Errorf("解压 %s 失败: %w", base, err)
		}
		ref := webArtifactRef{
			webArtifact: webArtifact{
				Store:       meta.Store,
				DisplayName: meta.DisplayName,
				Channel:     meta.Channel,
				FileName:    base,
			},
			Path: dest,
		}
		bundle.Artifacts[meta.Store] = ref
	}

	if len(bundle.Artifacts) == 0 {
		cleanup()
		return nil, fmt.Errorf("压缩包里没有识别到支持的渠道 APK；请按别名重命名后重试，或改为上传单个 APK 手动选择市场")
	}

	for _, ref := range bundle.Artifacts {
		bundle.Summary = append(bundle.Summary, ref.webArtifact)
	}
	slices.SortFunc(bundle.Summary, func(a, b webArtifact) int {
		return strings.Compare(a.Store, b.Store)
	})
	return bundle, nil
}

func uniqueBundlePath(dir, base string) string {
	base = filepath.Base(base)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	candidate := filepath.Join(dir, base)
	for i := 1; ; i++ {
		if _, err := os.Stat(candidate); errors.Is(err, os.ErrNotExist) {
			return candidate
		}
		candidate = filepath.Join(dir, fmt.Sprintf("%s-%d%s", name, i, ext))
	}
}

func extractZipFile(zf *zip.File, dest string) error {
	rc, err := zf.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return err
}

func copyLocalFile(srcPath, destPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func detectChannelFromName(cfg *config.Config, name string) (webArtifact, bool) {
	match, ok := cfg.MatchMarketByFilename(name)
	if !ok {
		return webArtifact{}, false
	}
	return webArtifact{
		Store:       match.Store,
		DisplayName: storeDisplayName(match.Store),
		Channel:     match.Alias,
	}, true
}

func annotateBundleConfigured(bundle *webArchiveBundle, configured map[string]bool) *webArchiveBundle {
	out := &webArchiveBundle{
		Dir:          bundle.Dir,
		Cleanup:      bundle.Cleanup,
		Artifacts:    map[string]webArtifactRef{},
		AutoDetected: bundle.AutoDetected,
	}
	for store, ref := range bundle.Artifacts {
		ref.Configured = configured[store]
		out.Artifacts[store] = ref
		out.Summary = append(out.Summary, ref.webArtifact)
	}
	slices.SortFunc(out.Summary, func(a, b webArtifact) int {
		return strings.Compare(a.Store, b.Store)
	})
	return out
}

func buildManualSelectionBundle(dir string, cleanup func(), file webUploadedFile, cfg *config.Config) (*webArchiveBundle, error) {
	base := filepath.Base(file.Name)
	dest := uniqueBundlePath(dir, base)
	if err := copyLocalFile(file.Path, dest); err != nil {
		cleanup()
		return nil, fmt.Errorf("处理 %s 失败: %w", base, err)
	}

	bundle := &webArchiveBundle{
		Dir:          dir,
		Cleanup:      cleanup,
		Artifacts:    map[string]webArtifactRef{},
		AutoDetected: false,
	}
	for _, meta := range visibleWebStores(cfg) {
		ref := webArtifactRef{
			webArtifact: webArtifact{
				Store:       meta.Key,
				DisplayName: meta.DisplayName,
				Channel:     "",
				FileName:    base,
			},
			Path: dest,
		}
		bundle.Artifacts[meta.Key] = ref
		bundle.Summary = append(bundle.Summary, ref.webArtifact)
	}
	slices.SortFunc(bundle.Summary, func(a, b webArtifact) int {
		return strings.Compare(a.Store, b.Store)
	})
	return bundle, nil
}

func webSelectionMode(bundle *webArchiveBundle) string {
	if bundle.AutoDetected {
		return "auto"
	}
	return "manual"
}

func filterBundleForUpload(bundle *webArchiveBundle, configured map[string]bool, selected []string) *webArchiveBundle {
	allowed := map[string]bool{}
	if len(selected) > 0 {
		for _, store := range selected {
			allowed[store] = true
		}
	}
	out := &webArchiveBundle{
		Dir:       bundle.Dir,
		Cleanup:   bundle.Cleanup,
		Artifacts: map[string]webArtifactRef{},
	}
	for store, ref := range bundle.Artifacts {
		if !configured[store] {
			continue
		}
		if len(allowed) > 0 && !allowed[store] {
			continue
		}
		ref.Configured = true
		out.Artifacts[store] = ref
		out.Summary = append(out.Summary, ref.webArtifact)
	}
	slices.SortFunc(out.Summary, func(a, b webArtifact) int {
		return strings.Compare(a.Store, b.Store)
	})
	return out
}

func csvToSlice(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

type webStreamEvent struct {
	Type    string `json:"type"`
	Store   string `json:"store,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type webHistoryItem struct {
	Timestamp    string                `json:"timestamp"`
	PublishedAt  string                `json:"published_at"`
	Status       string                `json:"status"`
	Notes        string                `json:"notes,omitempty"`
	PublishMode  string                `json:"publish_mode,omitempty"`
	PublishTime  string                `json:"publish_time,omitempty"`
	PackageName  string                `json:"package_name,omitempty"`
	VersionName  string                `json:"version_name,omitempty"`
	VersionCode  int32                 `json:"version_code,omitempty"`
	AppName      string                `json:"app_name,omitempty"`
	SuccessCount int                   `json:"success_count"`
	FailureCount int                   `json:"failure_count"`
	Stores       []string              `json:"stores,omitempty"`
	Results      []*store.UploadResult `json:"results,omitempty"`
}

func streamWebHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/x-ndjson; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
}

func writeWebStreamEvent(w http.ResponseWriter, ev webStreamEvent) {
	_ = json.NewEncoder(w).Encode(ev)
}

type webLineWriter struct {
	mu   sync.Mutex
	emit func(webStreamEvent)
	buf  string
}

func (w *webLineWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buf += string(p)
	for {
		idx := strings.IndexByte(w.buf, '\n')
		if idx < 0 {
			break
		}
		line := strings.TrimSpace(w.buf[:idx])
		w.buf = w.buf[idx+1:]
		if line != "" {
			w.emit(webStreamEvent{Type: "log", Message: line})
		}
	}
	return len(p), nil
}

func runWebBundleStream(ctx context.Context, cfg *config.Config, bundle *webArchiveBundle, notes, publishMode, publishTime string, dryRun bool, emit func(webStreamEvent)) (map[string]any, error) {
	var emitMu sync.Mutex
	safeEmit := func(ev webStreamEvent) {
		emitMu.Lock()
		defer emitMu.Unlock()
		emit(ev)
	}

	results := make([]webStoreRunResult, len(bundle.Summary))
	var wg sync.WaitGroup
	for i, artifact := range bundle.Summary {
		wg.Add(1)
		go func(index int, artifact webArtifact) {
			defer wg.Done()

			ref := bundle.Artifacts[artifact.Store]
			safeEmit(webStreamEvent{
				Type:    "store.start",
				Store:   artifact.Store,
				Message: fmt.Sprintf("开始发布 %s（%s）", artifact.DisplayName, artifact.FileName),
			})

			progressBuf := &webLineWriter{emit: safeEmit}
			loggerWriter := &webLineWriter{emit: safeEmit}
			logger := slog.New(slog.NewTextHandler(loggerWriter, &slog.HandlerOptions{Level: slog.LevelInfo}))
			nd := uploader.NewNDJSONManager(progressBuf)
			res, err := apkgo.Run(ctx, apkgo.Job{
				APKFile:     ref.Path,
				Stores:      []string{artifact.Store},
				Notes:       notes,
				PublishMode: publishMode,
				PublishTime: publishTime,
				Config:      cfg,
				Timeout:     flagTimeout,
				DryRun:      dryRun,
				Progress:    nd,
				Logger:      logger,
			})
			entry := webStoreRunResult{
				Store:       artifact.Store,
				DisplayName: artifact.DisplayName,
				Channel:     artifact.Channel,
				FileName:    artifact.FileName,
			}
			if err != nil {
				entry.Error = err.Error()
				safeEmit(webStreamEvent{
					Type:    "store.error",
					Store:   artifact.Store,
					Message: fmt.Sprintf("%s 发布失败：%s", artifact.DisplayName, err.Error()),
				})
			} else {
				entry.Result = res
				safeEmit(webStreamEvent{
					Type:    "store.done",
					Store:   artifact.Store,
					Message: fmt.Sprintf("%s 发布完成", artifact.DisplayName),
					Data:    res,
				})
			}
			if progressBuf.buf != "" {
				safeEmit(webStreamEvent{Type: "log", Message: strings.TrimSpace(progressBuf.buf)})
			}
			if loggerWriter.buf != "" {
				safeEmit(webStreamEvent{Type: "log", Message: strings.TrimSpace(loggerWriter.buf)})
			}
			results[index] = entry
		}(i, artifact)
	}
	wg.Wait()

	final := map[string]any{
		"archive": map[string]any{
			"artifacts": bundle.Summary,
		},
		"results": results,
	}
	appendWebHistory(results, notes, publishMode, publishTime, dryRun, emit)
	emit(webStreamEvent{Type: "result", Data: final})
	return final, nil
}

func appendWebHistory(results []webStoreRunResult, notes, publishMode, publishTime string, dryRun bool, emit func(webStreamEvent)) {
	if dryRun {
		return
	}

	for _, entry := range results {
		if entry.Result == nil || entry.Result.APK == nil || len(entry.Result.Results) == 0 {
			continue
		}
		err := history.AppendWithMeta(history.DefaultPath(), entry.Result.APK, entry.Result.Results, history.Meta{
			Notes:       strings.TrimSpace(notes),
			PublishMode: strings.TrimSpace(publishMode),
			PublishTime: strings.TrimSpace(publishTime),
		})
		if err != nil {
			emit(webStreamEvent{
				Type:    "log",
				Message: fmt.Sprintf("[history] 保存 %s 发布记录失败：%v", entry.DisplayName, err),
			})
		}
	}
}

func newWebHistoryItem(record history.Record) webHistoryItem {
	item := webHistoryItem{
		Timestamp:   record.Timestamp,
		PublishedAt: formatWebTime(record.Timestamp),
		Notes:       record.Notes,
		PublishMode: record.PublishMode,
		PublishTime: record.PublishTime,
		Results:     record.Results,
	}
	if record.APK != nil {
		item.PackageName = record.APK.PackageName
		item.VersionName = record.APK.VersionName
		item.VersionCode = record.APK.VersionCode
		item.AppName = record.APK.AppName
	}
	allFailed := len(record.Results) > 0
	for _, result := range record.Results {
		if result == nil {
			continue
		}
		item.Stores = append(item.Stores, result.Store)
		if result.Success {
			item.SuccessCount++
			allFailed = false
		} else {
			item.FailureCount++
		}
	}
	switch {
	case len(record.Results) == 0:
		item.Status = "unknown"
	case item.SuccessCount == len(record.Results):
		item.Status = "success"
	case allFailed:
		item.Status = "failed"
	default:
		item.Status = "partial"
	}
	if item.PublishedAt == "" {
		item.PublishedAt = record.Timestamp
	}
	return item
}

func formatWebTime(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}
	return t.In(time.Local).Format("2006-01-02 15:04:05")
}

func hasConfiguredValues(values map[string]string) bool {
	for _, v := range values {
		if hasConfiguredStoreValue(v) {
			return true
		}
	}
	return false
}

func writeWebJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeWebError(w http.ResponseWriter, status int, msg string) {
	writeWebJSON(w, status, map[string]string{"error": msg})
}
