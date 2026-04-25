package watcher

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/snapshot"
)

type mockRecorder struct {
	entries []history.Entry
	failOn  int // fail after this many successful records (-1 = never)
	calls   int
}

func (m *mockRecorder) Record(e history.Entry) error {
	m.calls++
	if m.failOn >= 0 && m.calls > m.failOn {
		return errors.New("mock record error")
	}
	m.entries = append(m.entries, e)
	return nil
}

func TestRecordDiffAdded(t *testing.T) {
	rec := &mockRecorder{failOn: -1}
	added := []snapshot.PortEntry{{Proto: "tcp", Port: 8080, LocalAddr: "0.0.0.0"}}
	errs := recordDiff(rec, added, nil)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(rec.entries) != 1 || rec.entries[0].Event != "opened" {
		t.Errorf("expected one 'opened' entry, got %+v", rec.entries)
	}
}

func TestRecordDiffRemoved(t *testing.T) {
	rec := &mockRecorder{failOn: -1}
	removed := []snapshot.PortEntry{{Proto: "udp", Port: 53, LocalAddr: "127.0.0.1"}}
	errs := recordDiff(rec, nil, removed)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(rec.entries) != 1 || rec.entries[0].Event != "closed" {
		t.Errorf("expected one 'closed' entry, got %+v", rec.entries)
	}
}

func TestRecordDiffRecordError(t *testing.T) {
	rec := &mockRecorder{failOn: 0} // fail immediately
	added := []snapshot.PortEntry{{Proto: "tcp", Port: 443}}
	errs := recordDiff(rec, added, nil)
	if len(errs) == 0 {
		t.Fatal("expected an error when recorder fails")
	}
}

func TestRecordDiffAddrFormat(t *testing.T) {
	rec := &mockRecorder{failOn: -1}
	added := []snapshot.PortEntry{{Proto: "tcp", Port: 9090, LocalAddr: "10.0.0.1", Process: "myapp"}}
	_ = recordDiff(rec, added, nil)
	if rec.entries[0].Addr != "10.0.0.1:9090" {
		t.Errorf("unexpected addr format: %q", rec.entries[0].Addr)
	}
	if rec.entries[0].Process != "myapp" {
		t.Errorf("process not forwarded: %q", rec.entries[0].Process)
	}
}
