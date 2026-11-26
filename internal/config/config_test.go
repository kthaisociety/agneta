package config

import (
	"maps"
	"strings"
	"testing"
)

// env returns a getenv that answers from vars and returns "" for anything else.
func env(vars map[string]string) func(string) string {
	return func(name string) string { return vars[name] }
}

// authVars returns a fresh, complete set of sign-in settings; tests change or delete entries in it.
func authVars() map[string]string {
	return map[string]string{
		"PUBLIC_URL":             "https://agneta.example.com",
		"IDENTITY_ISSUER":        "https://identity.example.com",
		"IDENTITY_INTERNAL_URL":  "http://identity.internal:8080",
		"IDENTITY_CLIENT_ID":     "agneta",
		"IDENTITY_CLIENT_SECRET": "test-secret",
	}
}

// with returns authVars plus extra, which overrides it.
func with(extra map[string]string) map[string]string {
	vars := authVars()
	maps.Copy(vars, extra)
	return vars
}

// TestLoadReadsSettings checks that HOST and PORT come from the environment.
func TestLoadReadsSettings(t *testing.T) {
	cfg, err := Load(env(with(map[string]string{"HOST": "127.0.0.1", "PORT": "9123"})))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Host != "127.0.0.1" || cfg.Port != 9123 {
		t.Errorf("cfg = %+v, want host 127.0.0.1 and port 9123", cfg)
	}
}

// TestLoadDefaultsPort checks that an unset PORT gives port 8000.
func TestLoadDefaultsPort(t *testing.T) {
	cfg, err := Load(env(authVars()))
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
			_, err := Load(env(with(map[string]string{"PORT": value})))
			if err == nil {
				t.Fatalf("Load accepted PORT=%q", value)
			}
			if !strings.Contains(err.Error(), "PORT") {
				t.Errorf("error %q does not name PORT", err)
			}
		})
	}
}

// TestLoadReadsAuthSettings checks that every sign-in setting lands in its field.
func TestLoadReadsAuthSettings(t *testing.T) {
	cfg, err := Load(env(authVars()))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := Config{
		Port:                 8000,
		PublicURL:            "https://agneta.example.com",
		IdentityIssuer:       "https://identity.example.com",
		IdentityInternalURL:  "http://identity.internal:8080",
		IdentityClientID:     "agneta",
		IdentityClientSecret: "test-secret",
	}
	if cfg != want {
		t.Errorf("cfg =\n %+v\nwant\n %+v", cfg, want)
	}
}

// TestLoadRequiresAuthSettings checks that each missing or blank sign-in setting is an error naming it.
func TestLoadRequiresAuthSettings(t *testing.T) {
	for name := range authVars() {
		for _, value := range []string{"", "   "} {
			t.Run(name+"="+value, func(t *testing.T) {
				_, err := Load(env(with(map[string]string{name: value})))
				if err == nil {
					t.Fatalf("Load accepted %s=%q", name, value)
				}
				if !strings.Contains(err.Error(), name) {
					t.Errorf("error %q does not name %s", err, name)
				}
			})
		}
	}
}

// TestLoadAuthDisabled checks that AUTH_DISABLED=true needs no sign-in settings.
func TestLoadAuthDisabled(t *testing.T) {
	cfg, err := Load(env(map[string]string{"AUTH_DISABLED": "true"}))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.AuthDisabled {
		t.Error("AuthDisabled = false, want true")
	}
}

// TestLoadRejectsInvalidAuthDisabled checks that a value other than true/false is an error naming it.
func TestLoadRejectsInvalidAuthDisabled(t *testing.T) {
	for _, value := range []string{"yes", "on", "tru"} {
		t.Run(value, func(t *testing.T) {
			_, err := Load(env(with(map[string]string{"AUTH_DISABLED": value})))
			if err == nil {
				t.Fatalf("Load accepted AUTH_DISABLED=%q", value)
			}
			if !strings.Contains(err.Error(), "AUTH_DISABLED") {
				t.Errorf("error %q does not name AUTH_DISABLED", err)
			}
		})
	}
}
