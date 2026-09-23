// Package config loads service configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"time"
)

// Default configuration values used when the environment variable is unset.
const (
	DefaultAddr            = ":8080"
	DefaultReadTimeout     = 5 * time.Second
	DefaultWriteTimeout    = 10 * time.Second
	DefaultShutdownTimeout = 10 * time.Second
)

// Config holds all runtime configuration of the service.
type Config struct {
	// Addr is the listen address, e.g. ":8080".
	Addr string
	// ReadTimeout bounds the time spent reading a request.
	ReadTimeout time.Duration
	// WriteTimeout bounds the time spent writing a response.
	WriteTimeout time.Duration
	// ShutdownTimeout bounds the graceful shutdown grace period.
	ShutdownTimeout time.Duration
}

// Load reads configuration from environment variables and validates it.
// Fields absent from the environment fall back to Default* constants.
func Load() (*Config, error) {
	cfg := &Config{
		Addr:            envString("APP_ADDR", DefaultAddr),
		ReadTimeout:     DefaultReadTimeout,
		WriteTimeout:    DefaultWriteTimeout,
		ShutdownTimeout: DefaultShutdownTimeout,
	}

	for name, target := range map[string]*time.Duration{
		"APP_READ_TIMEOUT":     &cfg.ReadTimeout,
		"APP_WRITE_TIMEOUT":    &cfg.WriteTimeout,
		"APP_SHUTDOWN_TIMEOUT": &cfg.ShutdownTimeout,
	} {
		raw := os.Getenv(name)
		if raw == "" {
			continue
		}
		d, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("config: invalid %s %q: %w", name, raw, err)
		}
		if d <= 0 {
			return nil, fmt.Errorf("config: %s must be positive, got %s", name, d)
		}
		*target = d
	}

	if cfg.Addr == "" {
		return nil, fmt.Errorf("config: APP_ADDR must not be empty")
	}
	return cfg, nil
}

// envString returns the value of key, or fallback when key is unset or empty.
func envString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
