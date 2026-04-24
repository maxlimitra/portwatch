package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/snapshot"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snap.json")
}

func TestStoreRoundTrip(t *testing.T) {
	store := snapshot.NewStore(tempPath(t))

	snap := snapshot.Snapshot{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Ports: []portscanner.PortInfo{
			{Address: "0.0.0.0:8080", Protocol: "tcp", PID: 42},
		},
	}

	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.Ports) != 1 || loaded.Ports[0].Address != "0.0.0.0:8080" {
		t.Errorf("unexpected ports: %+v", loaded.Ports)
	}
}

func TestStoreLoadMissingFile(t *testing.T) {
	store := snapshot.NewStore(filepath.Join(t.TempDir(), "nonexistent.json"))
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty snapshot, got %+v", snap)
	}
}

func TestStoreLoadCorrupt(t *testing.T) {
	p := tempPath(t)
	if err := os.WriteFile(p, []byte("not-json{"), 0644); err != nil {
		t.Fatal(err)
	}
	store := snapshot.NewStore(p)
	_, err := store.Load()
	if err == nil {
		t.Fatal("expected error for corrupt JSON")
	}
}

func TestDiff(t *testing.T) {
	portA := portscanner.PortInfo{Address: "0.0.0.0:80", Protocol: "tcp"}
	portB := portscanner.PortInfo{Address: "0.0.0.0:443", Protocol: "tcp"}
	portC := portscanner.PortInfo{Address: "0.0.0.0:8080", Protocol: "tcp"}

	prev := snapshot.Snapshot{Ports: []portscanner.PortInfo{portA, portB}}
	curr := snapshot.Snapshot{Ports: []portscanner.PortInfo{portB, portC}}

	added, removed := snapshot.Diff(prev, curr)

	if len(added) != 1 || added[0].Address != portC.Address {
		t.Errorf("expected added=[%s], got %+v", portC.Address, added)
	}
	if len(removed) != 1 || removed[0].Address != portA.Address {
		t.Errorf("expected removed=[%s], got %+v", portA.Address, removed)
	}
}

func TestDiffNoChange(t *testing.T) {
	port := portscanner.PortInfo{Address: "0.0.0.0:80", Protocol: "tcp"}
	snap := snapshot.Snapshot{Ports: []portscanner.PortInfo{port}}

	added, removed := snapshot.Diff(snap, snap)
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", added, removed)
	}
}
