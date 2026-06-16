package config

import (
	"regexp"
	"sort"
	"strings"
)

// MarketMatch is a resolved file-name match against a configured market alias.
type MarketMatch struct {
	Store string
	Alias string
}

// DefaultMarketAliases returns the built-in file-name aliases used to map APK
// packages to target stores. Config values may override these per store.
func DefaultMarketAliases() map[string][]string {
	return map[string][]string{
		"fir":        {"fir"},
		"googleplay": {"googleplay"},
		"honor":      {"honor"},
		"huawei":     {"huawei"},
		"oppo":       {"oppo"},
		"pgyer":      {"pgyer", "merit"},
		"samsung":    {"samsung"},
		"script":     {"script"},
		"tencent":    {"tencent", "qq"},
		"vivo":       {"vivo"},
		"xiaomi":     {"xiaomi", "xm"},
	}
}

// EffectiveMarketAliases returns the normalized market alias table with config
// overrides applied on top of the built-in defaults.
func (c *Config) EffectiveMarketAliases() map[string][]string {
	out := cloneMarketAliases(DefaultMarketAliases())
	for store, aliases := range normalizeMarketAliases(c.MarketAliases) {
		out[store] = aliases
	}
	return out
}

// MatchMarketByFilename detects a target store from a file name using the
// configured alias table. A match succeeds when the filename contains the
// alias anywhere, using a case-insensitive regular expression.
func (c *Config) MatchMarketByFilename(name string) (MarketMatch, bool) {
	type candidate struct {
		store string
		alias string
	}

	lower := strings.ToLower(name)
	var matches []candidate
	for store, aliases := range c.EffectiveMarketAliases() {
		for _, alias := range aliases {
			pattern := regexp.MustCompile("(?i)" + regexp.QuoteMeta(alias))
			if pattern.MatchString(lower) {
				matches = append(matches, candidate{store: store, alias: alias})
			}
		}
	}
	if len(matches) == 0 {
		return MarketMatch{}, false
	}

	sort.Slice(matches, func(i, j int) bool {
		if len(matches[i].alias) != len(matches[j].alias) {
			return len(matches[i].alias) > len(matches[j].alias)
		}
		if matches[i].alias != matches[j].alias {
			return matches[i].alias < matches[j].alias
		}
		return matches[i].store < matches[j].store
	})
	return MarketMatch{Store: matches[0].store, Alias: matches[0].alias}, true
}

func normalizeMarketAliases(raw map[string][]string) map[string][]string {
	if len(raw) == 0 {
		return map[string][]string{}
	}
	out := make(map[string][]string, len(raw))
	for store, aliases := range raw {
		store = strings.ToLower(strings.TrimSpace(store))
		if store == "" {
			continue
		}
		seen := map[string]bool{}
		var clean []string
		for _, alias := range aliases {
			alias = strings.ToLower(strings.TrimSpace(alias))
			if alias == "" || seen[alias] {
				continue
			}
			seen[alias] = true
			clean = append(clean, alias)
		}
		if len(clean) == 0 {
			continue
		}
		out[store] = clean
	}
	return out
}

func cloneMarketAliases(src map[string][]string) map[string][]string {
	out := make(map[string][]string, len(src))
	for store, aliases := range src {
		out[store] = append([]string(nil), aliases...)
	}
	return out
}
