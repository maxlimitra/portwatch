package portscanner

import (
	"os"
	"testing"
)

func TestHexToIP(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0100007F", "127.0.0.1"},
		{"00000000", "0.0.0.0"},
		{"0F02000A", "10.0.2.15"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := hexToIP(tt.input)
			if got != tt.expected {
				t.Errorf("hexToIP(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseProcNetMissingFile(t *testing.T) {
	_, err := parseProcNet("/nonexistent/path")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseProcNetValidFile(t *testing.T) {
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0
   1: 00000000:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 67890 1 0000000000000000 100 0 0 10 0
`
	tmpFile, err := os.CreateTemp("", "proc_net_tcp")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	entries, err := parseProcNet(tmpFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", entries[0].Port)
	}
	if entries[1].Port != 80 {
		t.Errorf("expected port 80, got %d", entries[1].Port)
	}
}

func TestNewScanner(t *testing.T) {
	s := NewScanner()
	if s == nil {
		t.Error("NewScanner returned nil")
	}
}
