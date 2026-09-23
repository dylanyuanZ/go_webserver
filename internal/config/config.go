// Package config loads service configuration from a YAML file.
package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// DefaultFile is the config file path used when no path is given.
const DefaultFile = "configs/config.yaml"

// Default configuration values used when the file omits a field.
const (
	DefaultAddr            = ":8080"
	DefaultReadTimeout     = 5 * time.Second
	DefaultWriteTimeout    = 10 * time.Second
	DefaultShutdownTimeout = 10 * time.Second
	DefaultLogLevel        = "info"
	DefaultLogFormat       = "json"
)

// Config is the fully resolved runtime configuration.
type Config struct {
	// Server holds HTTP server settings.
	Server ServerConfig
	// Log holds logging settings.
	Log LogConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	// Addr is the listen address, e.g. ":8080".
	Addr string
	// ReadTimeout bounds the time spent reading a request.
	ReadTimeout time.Duration
	// WriteTimeout bounds the time spent writing a response.
	WriteTimeout time.Duration
	// ShutdownTimeout bounds the graceful shutdown grace period.
	ShutdownTimeout time.Duration
}

// LogConfig holds logging settings.
type LogConfig struct {
	// Level is one of debug, info, warn, error.
	Level string
	// Format is one of json, text.
	Format string
}

// fileConfig mirrors Config with pointers so that "omitted" is distinguishable from "zero".
type fileConfig struct {
	Server struct {
		Addr            *string        `yaml:"addr"`
		ReadTimeout     *time.Duration `yaml:"read_timeout"`
		WriteTimeout    *time.Duration `yaml:"write_timeout"`
		ShutdownTimeout *time.Duration `yaml:"shutdown_timeout"`
	} `yaml:"server"`
	Log struct {
		Level  *string `yaml:"level"`
		Format *string `yaml:"format"`
	} `yaml:"log"`
}

// Load reads the config file at path and overlays it on top of the defaults.
// A missing or unreadable file is an error: the file is the only configuration source.
// An empty file is accepted and yields the defaults.
func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultFile
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("config: file %s not found: %w", path, err)
		}
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}

	var fc fileConfig
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&fc); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("config: parse %s: %w", path, err)
	}

	cfg := defaults()
	cfg.applyFile(&fc)
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// defaults returns the built-in configuration.
func defaults() *Config {
	return &Config{
		Server: ServerConfig{
			Addr:            DefaultAddr,
			ReadTimeout:     DefaultReadTimeout,
			WriteTimeout:    DefaultWriteTimeout,
			ShutdownTimeout: DefaultShutdownTimeout,
		},
		Log: LogConfig{
			Level:  DefaultLogLevel,
			Format: DefaultLogFormat,
		},
	}
}

// applyFile overlays fields present in the file on top of the current values.
func (c *Config) applyFile(fc *fileConfig) {
	if fc.Server.Addr != nil {
		c.Server.Addr = *fc.Server.Addr
	}
	if fc.Server.ReadTimeout != nil {
		c.Server.ReadTimeout = *fc.Server.ReadTimeout
	}
	if fc.Server.WriteTimeout != nil {
		c.Server.WriteTimeout = *fc.Server.WriteTimeout
	}
	if fc.Server.ShutdownTimeout != nil {
		c.Server.ShutdownTimeout = *fc.Server.ShutdownTimeout
	}
	if fc.Log.Level != nil {
		c.Log.Level = *fc.Log.Level
	}
	if fc.Log.Format != nil {
		c.Log.Format = *fc.Log.Format
	}
}

// validate rejects invalid values.
func (c *Config) validate() error {
	if strings.TrimSpace(c.Server.Addr) == "" {
		return errors.New("config: server.addr must not be empty")
	}
	timeouts := map[string]time.Duration{
		"server.read_timeout":     c.Server.ReadTimeout,
		"server.write_timeout":    c.Server.WriteTimeout,
		"server.shutdown_timeout": c.Server.ShutdownTimeout,
	}
	for name, d := range timeouts {
		if d <= 0 {
			return fmt.Errorf("config: %s must be positive, got %s", name, d)
		}
	}
	if _, err := ParseLogLevel(c.Log.Level); err != nil {
		return err
	}
	if c.Log.Format != "json" && c.Log.Format != "text" {
		return fmt.Errorf("config: log.format must be json or text, got %q", c.Log.Format)
	}
	return nil
}

// ParseLogLevel converts a level name into slog.Level.
func ParseLogLevel(level string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("config: unknown log level %q", level)
	}
}
