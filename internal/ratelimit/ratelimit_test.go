package ratelimit_test

import (
	"testing"
	"time"

	"portwatch/internal/ratelimit"
)

func TestAllowFirstCallAlwaysPasses(t *testing.T) {
	l := ratelimit.NewLimiter(5 * time.Second)
	if !l.Allow("127.0.0.1:8080") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllowSecondCallWithinCooldownBlocked(t *testing.T) {
	l := ratelimit.NewLimiter(5 * time.Second)
	l.Allow("127.0.0.1:8080")
	if l.Allow("127.0.0.1:8080") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllowDifferentKeysAreIndependent(t *testing.T) {
	l := ratelimit.NewLimiter(5 * time.Second)
	l.Allow("127.0.0.1:8080")
	if !l.Allow("127.0.0.1:9090") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestAllowZeroCooldownNeverBlocks(t *testing.T) {
	l := ratelimit.NewLimiter(0)
	for i := 0; i < 10; i++ {
		if !l.Allow("key") {
			t.Fatalf("expected zero-cooldown limiter to always allow (iteration %d)", i)
		}
	}
}

func TestResetAllowsKeyAgain(t *testing.T) {
	l := ratelimit.NewLimiter(5 * time.Second)
	l.Allow("key")
	l.Reset("key")
	if !l.Allow("key") {
		t.Fatal("expected key to be allowed after Reset")
	}
}

func TestFlushClearsAllState(t *testing.T) {
	l := ratelimit.NewLimiter(5 * time.Second)
	l.Allow("a")
	l.Allow("b")
	l.Flush()
	if !l.Allow("a") || !l.Allow("b") {
		t.Fatal("expected all keys to be allowed after Flush")
	}
}

func TestAllowAfterCooldownExpiry(t *testing.T) {
	l := ratelimit.NewLimiter(10 * time.Millisecond)
	l.Allow("key")
	time.Sleep(20 * time.Millisecond)
	if !l.Allow("key") {
		t.Fatal("expected key to be allowed after cooldown expired")
	}
}
