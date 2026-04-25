package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "events.jsonl")
}

func TestNewRecorderCreatesFile(t *testing.T) {
	p := tempFile(t)
	_, err := NewRecorder(p, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestNewRecorderInvalidPath(t *testing.T) {
	_, err := NewRecorder("/no/such/dir/events.jsonl", 10)
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestRecordPersistsToFile(t *testing.T) {
	p := tempFile(t)
	rec, _ := NewRecorder(p, 10)

	e := Entry{Proto: "tcp", Port: 8080, Event: "opened"}
	if err := rec.Record(e); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	data, _ := os.ReadFile(p)
	var got Entry
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("could not unmarshal persisted entry: %v", err)
	}
	if got.Port != 8080 || got.Proto != "tcp" {
		t.Errorf("unexpected entry: %+v", got)
	}
}

func TestRecordSetsTimestampIfZero(t *testing.T) {
	rec, _ := NewRecorder(tempFile(t), 10)
	e := Entry{Port: 443, Event: "opened"}
	before := time.Now().UTC()
	_ = rec.Record(e)
	recent := rec.Recent(1)
	if recent[0].Timestamp.Before(before) {
		t.Error("timestamp should have been set to approximately now")
	}
}

func TestRecentReturnsLatestEntries(t *testing.T) {
	rec, _ := NewRecorder(tempFile(t), 5)
	for i := uint16(1); i <= 7; i++ {
		_ = rec.Record(Entry{Port: i, Event: "opened"})
	}
	got := rec.Recent(3)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	// Ring buffer capped at 5; most recent 3 should be ports 5,6,7.
	if got[0].Port != 5 || got[2].Port != 7 {
		t.Errorf("unexpected ports: %v %v", got[0].Port, got[2].Port)
	}
}

func TestRecentAllWhenNLargerThanBuffer(t *testing.T) {
	rec, _ := NewRecorder(tempFile(t), 10)
	for i := uint16(1); i <= 3; i++ {
		_ = rec.Record(Entry{Port: i, Event: "opened"})
	}
	got := rec.Recent(100)
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
}
