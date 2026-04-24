package alerting

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// StdoutHandler returns a Handler that writes human-readable alerts to stdout.
func StdoutHandler() Handler {
	return func(a Alert) {
		fmt.Println(a.String())
	}
}

// JSONHandler returns a Handler that writes alerts as newline-delimited JSON
// to the given writer. Useful for log aggregation pipelines.
func JSONHandler(w io.Writer) Handler {
	encoder := json.NewEncoder(w)
	return func(a Alert) {
		payload := struct {
			Timestamp string `json:"timestamp"`
			Level     string `json:"level"`
			Port      int    `json:"port"`
			PID       int    `json:"pid"`
			Process   string `json:"process"`
			Message   string `json:"message"`
		}{
			Timestamp: a.Timestamp.UTC().Format("2006-01-02T15:04:05Z"),
			Level:     a.Level.String(),
			Port:      a.Port,
			PID:       a.PID,
			Process:   a.Process,
			Message:   a.Message,
		}
		// Best-effort; errors are silently dropped for non-critical log paths.
		_ = encoder.Encode(payload)
	}
}

// FileHandler returns a Handler that appends alerts as newline-delimited JSON
// to the file at the given path, creating it if it does not exist.
func FileHandler(path string) (Handler, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("alerting: open log file %q: %w", path, err)
	}
	return JSONHandler(f), nil
}
