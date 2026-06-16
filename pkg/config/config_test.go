package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShouldPreferJSONKeys(t *testing.T) {
	cases := map[string]bool{
		"apkgo.yaml":         true,
		".":                  true,
		"":                   true,
		"config/config.json": true,
		"other.yaml":         false,
	}
	for input, want := range cases {
		if got := shouldPreferJSONKeys(input); got != want {
			t.Fatalf("shouldPreferJSONKeys(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestLoad_ConfigJSONPath(t *testing.T) {
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
	if err := os.MkdirAll(filepath.Dir(DefaultJSONKeysPath), 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	body := `{
		"hooks": {"after": "echo ok"},
		"market_aliases": {
			"xiaomi": ["xiaomi", "xm"]
		},
		"ui": {"default_audit_package": "com.example.app"},
		"pgyer": {"api_key": "k1"},
		"xiaomi": {"email": "demo@example.com", "private_key": "k2"}
	}`
	if err := os.WriteFile(DefaultJSONKeysPath, []byte(body), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := Load(DefaultJSONKeysPath)
	if err != nil {
		t.Fatalf("Load(%q): %v", DefaultJSONKeysPath, err)
	}
	if got := cfg.Stores["pgyer"]["api_key"]; got != "k1" {
		t.Fatalf("pgyer api_key = %q, want k1", got)
	}
	if got := cfg.Stores["xiaomi"]["private_key"]; got != "k2" {
		t.Fatalf("xiaomi private_key = %q, want k2", got)
	}
	if got := cfg.MarketAliases["xiaomi"]; len(got) != 2 || got[0] != "xiaomi" || got[1] != "xm" {
		t.Fatalf("xiaomi market_aliases = %#v, want [xiaomi xm]", got)
	}
	if _, ok := cfg.Stores["ui"]; ok {
		t.Fatalf("ui should not be parsed as a store")
	}
}

func TestLoad_FallbackToYAMLWhenJSONMissing(t *testing.T) {
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

	body := "stores:\n  pgyer:\n    api_key: \"k1\"\n"
	if err := os.WriteFile("apkgo.yaml", []byte(body), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := Load("apkgo.yaml")
	if err != nil {
		t.Fatalf("Load(apkgo.yaml): %v", err)
	}
	if got := cfg.Stores["pgyer"]["api_key"]; got != "k1" {
		t.Fatalf("pgyer api_key = %q, want k1", got)
	}
}
