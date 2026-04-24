// Package ratelimit provides alert rate-limiting to suppress duplicate
// alerts for the same port within a configurable cooldown window.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter suppresses repeated alerts for the same key within a cooldown
// duration. It is safe for concurrent use.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
}

// NewLimiter creates a Limiter with the given cooldown window.
// A cooldown of zero disables rate-limiting (all events pass through).
func NewLimiter(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
	}
}

// Allow reports whether an alert identified by key should be dispatched.
// It returns true the first time a key is seen and again only after the
// cooldown period has elapsed since the previous allowed event.
func (l *Limiter) Allow(key string) bool {
	if l.cooldown <= 0 {
		return true
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if t, ok := l.last[key]; ok && now.Sub(t) < l.cooldown {
		return false
	}
	l.last[key] = now
	return true
}

// Reset clears the rate-limit state for a specific key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// Flush clears all recorded state, effectively resetting the limiter.
func (l *Limiter) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}
