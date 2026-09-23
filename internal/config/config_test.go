package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_ADDR", "")
	t.Setenv("APP_READ_TIMEOUT", "")
	t.Setenv("APP_WRITE_TIMEOUT", "")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if cfg.Addr != DefaultAddr {
		t.Errorf("Addr = %q, want %q", cfg.Addr, DefaultAddr)
	}
	if cfg.ReadTimeout != DefaultReadTimeout {
		t.Errorf("ReadTimeout = %s, want %s", cfg.ReadTimeout, DefaultReadTimeout)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("APP_ADDR", "127.0.0.1:9090")
	t.Setenv("APP_READ_TIMEOUT", "3s")
	t.Setenv("APP_WRITE_TIMEOUT", "")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if cfg.Addr != "127.0.0.1:9090" {
		t.Errorf("Addr = %q, want 127.0.0.1:9090", cfg.Addr)
	}
	if cfg.ReadTimeout != 3*time.Second {
		t.Errorf("ReadTimeout = %s, want 3s", cfg.ReadTimeout)
	}
}

func TestLoadInvalidDuration(t *testing.T) {
	t.Setenv("APP_ADDR", ":8080")
	t.Setenv("APP_READ_TIMEOUT", "abc")
	t.Setenv("APP_WRITE_TIMEOUT", "")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want error for invalid duration")
	}
}

func TestLoadNegativeDuration(t *testing.T) {
	t.Setenv("APP_ADDR", ":8080")
	t.Setenv("APP_READ_TIMEOUT", "-1s")
	t.Setenv("APP_WRITE_TIMEOUT", "")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want error for negative duration")
	}
}
