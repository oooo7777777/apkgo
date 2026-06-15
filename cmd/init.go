package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KevinGong2013/apkgo/v3/pkg/store"
)

var flagInitStore string

func init() {
	initCmd.Flags().StringVarP(&flagInitStore, "store", "s", "", "comma-separated store names (default: all)")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a config template",
	Example: `  apkgo init
  apkgo init --store huawei,xiaomi
  apkgo init -c config/config.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if file already exists
		if _, err := os.Stat(flagConfig); err == nil {
			return fmt.Errorf("%s already exists (use -c to specify a different path)", flagConfig)
		}

		// Determine which stores to include
		wanted := map[string]bool{}
		if flagInitStore != "" {
			for _, s := range strings.Split(flagInitStore, ",") {
				wanted[strings.TrimSpace(s)] = true
			}
		}

		schemas := store.Schemas()

		content, included, err := buildInitConfig(flagConfig, schemas, wanted)
		if err != nil {
			return err
		}

		if included == 0 {
			return fmt.Errorf("no matching stores found; available: %s", strings.Join(store.Names(), ", "))
		}

		if err := os.WriteFile(flagConfig, content, 0644); err != nil {
			return fmt.Errorf("write config: %w", err)
		}

		slog.Info("config created", "path", flagConfig)
		writeOutput(map[string]string{
			"created": flagConfig,
			"stores":  fmt.Sprintf("%d", included),
		})
		return nil
	},
}

func buildInitConfig(path string, schemas []store.ConfigSchema, wanted map[string]bool) ([]byte, int, error) {
	if filepath.Ext(path) == ".json" {
		cfg := map[string]any{
			"hooks": map[string]string{
				"after": "",
			},
			"ui": map[string]any{
				"default_audit_package": "com.example.app",
				"manual_urls": map[string]string{
					"huawei":  "https://developer.huawei.com/consumer/cn/",
					"tencent": "https://open.tencent.com/",
					"oppo":    "https://open.oppomobile.com/",
					"honor":   "https://developer.honor.com/cn/",
					"vivo":    "https://developer.vivo.com.cn/",
					"xiaomi":  "https://dev.mi.com/xiaomihyperos",
					"pgyer":   "https://www.pgyer.com/",
				},
			},
		}

		included := 0
		for _, schema := range schemas {
			if len(wanted) > 0 && !wanted[schema.Name] {
				continue
			}
			included++
			values := map[string]string{}
			for _, f := range schema.Fields {
				values[f.Key] = ""
			}
			if schema.Name == "huawei" {
				values["service_account_file"] = "./config/huawei.json"
			}
			if schema.Name == "xiaomi" {
				values["cert_file"] = "./config/xiaomi.cer"
			}
			cfg[schema.Name] = values
		}

		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return nil, 0, fmt.Errorf("marshal config: %w", err)
		}
		data = append(data, '\n')
		return data, included, nil
	}

	var b strings.Builder
	b.WriteString("# apkgo configuration\n")
	b.WriteString("# Docs: https://github.com/KevinGong2013/apkgo\n\n")
	b.WriteString("stores:\n")

	included := 0
	for _, schema := range schemas {
		if len(wanted) > 0 && !wanted[schema.Name] {
			continue
		}
		included++
		b.WriteString(fmt.Sprintf("  %s:\n", schema.Name))
		for _, f := range schema.Fields {
			req := ""
			if f.Required {
				req = " (required)"
			}
			b.WriteString(fmt.Sprintf("    # %s%s\n", f.Desc, req))
			b.WriteString(fmt.Sprintf("    %s: \"\"\n", f.Key))
		}
		b.WriteString("\n")
	}

	b.WriteString("# Update check interval: 30d (default), 7d, 0 to disable\n")
	b.WriteString("# update_check: 30d\n")
	return []byte(b.String()), included, nil
}
