// Package metrics provides lightweight in-process counters and gauges
// that portwatch exposes for observability (e.g. number of alerts fired,
// ports currently tracked, scan cycles completed).
package metrics

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"sync/atomic"
)

// Counter is a monotonically increasing integer counter.
type Counter struct{ v uint64 }

// Inc increments the counter by 1.
func (c *Counter) Inc() { atomic.AddUint64(&c.v, 1) }

// Add increments the counter by delta.
func (c *Counter) Add(delta uint64) { atomic.AddUint64(&c.v, delta) }

// Value returns the current counter value.
func (c *Counter) Value() uint64 { return atomic.LoadUint64(&c.v) }

// Gauge is an integer value that can go up or down.
type Gauge struct{ v int64 }

// Set sets the gauge to an absolute value.
func (g *Gauge) Set(v int64) { atomic.StoreInt64(&g.v, v) }

// Value returns the current gauge value.
func (g *Gauge) Value() int64 { return atomic.LoadInt64(&g.v) }

// Registry holds named counters and gauges.
type Registry struct {
	mu       sync.RWMutex
	counters map[string]*Counter
	gauges   map[string]*Gauge
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		counters: make(map[string]*Counter),
		gauges:   make(map[string]*Gauge),
	}
}

// Counter returns (creating if necessary) the named counter.
func (r *Registry) Counter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.counters[name]; ok {
		return c
	}
	c := &Counter{}
	r.counters[name] = c
	return c
}

// Gauge returns (creating if necessary) the named gauge.
func (r *Registry) Gauge(name string) *Gauge {
	r.mu.Lock()
	defer r.mu.Unlock()
	if g, ok := r.gauges[name]; ok {
		return g
	}
	g := &Gauge{}
	r.gauges[name] = g
	return g
}

// WriteTo writes a human-readable snapshot of all metrics to w.
func (r *Registry) WriteTo(w io.Writer) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cnames := make([]string, 0, len(r.counters))
	for k := range r.counters {
		cnames = append(cnames, k)
	}
	sort.Strings(cnames)
	for _, name := range cnames {
		fmt.Fprintf(w, "counter %s %d\n", name, r.counters[name].Value())
	}

	gnames := make([]string, 0, len(r.gauges))
	for k := range r.gauges {
		gnames = append(gnames, k)
	}
	sort.Strings(gnames)
	for _, name := range gnames {
		fmt.Fprintf(w, "gauge   %s %d\n", name, r.gauges[name].Value())
	}
}
