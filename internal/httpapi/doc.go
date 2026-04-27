// Package httpapi exposes a small HTTP interface for portwatch.
//
// Endpoints:
//
//	 GET /healthz  – liveness probe; returns {"status":"ok"} with HTTP 200.
//	 GET /metrics  – Prometheus-compatible plain-text metrics snapshot
//	                 sourced from the shared metrics.Registry.
//
// Usage:
//
//	srv := httpapi.NewServer(":9090", registry)
//	if err := srv.Start(); err != nil {
//	    log.Fatal(err)
//	}
//	defer srv.Stop(ctx)
package httpapi
