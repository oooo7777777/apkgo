package config

import "testing"

func TestEffectiveMarketAliases_OverridesDefaults(t *testing.T) {
	cfg := &Config{
		MarketAliases: map[string][]string{
			"xiaomi": {"xiaomi", "mi"},
			"custom": {"alpha"},
		},
	}

	aliases := cfg.EffectiveMarketAliases()
	if got := aliases["xiaomi"]; len(got) != 2 || got[0] != "xiaomi" || got[1] != "mi" {
		t.Fatalf("xiaomi aliases = %#v, want [xiaomi mi]", got)
	}
	if got := aliases["tencent"]; len(got) != 2 || got[0] != "tencent" || got[1] != "qq" {
		t.Fatalf("tencent aliases = %#v, want default aliases", got)
	}
	if got := aliases["custom"]; len(got) != 1 || got[0] != "alpha" {
		t.Fatalf("custom aliases = %#v, want [alpha]", got)
	}
}

func TestMatchMarketByFilename(t *testing.T) {
	cfg := &Config{}
	match, ok := cfg.MatchMarketByFilename("myapp_xm_release.apk")
	if !ok {
		t.Fatalf("expected xiaomi alias to match")
	}
	if match.Store != "xiaomi" || match.Alias != "xm" {
		t.Fatalf("match = %#v, want xiaomi/xm", match)
	}

	match, ok = cfg.MatchMarketByFilename("myapp.qq.release.apk")
	if !ok {
		t.Fatalf("expected tencent alias to match")
	}
	if match.Store != "tencent" || match.Alias != "qq" {
		t.Fatalf("match = %#v, want tencent/qq", match)
	}
}

func TestMatchMarketByFilename_PrefersLongestAlias(t *testing.T) {
	cfg := &Config{
		MarketAliases: map[string][]string{
			"xiaomi": {"mi"},
			"custom": {"xiaomi"},
		},
	}

	match, ok := cfg.MatchMarketByFilename("demo_xiaomi_prod.apk")
	if !ok {
		t.Fatalf("expected a match")
	}
	if match.Store != "custom" || match.Alias != "xiaomi" {
		t.Fatalf("match = %#v, want custom/xiaomi", match)
	}
}

func TestMatchMarketByFilename_RegexEscapesAlias(t *testing.T) {
	cfg := &Config{
		MarketAliases: map[string][]string{
			"custom": {"mi+"},
		},
	}

	match, ok := cfg.MatchMarketByFilename("demo-mi+-prod.apk")
	if !ok {
		t.Fatalf("expected literal alias with regex chars to match")
	}
	if match.Store != "custom" || match.Alias != "mi+" {
		t.Fatalf("match = %#v, want custom/mi+", match)
	}
}
