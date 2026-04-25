// Package history provides persistent audit-trail recording for
// portwatch.
//
// # Overview
//
// A [Recorder] writes every port open/close event to a newline-delimited
// JSON (JSONL) file so that events survive daemon restarts. It also
// maintains a bounded in-memory ring buffer that allows callers to
// retrieve the most recent N events without touching the disk.
//
// # Usage
//
//	rec, err := history.NewRecorder("/var/lib/portwatch/events.jsonl", 512)
//	if err != nil { ... }
//
//	rec.Record(history.Entry{
//		Proto:   "tcp",
//		Port:    8080,
//		Process: "nginx",
//		Event:   "opened",
//	})
//
//	recent := rec.Recent(10)
package history
