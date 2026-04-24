package alerting

import (
	"fmt"

	"portwatch/internal/ratelimit"
)

// RateLimitedAlerter wraps an Alerter and suppresses duplicate alerts
// for the same port within the limiter's cooldown window.
type RateLimitedAlerter struct {
	inner   *Alerter
	limiter *ratelimit.Limiter
}

// NewRateLimitedAlerter creates a RateLimitedAlerter that delegates to
// inner while using limiter to deduplicate alerts by port key.
func NewRateLimitedAlerter(inner *Alerter, limiter *ratelimit.Limiter) *RateLimitedAlerter {
	return &RateLimitedAlerter{inner: inner, limiter: limiter}
}

// Send dispatches the alert only if the rate limiter allows it.
// The deduplication key is composed of the alert's IP, port and level.
func (r *RateLimitedAlerter) Send(a Alert) {
	key := fmt.Sprintf("%s:%d:%s", a.IP, a.Port, a.Level)
	if !r.limiter.Allow(key) {
		return
	}
	r.inner.Send(a)
}
