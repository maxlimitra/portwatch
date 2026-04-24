package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alerting"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/watcher"
)

func main() {
	configPath := flag.String("config", "", "path to config file (optional)")
	flag.Parse()

	var cfg config.Config
	var err error

	if *configPath != "" {
		cfg, err = config.LoadConfig(*configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
			os.Exit(1)
		}
	} else {
		cfg = config.DefaultConfig()
	}

	alerter := alerting.NewAlerter()

	switch cfg.OutputFormat {
	case "json":
		alerter.AddHandler(alerting.JSONHandler(os.Stdout))
	default:
		alerter.AddHandler(alerting.StdoutHandler(os.Stdout))
	}

	if cfg.LogFile != "" {
		f, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening log file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		alerter.AddHandler(alerting.FileHandler(f))
	}

	scanner := portscanner.NewScanner()
	w := watcher.NewWatcher(cfg, scanner, alerter)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("portwatch started. Press Ctrl+C to stop.")

	w.Start()

	<-sigCh
	fmt.Println("\nshutting down portwatch...")
	w.Stop()
}
