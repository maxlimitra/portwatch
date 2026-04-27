package portscanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseProcNetUDPMissingFile(t *testing.T) {
	_, err := ParseProcNetUDP("/nonexistent/proc/net/udp")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseProcNetUDP6MissingFile(t *testing.T) {
	_, err := ParseProcNetUDP6("/nonexistent/proc/net/udp6")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseProcNetUDPValidFile(t *testing.T) {
	// Same format as /proc/net/tcp — reuse a minimal fixture.
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000   101        0 12345 2 0000000000000000
   1: 00000000:0277 00000000:0000 07 00000000:00000000 00:00000000 00000000     0        0 23456 2 0000000000000000
`
	dir := t.TempDir()
	path := filepath.Join(dir, "udp")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	entries, err := ParseProcNetUDP(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	for _, e := range entries {
		if e.Protocol != "udp" {
			t.Errorf("expected protocol 'udp', got %q", e.Protocol)
		}
		if e.LocalPort == 0 {
			t.Error("expected non-zero local port")
		}
		if e.LocalIP == nil {
			t.Error("expected non-nil local IP")
		}
	}
}

func TestUDPEntryString(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "udp")
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000   101        0 12345 2 0000000000000000
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	entries, err := ParseProcNetUDP(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}

	s := entries[0].String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
	if entries[0].Protocol != "udp" {
		t.Errorf("expected 'udp' in string, got %q", s)
	}
}
