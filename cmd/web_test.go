package cmd

import (
	"os"
	"path/filepath"
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
