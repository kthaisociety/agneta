package config

import (
	"strings"
	"testing"
)

// env returns a getenv that answers from vars and returns "" for anything else.
func env(vars map[string]string) func(string) string {
	return func(name string) string { return vars[name] }
}

// TestLoadReadsSettings checks that HOST and PORT come from the environment.
func TestLoadReadsSettings(t *testing.T) {
	cfg, err := Load(env(map[string]string{"HOST": "127.0.0.1", "PORT": "9123"}))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Host != "127.0.0.1" || cfg.Port != 9123 {
		t.Errorf("cfg = %+v, want host 127.0.0.1 and port 9123", cfg)
	}
}

// TestLoadDefaultsPort checks that an unset PORT gives port 8000.
func TestLoadDefaultsPort(t *testing.T) {
	cfg, err := Load(env(nil))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Port != 8000 {
		t.Errorf("Port = %d, want 8000", cfg.Port)
	}
}

// TestLoadRejectsInvalidPort checks that a bad PORT is an error naming the setting.
func TestLoadRejectsInvalidPort(t *testing.T) {
	for _, value := range []string{"abc", "80a", "0", "-1", "65536", " "} {
		t.Run(value, func(t *testing.T) {
			_, err := Load(env(map[string]string{"PORT": value}))
			if err == nil {
				t.Fatalf("Load accepted PORT=%q", value)
			}
			if !strings.Contains(err.Error(), "PORT") {
				t.Errorf("error %q does not name PORT", err)
			}
		})
	}
}
