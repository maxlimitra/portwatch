package portscanner

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestHexToIPv6Valid(t *testing.T) {
	// ::1 in kernel little-endian word format
	// 00000000 00000000 00000000 01000000 → reversed per word → 00000000 00000000 00000000 00000001
	ip, err := hexToIPv6("00000000000000000000000001000000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := net.ParseIP("::1").To16()
	if !ip.Equal(expected) {
		t.Errorf("got %v, want %v", ip, expected)
	}
}

func TestHexToIPv6InvalidLength(t *testing.T) {
	_, err := hexToIPv6("0000")
	if err == nil {
		t.Fatal("expected error for short input")
	}
}

func TestHexToIPv6InvalidChars(t *testing.T) {
	_, err := hexToIPv6("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ")
	if err == nil {
		t.Fatal("expected error for non-hex input")
	}
}

func TestParseProcNet6MissingFile(t *testing.T) {
	_, err := parseProcNet6("/nonexistent/path/tcp6", "tcp6", nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseProcNet6ValidFile(t *testing.T) {
	// Minimal /proc/net/tcp6 fixture: loopback ::1 on port 8080 (0x1F90)
	// Address column: 00000000000000000000000001000000:1F90
	content := `  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000000000000000000001000000:1F90 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0
`
	tmp := filepath.Join(t.TempDir(), "tcp6")
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	entries, err := parseProcNet6(tmp, "tcp6", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Port != 8080 {
		t.Errorf("port: got %d, want 8080", e.Port)
	}
	if e.Protocol != "tcp6" {
		t.Errorf("protocol: got %q, want \"tcp6\"", e.Protocol)
	}
	expectedIP := net.ParseIP("::1").To16()
	if !e.IP.Equal(expectedIP) {
		t.Errorf("IP: got %v, want %v", e.IP, expectedIP)
	}
}
