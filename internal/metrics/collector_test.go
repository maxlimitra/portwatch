package metrics

import (
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func makeEntry(ip string, port uint16, proto string) snapshot.Entry {
	return snapshot.Entry{IP: ip, Port: port, Protocol: proto}
}

func TestCollectorRecordScanIncrementsTotalScans(t *testing.T) {
	reg := NewRegistry()
	col := NewCollector(reg)

	col.RecordScan(snapshot.Diff{}, 0)
	col.RecordScan(snapshot.Diff{}, 0)

	if got := reg.Counter("portwatch_scans_total").Value(); got != 2 {
		t.Fatalf("expected 2 scans, got %d", got)
	}
}

func TestCollectorRecordScanCountsAddedPorts(t *testing.T) {
	reg := NewRegistry()
	col := NewCollector(reg)

	diff := snapshot.Diff{
		Added: []snapshot.Entry{
			makeEntry("127.0.0.1", 8080, "tcp"),
			makeEntry("127.0.0.1", 9090, "tcp"),
		},
	}
	col.RecordScan(diff, 2)

	if got := reg.Counter("portwatch_ports_added_total").Value(); got != 2 {
		t.Fatalf("expected 2 added, got %d", got)
	}
}

func TestCollectorRecordScanCountsRemovedPorts(t *testing.T) {
	reg := NewRegistry()
	col := NewCollector(reg)

	diff := snapshot.Diff{
		Removed: []snapshot.Entry{
			makeEntry("0.0.0.0", 22, "tcp"),
		},
	}
	col.RecordScan(diff, 0)

	if got := reg.Counter("portwatch_ports_removed_total").Value(); got != 1 {
		t.Fatalf("expected 1 removed, got %d", got)
	}
}

func TestCollectorRecordScanSetsActiveGauge(t *testing.T) {
	reg := NewRegistry()
	col := NewCollector(reg)

	col.RecordScan(snapshot.Diff{}, 42)

	if got := reg.Gauge("portwatch_active_ports").Value(); got != 42 {
		t.Fatalf("expected gauge 42, got %v", got)
	}
}

func TestCollectorAccumulatesAcrossScans(t *testing.T) {
	reg := NewRegistry()
	col := NewCollector(reg)

	col.RecordScan(snapshot.Diff{Added: []snapshot.Entry{makeEntry("127.0.0.1", 80, "tcp")}}, 1)
	col.RecordScan(snapshot.Diff{Added: []snapshot.Entry{makeEntry("127.0.0.1", 443, "tcp")}}, 2)

	if got := reg.Counter("portwatch_ports_added_total").Value(); got != 2 {
		t.Fatalf("expected cumulative 2 added, got %d", got)
	}
	if got := reg.Counter("portwatch_scans_total").Value(); got != 2 {
		t.Fatalf("expected 2 total scans, got %d", got)
	}
}
