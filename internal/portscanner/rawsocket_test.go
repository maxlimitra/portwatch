package portscanner

import (
	"os"
	"path/filepath"
	"testing"
)

const rawProcContent = `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0000 00000000:0000 07 00000000:00000000 00:00000000 00000000     0        0 12345 2 0000000000000000 0
   1: 00000000:0000 00000000:0000 07 00000000:00000000 00:00000000 00000000     0        0 67890 2 0000000000000000 0
`

func writeTempRaw(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "raw")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp raw: %v", err)
	}
	return p
}

func TestParseProcNetRawMissingFile(t *testing.T) {
	_, err := ParseProcNetRaw("/nonexistent/proc/net/raw")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseProcNetRaw6MissingFile(t *testing.T) {
	_, err := ParseProcNetRaw6("/nonexistent/proc/net/raw6")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseProcNetRawValidFile(t *testing.T) {
	path := writeTempRaw(t, rawProcContent)
	entries, err := ParseProcNetRaw(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one entry, got none")
	}
	// First entry: 0100007F -> 127.0.0.1
	got := entries[0].LocalIP.String()
	if got != "127.0.0.1" {
		t.Errorf("expected 127.0.0.1, got %s", got)
	}
}

func TestRawEntryString(t *testing.T) {
	path := writeTempRaw(t, rawProcContent)
	entries, err := ParseProcNetRaw(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) == 0 {
		t.Skip("no entries to test String()")
	}
	s := entries[0].String()
	if s == "" {
		t.Error("expected non-empty string from RawEntry.String()")
	}
}

func TestReadProcLinesSkipsHeader(t *testing.T) {
	content := "header line\ndata line one\ndata line two\n"
	dir := t.TempDir()
	p := filepath.Join(dir, "testfile")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	lines, err := readProcLines(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Errorf("expected 2 data lines, got %d", len(lines))
	}
}
