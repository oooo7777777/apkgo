package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/KevinGong2013/apkgo/v3/pkg/apk"
	"github.com/KevinGong2013/apkgo/v3/pkg/store"
)

// Record is a single upload history entry.
type Record struct {
	Timestamp   string                `json:"timestamp"`
	APK         *apk.Info             `json:"apk"`
	Results     []*store.UploadResult `json:"results"`
	Notes       string                `json:"notes,omitempty"`
	PublishMode string                `json:"publish_mode,omitempty"`
	PublishTime string                `json:"publish_time,omitempty"`
}

// Meta captures optional publish context saved alongside a history entry.
type Meta struct {
	Notes       string
	PublishMode string
	PublishTime string
}

// DefaultPath returns ~/.apkgo/history.jsonl
func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".apkgo", "history.jsonl")
}

// Append adds a record to the history file (JSONL format, one JSON object per line).
func Append(path string, apkInfo *apk.Info, results []*store.UploadResult) error {
	return AppendWithMeta(path, apkInfo, results, Meta{})
}

// AppendWithMeta adds a record to the history file with publish metadata.
func AppendWithMeta(path string, apkInfo *apk.Info, results []*store.UploadResult, meta Meta) error {
	record := Record{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		APK:         apkInfo,
		Results:     results,
		Notes:       meta.Notes,
		PublishMode: meta.PublishMode,
		PublishTime: meta.PublishTime,
	}
	return AppendRecord(path, record)
}

// AppendRecord appends a fully-formed history record.
func AppendRecord(path string, record Record) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create history dir: %w", err)
	}

	if record.Timestamp == "" {
		record.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

// Read returns all records from the history file. Returns empty slice if file doesn't exist.
func Read(path string) ([]Record, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var records []Record
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var r Record
		if err := json.Unmarshal(line, &r); err != nil {
			continue // skip malformed lines
		}
		records = append(records, r)
	}
	return records, nil
}

// DeleteByTimestamp removes a single record matched by timestamp.
func DeleteByTimestamp(path, timestamp string) error {
	records, err := Read(path)
	if err != nil {
		return err
	}

	filtered := make([]Record, 0, len(records))
	removed := false
	for _, record := range records {
		if !removed && record.Timestamp == timestamp {
			removed = true
			continue
		}
		filtered = append(filtered, record)
	}
	if !removed {
		return os.ErrNotExist
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create history dir: %w", err)
	}
	if len(filtered) == 0 {
		return os.WriteFile(path, nil, 0644)
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), "history-*.jsonl")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	for _, record := range filtered {
		data, err := json.Marshal(record)
		if err != nil {
			tmp.Close()
			return err
		}
		if _, err := tmp.Write(append(data, '\n')); err != nil {
			tmp.Close()
			return err
		}
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
