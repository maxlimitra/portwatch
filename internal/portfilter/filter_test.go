package portfilter_test

import (
	"testing"

	"github.com/user/portwatch/internal/portfilter"
	"github.com/user/portwatch/internal/portscanner"
)

func makeEntry(port uint16, proto string) portscanner.PortEntry {
	return portscanner.PortEntry{Port: port, Protocol: proto, LocalAddr: "127.0.0.1"}
}

func TestFilterEmptyRulesRetainsAll(t *testing.T) {
	f := portfilter.NewFilter(nil)
	entries := []portscanner.PortEntry{makeEntry(80, "tcp"), makeEntry(443, "tcp")}
	got := f.Apply(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestFilterDenySpecificPort(t *testing.T) {
	rules := []portfilter.Rule{
		{Port: 80, Protocol: "tcp", Allow: false},
	}
	f := portfilter.NewFilter(rules)
	entries := []portscanner.PortEntry{makeEntry(80, "tcp"), makeEntry(443, "tcp")}
	got := f.Apply(entries)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].Port != 443 {
		t.Errorf("expected port 443, got %d", got[0].Port)
	}
}

func TestFilterAllowOnlySpecificProtocol(t *testing.T) {
	rules := []portfilter.Rule{
		{Port: 0, Protocol: "udp", Allow: false},
	}
	f := portfilter.NewFilter(rules)
	entries := []portscanner.PortEntry{
		makeEntry(53, "udp"),
		makeEntry(80, "tcp"),
		makeEntry(123, "udp"),
	}
	got := f.Apply(entries)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].Protocol != "tcp" {
		t.Errorf("expected tcp, got %s", got[0].Protocol)
	}
}

func TestFilterFirstRuleWins(t *testing.T) {
	rules := []portfilter.Rule{
		{Port: 80, Protocol: "tcp", Allow: true},
		{Port: 0, Protocol: "tcp", Allow: false},
	}
	f := portfilter.NewFilter(rules)
	entries := []portscanner.PortEntry{
		makeEntry(80, "tcp"),
		makeEntry(443, "tcp"),
	}
	got := f.Apply(entries)
	// port 80 matches first rule (allow), port 443 matches second rule (deny)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].Port != 80 {
		t.Errorf("expected port 80, got %d", got[0].Port)
	}
}

func TestFilterNoMatchDefaultRetain(t *testing.T) {
	rules := []portfilter.Rule{
		{Port: 9999, Protocol: "tcp", Allow: false},
	}
	f := portfilter.NewFilter(rules)
	entries := []portscanner.PortEntry{makeEntry(80, "tcp")}
	got := f.Apply(entries)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry retained, got %d", len(got))
	}
}
