package watcher

import (
	"fmt"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/snapshot"
)

// historyRecorder is the subset of history.Recorder used by the watcher
// so it can be mocked in tests.
type historyRecorder interface {
	Record(history.Entry) error
}

// recordDiff converts snapshot diff results into history entries and
// forwards them to the provided recorder. Errors are non-fatal; they
// are returned so callers can log them if desired.
func recordDiff(rec historyRecorder, added, removed []snapshot.PortEntry) []error {
	var errs []error
	for _, p := range added {
		err := rec.Record(history.Entry{
			Proto:   p.Proto,
			Port:    p.Port,
			Addr:    fmt.Sprintf("%s:%d", p.LocalAddr, p.Port),
			Process: p.Process,
			Event:   "opened",
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, p := range removed {
		err := rec.Record(history.Entry{
			Proto:   p.Proto,
			Port:    p.Port,
			Addr:    fmt.Sprintf("%s:%d", p.LocalAddr, p.Port),
			Process: p.Process,
			Event:   "closed",
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
