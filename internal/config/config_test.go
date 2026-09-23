package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// writeConfig writes content to a temp file and returns its path.
func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func TestLoadFullFile(t *testing.T) {
	path := writeConfig(t, `
server:
  addr: "127.0.0.1:9090"
  read_timeout: 3s
  write_timeout: 20s
  shutdown_timeout: 30s
log:
  level: debug
  format: text
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if cfg.Server.Addr != "127.0.0.1:9090" {
		t.Errorf("Addr = %q, want 127.0.0.1:9090", cfg.Server.Addr)
	}
	if cfg.Server.ReadTimeout != 3*time.Second {
		t.Errorf("ReadTimeout = %s, want 3s", cfg.Server.ReadTimeout)
	}
	if cfg.Server.WriteTimeout != 20*time.Second {
		t.Errorf("WriteTimeout = %s, want 20s", cfg.Server.WriteTimeout)
	}
	if cfg.Server.ShutdownTimeout != 30*time.Second {
		t.Errorf("ShutdownTimeout = %s, want 30s", cfg.Server.ShutdownTimeout)
	}
	if cfg.Log.Level != "debug" || cfg.Log.Format != "text" {
		t.Errorf("Log = %+v, want debug/text", cfg.Log)
	}
}

func TestLoadPartialFileKeepsDefaults(t *testing.T) {
	path := writeConfig(t, "server:\n  addr: \":9099\"\n")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if cfg.Server.Addr != ":9099" {
		t.Errorf("Addr = %q, want :9099", cfg.Server.Addr)
	}
	if cfg.Server.WriteTimeout != DefaultWriteTimeout {
		t.Errorf("WriteTimeout = %s, want default %s", cfg.Server.WriteTimeout, DefaultWriteTimeout)
	}
	if cfg.Log.Level != DefaultLogLevel {
		t.Errorf("Log.Level = %q, want default %q", cfg.Log.Level, DefaultLogLevel)
	}
}

func TestLoadEmptyFileUsesDefaults(t *testing.T) {
	path := writeConfig(t, "")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v, want nil for empty file", err)
	}
	if cfg.Server.Addr != DefaultAddr {
		t.Errorf("Addr = %q, want default %q", cfg.Server.Addr, DefaultAddr)
	}
}

func TestLoadMissingFileIsError(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "absent.yaml")); err == nil {
		t.Fatal("Load() error = nil, want not-found error")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	if _, err := Load(writeConfig(t, "server: [unclosed\n")); err == nil {
		t.Fatal("Load() error = nil, want parse error")
	}
}

func TestLoadUnknownFieldRejected(t *testing.T) {
	if _, err := Load(writeConfig(t, "server:\n  addrr: \":8080\"\n")); err == nil {
		t.Fatal("Load() error = nil, want error for unknown field")
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"empty addr", "server:\n  addr: \"\"\n"},
		{"blank addr", "server:\n  addr: \"   \"\n"},
		{"zero timeout", "server:\n  read_timeout: 0s\n"},
		{"negative timeout", "server:\n  shutdown_timeout: -5s\n"},
		{"bad level", "log:\n  level: verbose\n"},
		{"bad format", "log:\n  format: xml\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Load(writeConfig(t, tt.content)); err == nil {
				t.Fatal("Load() error = nil, want error")
			}
		})
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		in   string
		want int
		fail bool
	}{
		{"debug", -4, false},
		{"INFO", 0, false},
		{"warn", 4, false},
		{"error", 8, false},
		{"verbose", 0, true},
	}

	for _, tt := range tests {
		got, err := ParseLogLevel(tt.in)
		if tt.fail {
			if err == nil {
				t.Errorf("ParseLogLevel(%q) error = nil, want error", tt.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseLogLevel(%q) error = %v, want nil", tt.in, err)
			continue
		}
		if int(got) != tt.want {
			t.Errorf("ParseLogLevel(%q) = %d, want %d", tt.in, int(got), tt.want)
		}
	}
}
