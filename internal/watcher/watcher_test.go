package watcher_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alerting"
	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/watcher"
)

// fakeScanner satisfies portscanner.Scanner's Scan interface via embedding.
type fakeScanner struct {
	mu      sync.Mutex
	entries []portscanner.PortEntry
}

func (f *fakeScanner) set(entries []portscanner.PortEntry) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.entries = entries
}

func (f *fakeScanner) Scan() ([]portscanner.PortEntry, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.entries, nil
}

func collectingAlerter() (*alerting.Alerter, *[]alerting.Alert) {
	var mu sync.Mutex
	var collected []alerting.Alert
	h := func(a alerting.Alert) {
		mu.Lock()
		defer mu.Unlock()
		collected = append(collected, a)
	}
	a := alerting.NewAlerter(alerting.HandlerFunc(h))
	return a, &collected
}

func TestWatcherDetectsNewPort(t *testing.T) {
	fs := &fakeScanner{}
	a, alerts := collectingAlerter()

	cfg := watcher.Config{
		Interval:     10 * time.Millisecond,
		AllowedPorts: map[int]bool{80: true},
	}
	w := watcher.NewWatcher(cfg, portscanner.ScannerFunc(fs.Scan), a)

	fs.set([]portscanner.PortEntry{{Port: 9999, Address: "0.0.0.0"}})

	go w.Start()
	time.Sleep(50 * time.Millisecond)
	w.Stop()

	if len(*alerts) == 0 {
		t.Fatal("expected at least one alert for unexpected port")
	}
	found := false
	for _, al := range *alerts {
		if al.Port == 9999 && al.Level == alerting.LevelWarn {
			found = true
		}
	}
	if !found {
		t.Errorf("expected warn alert for port 9999, got %+v", *alerts)
	}
}

func TestWatcherAllowedPortIsInfo(t *testing.T) {
	fs := &fakeScanner{}
	a, alerts := collectingAlerter()

	cfg := watcher.Config{
		Interval:     10 * time.Millisecond,
		AllowedPorts: map[int]bool{80: true},
	}
	w := watcher.NewWatcher(cfg, portscanner.ScannerFunc(fs.Scan), a)

	fs.set([]portscanner.PortEntry{{Port: 80, Address: "0.0.0.0"}})

	go w.Start()
	time.Sleep(50 * time.Millisecond)
	w.Stop()

	for _, al := range *alerts {
		if al.Port == 80 && al.Level == alerting.LevelWarn {
			t.Errorf("allowed port 80 should not produce a warn alert")
		}
	}
}
