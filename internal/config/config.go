package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	// PollInterval is the number of seconds between port scans.
	PollInterval int `json:"poll_interval"`

	// AllowedPorts is a list of ports that are expected to be bound.
	// Bindings on these ports will be reported at Info level rather than Warn.
	AllowedPorts []int `json:"allowed_ports"`

	// LogFile is an optional path to write alerts as JSON lines.
	// If empty, alerts are written to stdout only.
	LogFile string `json:"log_file"`

	// AlertFormat controls the output format: "text" or "json".
	AlertFormat string `json:"alert_format"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		PollInterval: 5,
		AllowedPorts: []int{},
		AlertFormat:  "text",
	}
}

// LoadConfig reads a JSON config file from path and returns a Config.
// Fields not present in the file retain their default values.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	if err := dec.Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}

// validate checks that the Config fields hold acceptable values.
func (c *Config) validate() error {
	if c.PollInterval <= 0 {
		return fmt.Errorf("poll_interval must be > 0, got %d", c.PollInterval)
	}
	if c.AlertFormat != "text" && c.AlertFormat != "json" {
		return fmt.Errorf("alert_format must be \"text\" or \"json\", got %q", c.AlertFormat)
	}
	for _, p := range c.AllowedPorts {
		if p < 1 || p > 65535 {
			return fmt.Errorf("allowed_ports contains invalid port %d", p)
		}
	}
	return nil
}
