package portscanner

import (
	"net"
	"testing"
)

func TestScanICMPReturnsWithoutError(t *testing.T) {
	// Use a temp proc root that has no icmp files — ScanICMP must not error.
	dir := t.TempDir()
	s := &Scanner{procRoot: dir}

	entries, err := s.ScanICMP()
	if err != nil {
		t.Fatalf("ScanICMP returned unexpected error: %v", err)
	}
	// No files present, so we expect an empty (but non-nil) slice.
	if entries == nil {
		t.Error("expected non-nil slice, got nil")
	}
}

func TestICMPEntryToPortEntry(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	e := ICMPEntry{
		LocalIP:   ip,
		LocalPort: 8080,
		State:     "0A",
		Inode:     "42",
		Protocol:  "icmp",
	}

	p := ICMPEntryToPortEntry(e)

	if p.Protocol != "icmp" {
		t.Errorf("expected protocol icmp, got %s", p.Protocol)
	}
	if p.LocalPort != 8080 {
		t.Errorf("expected port 8080, got %d", p.LocalPort)
	}
	if p.State != "0A" {
		t.Errorf("expected state 0A, got %s", p.State)
	}
	if p.Inode != "42" {
		t.Errorf("expected inode 42, got %s", p.Inode)
	}
	if !p.LocalIP.Equal(ip) {
		t.Errorf("expected IP %s, got %s", ip, p.LocalIP)
	}
}

func TestPortEntryString(t *testing.T) {
	p := PortEntry{
		Protocol:  "icmp",
		LocalIP:   net.ParseIP("0.0.0.0"),
		LocalPort: 0,
		State:     "0A",
		Inode:     "1",
	}
	s := p.String()
	if s == "" {
		t.Error("PortEntry.String() returned empty string")
	}
}
