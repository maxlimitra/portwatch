package portscanner

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestHexToIPv6Valid(t *testing.T) {
	// ::1 in /proc/net/tcp6 little-endian word encoding
	// Each 32-bit word of ::1 is 0x00000000 0x00000000 0x00000000 0x01000000
	input := "000000000000000000000000" + "01000000"
	ip, err := hexToIPv6(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := net.ParseIP("::1")
	if !ip.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, ip)
	}
}

func TestHexToIPv6InvalidLength(t *testing.T) {
	_, err := hexToIPv6("deadbeef")
	if err == nil {
		t.Fatal("expected error for short input")
	}
}

func TestHexToIPv6InvalidChars(t *testing.T) {
	_, err := hexToIPv6("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ")
	if err == nil {
		t.Fatal("expected error for invalid hex")
	}
}

func TestParseProcNet6MissingFile(t *testing.T) {
	_, err := parseProcNet6("/nonexistent/proc/net/tcp6", "tcp6")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseProcNet6ValidFile(t *testing.T) {
	// Minimal /proc/net/tcp6 content with loopback ::1 on port 8080 (0x1F90)
	content := `  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000000000000000000001000000:1F90 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0
`
	dir := t.TempDir()
	path := filepath.Join(dir, "tcp6")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	entries, err := parseProcNet6(path, "tcp6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", entries[0].Port)
	}
	if entries[0].Protocol != "tcp6" {
		t.Errorf("expected protocol tcp6, got %s", entries[0].Protocol)
	}
}
