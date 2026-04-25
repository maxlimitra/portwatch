// Package metrics provides a minimal, dependency-free metrics registry for
// portwatch.
//
// # Overview
//
// The registry tracks named [Counter] and [Gauge] values using atomic
// operations, making it safe to update from multiple goroutines without
// additional locking on the caller side.
//
// # Usage
//
//	reg := metrics.NewRegistry()
//
//	// Increment a counter each scan cycle.
//	reg.Counter("scan_cycles").Inc()
//
//	// Track the number of currently-open ports as a gauge.
//	reg.Gauge("ports_tracked").Set(int64(len(ports)))
//
//	// Dump a human-readable snapshot to stdout.
//	reg.WriteTo(os.Stdout)
//
// Metric names are arbitrary strings; the same name always returns the same
// underlying instance within a registry.
package metrics
