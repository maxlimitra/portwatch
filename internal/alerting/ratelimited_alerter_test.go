package alerting_test

import (
	"sync"
	"testing"
	"time"

	"portwatch/internal/alerting"
	"portwatch/internal/ratelimit"
)

type countingHandler struct {
	mu    sync.Mutex
	count int
}

func (c *countingHandler) Handle(a alerting.Alert) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *countingHandler) total() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func TestRateLimitedAlerterSuppressesDuplicates(t *testing.T) {
	handler := &countingHandler{}
	inner := alerting.NewAlerter(handler.Handle)
	limiter := ratelimit.NewLimiter(5 * time.Second)
	ra := alerting.NewRateLimitedAlerter(inner, limiter)

	a := alerting.NewAlert(alerting.LevelWarn, "127.0.0.1", 8080, "test")
	ra.Send(a)
	ra.Send(a)
	ra.Send(a)

	if got := handler.total(); got != 1 {
		t.Fatalf("expected 1 dispatch, got %d", got)
	}
}

func TestRateLimitedAlerterAllowsDifferentPorts(t *testing.T) {
	handler := &countingHandler{}
	inner := alerting.NewAlerter(handler.Handle)
	limiter := ratelimit.NewLimiter(5 * time.Second)
	ra := alerting.NewRateLimitedAlerter(inner, limiter)

	ra.Send(alerting.NewAlert(alerting.LevelWarn, "127.0.0.1", 8080, "test"))
	ra.Send(alerting.NewAlert(alerting.LevelWarn, "127.0.0.1", 9090, "test"))

	if got := handler.total(); got != 2 {
		t.Fatalf("expected 2 dispatches, got %d", got)
	}
}

func TestRateLimitedAlerterPassesAfterCooldown(t *testing.T) {
	handler := &countingHandler{}
	inner := alerting.NewAlerter(handler.Handle)
	limiter := ratelimit.NewLimiter(10 * time.Millisecond)
	ra := alerting.NewRateLimitedAlerter(inner, limiter)

	a := alerting.NewAlert(alerting.LevelWarn, "127.0.0.1", 8080, "test")
	ra.Send(a)
	time.Sleep(20 * time.Millisecond)
	ra.Send(a)

	if got := handler.total(); got != 2 {
		t.Fatalf("expected 2 dispatches after cooldown, got %d", got)
	}
}
