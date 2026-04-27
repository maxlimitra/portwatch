package metrics

import (
	"github.com/user/portwatch/internal/snapshot"
)

// Collector updates a Registry with metrics derived from a snapshot diff.
type Collector struct {
	reg           *Registry
	totalScans    Counter
	portsAdded    Counter
	portsRemoved  Counter
	activePortsG  Gauge
}

// NewCollector creates a Collector that records metrics into reg.
func NewCollector(reg *Registry) *Collector {
	return &Collector{
		reg:          reg,
		totalScans:   reg.Counter("portwatch_scans_total"),
		portsAdded:   reg.Counter("portwatch_ports_added_total"),
		portsRemoved: reg.Counter("portwatch_ports_removed_total"),
		activePortsG: reg.Gauge("portwatch_active_ports"),
	}
}

// RecordScan updates counters after each poll cycle.
// diff is the result of the latest snapshot comparison.
// activeCount is the total number of currently open ports.
func (c *Collector) RecordScan(diff snapshot.Diff, activeCount int) {
	c.totalScans.Inc()
	for range diff.Added {
		c.portsAdded.Inc()
	}
	for range diff.Removed {
		c.portsRemoved.Inc()
	}
	c.activePortsG.Set(float64(activeCount))
}
