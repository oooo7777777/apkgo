package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KevinGong2013/apkgo/v3/pkg/apk"
	"github.com/KevinGong2013/apkgo/v3/pkg/hooks"
	"github.com/KevinGong2013/apkgo/v3/pkg/store"
)

func init() {
	rootCmd.AddCommand(notifyCmd)
	notifyCmd.AddCommand(notifyFeishuCmd)
	notifyFeishuCmd.Flags().String("webhook", "", "Feishu bot webhook URL")
}

var notifyCmd = &cobra.Command{
	Use:    "notify",
	Short:  "Internal notification helpers",
	Hidden: true,
}

var notifyFeishuCmd = &cobra.Command{
	Use:    "feishu",
	Short:  "Send a Feishu card from hook stdin payload",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		webhook, _ := cmd.Flags().GetString("webhook")
		webhook = strings.TrimSpace(webhook)
		if webhook == "" {
			webhook = strings.TrimSpace(os.Getenv("APKGO_FEISHU_WEBHOOK"))
		}
		if webhook == "" {
			return fmt.Errorf("missing webhook: use --webhook or APKGO_FEISHU_WEBHOOK")
		}

		var payload hooks.AfterAllPayload
		if err := json.NewDecoder(os.Stdin).Decode(&payload); err != nil {
			return fmt.Errorf("decode hook payload: %w", err)
		}

		card, err := buildFeishuAfterAllCard(payload)
		if err != nil {
			return err
		}
		return postFeishuCard(webhook, card)
	},
}

func buildFeishuAfterAllCard(payload hooks.AfterAllPayload) (map[string]any, error) {
	info := payload.APK
	if info == nil {
		info = &apk.Info{}
	}

	success := 0
	failed := 0
	lines := make([]map[string]any, 0, len(payload.Results)+2)
	for _, r := range payload.Results {
		statusText := "成功"
		statusColor := "green"
		detail := fmt.Sprintf("耗时 %d ms", r.DurationMs)
		if !r.Success {
			statusText = "失败"
			statusColor = "red"
			detail = r.Error
			failed++
		} else {
			success++
		}
		lines = append(lines, map[string]any{
			"tag": "div",
			"text": map[string]any{
				"tag":     "lark_md",
				"content": fmt.Sprintf("**%s**  `<font color='%s'>%s</font>`\n%s", r.Store, statusColor, statusText, escapeFeishu(detail)),
			},
		})
	}

	summaryColor := "green"
	summaryText := "全部成功"
	if failed > 0 && success > 0 {
		summaryColor = "orange"
		summaryText = "部分成功"
	} else if failed > 0 {
		summaryColor = "red"
		summaryText = "全部失败"
	}

	versionName := fallbackText(info.VersionName, "未知版本")
	appName := fallbackText(info.AppName, "未识别应用名")
	notes := fallbackText(strings.TrimSpace(payload.Notes), "未填写")

	elements := []map[string]any{
		{
			"tag": "div",
			"text": map[string]any{
				"tag":     "lark_md",
				"content": fmt.Sprintf("**%s**\n<font color='%s'>%s</font>", escapeFeishu(appName), summaryColor, summaryText),
			},
		},
		{
			"tag": "div",
			"fields": []map[string]any{
				{
					"is_short": true,
					"text": map[string]any{
						"tag":     "lark_md",
						"content": fmt.Sprintf("**版本**\n%s", escapeFeishu(versionName)),
					},
				},
			},
		},
		{"tag": "hr"},
		{
			"tag": "div",
			"fields": []map[string]any{
				{
					"is_short": false,
					"text": map[string]any{
						"tag":     "lark_md",
						"content": fmt.Sprintf("**更新文案**\n%s", escapeFeishu(notes)),
					},
				},
			},
		},
		{
			"tag": "div",
			"text": map[string]any{
				"tag":     "lark_md",
				"content": fmt.Sprintf("**结果汇总**\n成功 %d 个，失败 %d 个", success, failed),
			},
		},
		{"tag": "hr"},
	}
	elements = append(elements, lines...)

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
					"content": "apkgo 发布通知",
				},
				"template": summaryColor,
			},
			"elements": elements,
		},
	}, nil
}

func postFeishuCard(webhook string, body map[string]any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal feishu card: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("post feishu webhook: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode feishu response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("feishu http %d: %s", resp.StatusCode, result.Msg)
	}
	if result.Code != 0 {
		return fmt.Errorf("feishu code %d: %s", result.Code, result.Msg)
	}
	return nil
}

func fallbackText(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func escapeFeishu(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

var _ = store.UploadResult{}
