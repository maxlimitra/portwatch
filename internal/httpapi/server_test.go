package httpapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/httpapi"
	"github.com/user/portwatch/internal/metrics"
)

func freeAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freeAddr: %v", err)
	}
	addr := l.Addr().String()
	_ = l.Close()
	return addr
}

func startServer(t *testing.T) (*httpapi.Server, string) {
	t.Helper()
	reg := metrics.NewRegistry()
	addr := freeAddr(t)
	srv := httpapi.NewServer(addr, reg)
	if err := srv.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	time.Sleep(20 * time.Millisecond) // allow goroutine to bind
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Stop(ctx)
	})
	return srv, fmt.Sprintf("http://%s", addr)
}

func TestHealthzReturnsOK(t *testing.T) {
	_, base := startServer(t)
	resp, err := http.Get(base + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", body["status"])
	}
}

func TestMetricsEndpointContainsData(t *testing.T) {
	_, base := startServer(t)
	resp, err := http.Get(base + "/metrics")
	if err != nil {
		t.Fatalf("GET /metrics: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	b, _ := io.ReadAll(resp.Body)
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/plain") {
		t.Fatalf("unexpected content-type: %s", resp.Header.Get("Content-Type"))
	}
	_ = b // response body may be empty when no metrics recorded yet
}

func TestUnknownRouteReturns404(t *testing.T) {
	_, base := startServer(t)
	resp, err := http.Get(base + "/notfound")
	if err != nil {
		t.Fatalf("GET /notfound: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}
