// Package config reads every setting from the environment, once, at startup (docs/pcr/0005).
package config

import (
	"fmt"
	"strconv"
	"strings"
)

// defaultPort matches the port the container image exposes.
const defaultPort = 8000

// Config holds the settings the service reads from its environment.
type Config struct {
	// Host is the bind address, empty binds all interfaces, as the container needs.
	Host string
	// TCP port the service listens on.
	Port int

	// AuthDisabled turns sign-in off, for local work only. When true, the settings below are not read.
	AuthDisabled bool
	// PublicURL is where Claude reaches agneta.
	PublicURL string
	// IdentityIssuer is identity-service's public URL. The Advisor's browser is sent to its /authorize.
	IdentityIssuer string
	// IdentityInternalURL is identity-service on the internal network, for token exchange and keys.
	IdentityInternalURL string
	// IdentityClientID and IdentityClientSecret are agneta's credentials as an identity-service client.
	IdentityClientID     string
	IdentityClientSecret string
}

// Load reads the settings through getenv; an invalid PORT is an error naming it, never the default.
func Load(getenv func(string) string) (Config, error) {
	port, err := loadPort(getenv("PORT"))
	if err != nil {
		return Config{}, err
	}
	cfg := Config{Host: strings.TrimSpace(getenv("HOST")), Port: port}

	if raw := getenv("AUTH_DISABLED"); raw != "" {
		disabled, err := strconv.ParseBool(raw)
		if err != nil {
			return Config{}, fmt.Errorf("AUTH_DISABLED %q is not true or false", raw)
		}
		cfg.AuthDisabled = disabled
	}
	if cfg.AuthDisabled {
		return cfg, nil // local mode, skip everything below
	}

	// Each required setting and the field it fills
	settings := []struct {
		name  string
		field *string
	}{
		{"PUBLIC_URL", &cfg.PublicURL},
		{"IDENTITY_ISSUER", &cfg.IdentityIssuer},
		{"IDENTITY_INTERNAL_URL", &cfg.IdentityInternalURL},
		{"IDENTITY_CLIENT_ID", &cfg.IdentityClientID},
		{"IDENTITY_CLIENT_SECRET", &cfg.IdentityClientSecret},
	}
	for _, setting := range settings {
		if *setting.field, err = required(getenv, setting.name); err != nil {
			return Config{}, err
		}
	}
	return cfg, nil
}

// required returns getenv(name) trimmed, or an error naming the setting if it is empty.
func required(getenv func(string) string, name string) (string, error) {
	value := strings.TrimSpace(getenv(name))
	if value == "" {
		return "", fmt.Errorf("%s is not set", name)
	}
	return value, nil
}

// loadPort parses raw as a port from 1 to 65535, or gives the default for "".
func loadPort(raw string) (int, error) {
	if raw == "" {
		return defaultPort, nil
	}
	port, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("PORT %q is not a number", raw)
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("PORT %d is out of range 1-65535", port)
	}
	return port, nil
}
