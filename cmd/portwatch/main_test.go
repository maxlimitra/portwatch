package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestMainHelp verifies the binary prints usage when --help is passed.
func TestMainHelp(t *testing.T) {
	if os.Getenv("PORTWATCH_INTEGRATION") != "1" {
		t.Skip("skipping integration test; set PORTWATCH_INTEGRATION=1 to run")
	}

	cmd := exec.Command("go", "run", ".", "--help")
	cmd.Dir = filepath.Join(".")
	out, _ := cmd.CombinedOutput()
	if len(out) == 0 {
		t.Error("expected help output, got nothing")
	}
}

// TestMainInvalidConfig ensures the binary exits non-zero for a bad config path.
func TestMainInvalidConfig(t *testing.T) {
	if os.Getenv("PORTWATCH_INTEGRATION") != "1" {
		t.Skip("skipping integration test; set PORTWATCH_INTEGRATION=1 to run")
	}

	cmd := exec.Command("go", "run", ".", "-config", "/nonexistent/path/config.yaml")
	cmd.Dir = filepath.Join(".")
	err := cmd.Run()
	if err == nil {
		t.Error("expected non-zero exit for invalid config, got nil")
	}
}

// TestMainRunsAndStops verifies the daemon starts and shuts down cleanly on SIGTERM.
func TestMainRunsAndStops(t *testing.T) {
	if os.Getenv("PORTWATCH_INTEGRATION") != "1" {
		t.Skip("skipping integration test; set PORTWATCH_INTEGRATION=1 to run")
	}

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = filepath.Join(".")
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start portwatch: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatalf("failed to send interrupt: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case <-done:
		// exited cleanly
	case <-time.After(3 * time.Second):
		cmd.Process.Kill()
		t.Error("portwatch did not exit within timeout after SIGINT")
	}
}
