package server

import (
	"strings"
	"testing"
)

// env returns a getenv that answers from vars and returns "" for anything else.
func env(vars map[string]string) func(string) string {
	return func(name string) string { return vars[name] }
}

// TestLoadConfigReadsPort checks that the port comes from the PORT variable.
func TestLoadConfigReadsPort(t *testing.T) {
	cfg, err := LoadConfig(env(map[string]string{"PORT": "9123"}))
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.Port != 9123 {
		t.Errorf("Port = %d, want 9123", cfg.Port)
	}
}

// TestLoadConfigDefaultsPort checks that an unset PORT gives port 8000.
func TestLoadConfigDefaultsPort(t *testing.T) {
	cfg, err := LoadConfig(env(nil))
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.Port != 8000 {
		t.Errorf("Port = %d, want 8000", cfg.Port)
	}
}

// TestLoadConfigRejectsInvalidPort checks that a PORT that is not a number,
// or is out of range, is an error naming the setting, never the default.
func TestLoadConfigRejectsInvalidPort(t *testing.T) {
	for _, value := range []string{"abc", "80a", "0", "-1", "65536", " "} {
		t.Run(value, func(t *testing.T) {
			_, err := LoadConfig(env(map[string]string{"PORT": value}))
			if err == nil {
				t.Fatalf("LoadConfig accepted PORT=%q", value)
			}
			if !strings.Contains(err.Error(), "PORT") {
				t.Errorf("error %q does not name PORT", err)
			}
		})
	}
}
