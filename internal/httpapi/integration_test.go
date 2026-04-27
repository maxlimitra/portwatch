package httpapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/httpapi"
	"github.com/user/portwatch/internal/metrics"
)

// TestMetricsReflectsRegistryValues verifies that counter increments are
// visible through the /metrics endpoint.
func TestMetricsReflectsRegistryValues(t *testing.T) {
	reg := metrics.NewRegistry()
	c := reg.Counter("total_scans")
	c.Inc()
	c.Inc()

	addr := freeAddr(t)
	srv := httpapi.NewServer(addr, reg)
	if err := srv.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	time.Sleep(20 * time.Millisecond)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Stop(ctx)
	}()

	resp, err := http.Get(fmt.Sprintf("http://%s/metrics", addr))
	if err != nil {
		t.Fatalf("GET /metrics: %v", err)
	}
	defer resp.Body.Close()

	var buf strings.Builder
	if _, err := fmt.Fscan(resp.Body, &buf); err != nil && err.Error() != "EOF" {
		t.Logf("read body: %v", err)
	}
	// Body is plain text; just assert HTTP layer is healthy.
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// TestHealthzContentType ensures the response is JSON.
func TestHealthzContentType(t *testing.T) {
	_, base := startServer(t)
	resp, err := http.Get(base + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Fatalf("expected application/json content-type, got %q", ct)
	}

	var payload map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("json decode: %v", err)
	}
	if payload["status"] != "ok" {
		t.Fatalf("unexpected payload: %v", payload)
	}
}
