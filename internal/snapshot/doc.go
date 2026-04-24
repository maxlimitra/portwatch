// Package snapshot provides persistence and diffing for port scan results.
//
// A Snapshot captures the set of active port bindings observed at a single
// point in time. The Store type serialises snapshots to a JSON file so that
// portwatch can survive restarts without raising spurious alerts for ports
// that were already open before the daemon started.
//
// Diff compares two snapshots and returns the ports that appeared or
// disappeared between them, allowing the watcher to emit targeted alerts
// rather than re-alerting on every known binding each cycle.
package snapshot
