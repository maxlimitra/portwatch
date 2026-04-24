package watcher

import (
	"time"

	"github.com/user/portwatch/internal/alerting"
	"github.com/user/portwatch/internal/portscanner"
)

// Config holds configuration for the Watcher.
type Config struct {
	Interval    time.Duration
	AllowedPorts map[int]bool
}

// Watcher periodically scans ports and emits alerts on unexpected bindings.
type Watcher struct {
	config  Config
	scanner *portscanner.Scanner
	alerter *alerting.Alerter
	prev    map[int]portscanner.PortEntry
	stop    chan struct{}
}

// NewWatcher creates a new Watcher with the given config, scanner, and alerter.
func NewWatcher(cfg Config, s *portscanner.Scanner, a *alerting.Alerter) *Watcher {
	return &Watcher{
		config:  cfg,
		scanner: s,
		alerter: a,
		prev:    make(map[int]portscanner.PortEntry),
		stop:    make(chan struct{}),
	}
}

// Start begins the watch loop in the current goroutine.
func (w *Watcher) Start() {
	ticker := time.NewTicker(w.config.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.scan()
		case <-w.stop:
			return
		}
	}
}

// Stop signals the watch loop to exit.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) scan() {
	entries, err := w.scanner.Scan()
	if err != nil {
		w.alerter.Send(alerting.NewAlert(alerting.LevelError, "scan failed: "+err.Error()))
		return
	}

	current := make(map[int]portscanner.PortEntry, len(entries))
	for _, e := range entries {
		current[e.Port] = e
	}

	// Detect newly bound ports.
	for port, entry := range current {
		if _, existed := w.prev[port]; !existed {
			level := alerting.LevelInfo
			msg := "new port binding detected"
			if !w.config.AllowedPorts[port] {
				level = alerting.LevelWarn
				msg = "unexpected port binding detected"
			}
			a := alerting.NewAlert(level, msg)
			a.Port = entry.Port
			a.Address = entry.Address
			w.alerter.Send(a)
		}
	}

	// Detect released ports.
	for port, entry := range w.prev {
		if _, stillOpen := current[port]; !stillOpen {
			a := alerting.NewAlert(alerting.LevelInfo, "port released")
			a.Port = entry.Port
			a.Address = entry.Address
			w.alerter.Send(a)
		}
	}

	w.prev = current
}
