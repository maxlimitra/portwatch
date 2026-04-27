package portscanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseProcNetICMPMissingFile(t *testing.T) {
	_, err := ParseProcNetICMP("/nonexistent/proc/net/icmp")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseProcNetICMP6MissingFile(t *testing.T) {
	_, err := ParseProcNetICMP6("/nonexistent/proc/net/icmp6")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseProcNetICMPValidFile(t *testing.T) {
	// Re-use the same hex format as /proc/net/tcp
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0035 00000000:0000 0A 00000000:00000000 00:00000000 00000000   101        0 12345 1 0000000000000000 100 0 0 10 0
`
	dir := t.TempDir()
	path := filepath.Join(dir, "icmp")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	entries, err := ParseProcNetICMP(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Protocol != "icmp" {
		t.Errorf("expected protocol icmp, got %s", entries[0].Protocol)
	}
	if entries[0].LocalPort != 53 {
		t.Errorf("expected local port 53, got %d", entries[0].LocalPort)
	}
}

func TestICMPEntryString(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "icmp")
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000   0        0 99 1 0000000000000000 100 0 0 10 0
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	entries, err := ParseProcNetICMP(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}
	s := entries[0].String()
	if s == "" {
		t.Error("String() returned empty string")
	}
	// Should contain protocol prefix
	if len(s) < 4 || s[:4] != "icmp" {
		t.Errorf("String() does not start with 'icmp': %s", s)
	}
}
