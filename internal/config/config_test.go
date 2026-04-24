package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.json")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.PollInterval != 5 {
		t.Errorf("expected poll_interval 5, got %d", cfg.PollInterval)
	}
	if cfg.AlertFormat != "text" {
		t.Errorf("expected alert_format text, got %q", cfg.AlertFormat)
	}
	if len(cfg.AllowedPorts) != 0 {
		t.Errorf("expected empty allowed_ports, got %v", cfg.AllowedPorts)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/portwatch.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadConfigValid(t *testing.T) {
	p := writeTemp(t, `{
		"poll_interval": 10,
		"allowed_ports": [80, 443, 8080],
		"alert_format": "json",
		"log_file": "/var/log/portwatch.log"
	}`)

	cfg, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PollInterval != 10 {
		t.Errorf("expected 10, got %d", cfg.PollInterval)
	}
	if len(cfg.AllowedPorts) != 3 {
		t.Errorf("expected 3 allowed ports, got %d", len(cfg.AllowedPorts))
	}
	if cfg.AlertFormat != "json" {
		t.Errorf("expected json, got %q", cfg.AlertFormat)
	}
	if cfg.LogFile != "/var/log/portwatch.log" {
		t.Errorf("unexpected log_file %q", cfg.LogFile)
	}
}

func TestLoadConfigInvalidPollInterval(t *testing.T) {
	p := writeTemp(t, `{"poll_interval": 0}`)
	_, err := LoadConfig(p)
	if err == nil {
		t.Fatal("expected validation error for poll_interval 0")
	}
}

func TestLoadConfigInvalidAlertFormat(t *testing.T) {
	p := writeTemp(t, `{"alert_format": "xml"}`)
	_, err := LoadConfig(p)
	if err == nil {
		t.Fatal("expected validation error for unknown alert_format")
	}
}

func TestLoadConfigInvalidPort(t *testing.T) {
	p := writeTemp(t, `{"allowed_ports": [80, 99999]}`)
	_, err := LoadConfig(p)
	if err == nil {
		t.Fatal("expected validation error for out-of-range port")
	}
}

func TestLoadConfigUnknownField(t *testing.T) {
	p := writeTemp(t, `{"unknown_key": true}`)
	_, err := LoadConfig(p)
	if err == nil {
		t.Fatal("expected error for unknown JSON field")
	}
}
