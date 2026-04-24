package alerting

import (
	"bytes"
	"log"
	"testing"
	"time"
)

func TestAlertLevelString(t *testing.T) {
	cases := []struct {
		level    AlertLevel
		expected string
	}{
		{AlertInfo, "INFO"},
		{AlertWarn, "WARN"},
		{AlertCritical, "CRITICAL"},
		{AlertLevel(99), "UNKNOWN"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.expected {
			t.Errorf("AlertLevel(%d).String() = %q, want %q", tc.level, got, tc.expected)
		}
	}
}

func TestNewAlert(t *testing.T) {
	before := time.Now()
	a := NewAlert(AlertWarn, 8080, 1234, "nginx", "unexpected binding")
	after := time.Now()

	if a.Level != AlertWarn {
		t.Errorf("expected level WARN, got %s", a.Level)
	}
	if a.Port != 8080 {
		t.Errorf("expected port 8080, got %d", a.Port)
	}
	if a.PID != 1234 {
		t.Errorf("expected pid 1234, got %d", a.PID)
	}
	if a.Process != "nginx" {
		t.Errorf("expected process nginx, got %s", a.Process)
	}
	if a.Timestamp.Before(before) || a.Timestamp.After(after) {
		t.Error("timestamp out of expected range")
	}
}

func TestAlerterSendDispatchesHandlers(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	received := make([]Alert, 0)
	handler := func(a Alert) {
		received = append(received, a)
	}

	alerter := NewAlerter(logger, handler)
	alert := NewAlert(AlertCritical, 443, 999, "unknown", "port conflict")
	alerter.Send(alert)

	if len(received) != 1 {
		t.Fatalf("expected 1 alert dispatched, got %d", len(received))
	}
	if received[0].Port != 443 {
		t.Errorf("expected port 443, got %d", received[0].Port)
	}
	if buf.Len() == 0 {
		t.Error("expected logger output, got none")
	}
}

func TestAlerterSendSetsTimestampIfZero(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	var got Alert
	alerter := NewAlerter(logger, func(a Alert) { got = a })

	alerter.Send(Alert{Level: AlertInfo, Port: 80})

	if got.Timestamp.IsZero() {
		t.Error("expected timestamp to be set automatically")
	}
}
