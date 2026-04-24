package alerting

import (
	"fmt"
	"log"
	"time"
)

// AlertLevel represents the severity of an alert.
type AlertLevel int

const (
	AlertInfo AlertLevel = iota
	AlertWarn
	AlertCritical
)

func (a AlertLevel) String() string {
	switch a {
	case AlertInfo:
		return "INFO"
	case AlertWarn:
		return "WARN"
	case AlertCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// Alert represents a single port-related alert event.
type Alert struct {
	Level     AlertLevel
	Port      int
	PID       int
	Process   string
	Message   string
	Timestamp time.Time
}

func (a Alert) String() string {
	return fmt.Sprintf("[%s] %s — port=%d pid=%d process=%q",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Port,
		a.PID,
		a.Process,
	)
}

// Handler is a function that receives and processes an Alert.
type Handler func(Alert)

// Alerter dispatches alerts to one or more handlers.
type Alerter struct {
	handlers []Handler
	logger   *log.Logger
}

// NewAlerter creates a new Alerter with the given handlers.
func NewAlerter(logger *log.Logger, handlers ...Handler) *Alerter {
	return &Alerter{
		handlers: handlers,
		logger:   logger,
	}
}

// Send dispatches an alert to all registered handlers.
func (a *Alerter) Send(alert Alert) {
	if alert.Timestamp.IsZero() {
		alert.Timestamp = time.Now()
	}
	a.logger.Println(alert.String())
	for _, h := range a.handlers {
		h(alert)
	}
}

// NewAlert is a convenience constructor for building Alert values.
func NewAlert(level AlertLevel, port, pid int, process, message string) Alert {
	return Alert{
		Level:     level,
		Port:      port,
		PID:       pid,
		Process:   process,
		Message:   message,
		Timestamp: time.Now(),
	}
}
