// Package history records port binding events over time, enabling
// trend analysis and persistent audit trails across daemon restarts.
package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single recorded port event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Proto     string    `json:"proto"`
	Port      uint16    `json:"port"`
	Addr      string    `json:"addr"`
	Process   string    `json:"process,omitempty"`
	Event     string    `json:"event"` // "opened" | "closed"
}

// Recorder appends entries to an on-disk JSONL file and keeps a
// bounded in-memory ring buffer for fast recent-event queries.
type Recorder struct {
	mu      sync.Mutex
	path    string
	buf     []Entry
	cap     int
}

// NewRecorder creates a Recorder that writes to path and keeps at
// most bufCap entries in memory. bufCap must be > 0.
func NewRecorder(path string, bufCap int) (*Recorder, error) {
	if bufCap <= 0 {
		bufCap = 256
	}
	r := &Recorder{path: path, cap: bufCap}
	// Ensure the file is writable / creatable.
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	f.Close()
	return r, nil
}

// Record appends an entry to the JSONL file and the in-memory buffer.
func (r *Recorder) Record(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	// Append to ring buffer.
	if len(r.buf) >= r.cap {
		r.buf = r.buf[1:]
	}
	r.buf = append(r.buf, e)

	// Persist.
	f, err := os.OpenFile(r.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(e)
}

// Recent returns up to n most-recent in-memory entries.
func (r *Recorder) Recent(n int) []Entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	if n <= 0 || n > len(r.buf) {
		n = len(r.buf)
	}
	result := make([]Entry, n)
	copy(result, r.buf[len(r.buf)-n:])
	return result
}
