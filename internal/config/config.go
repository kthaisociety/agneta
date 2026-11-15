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
	// Host is the bind address; empty binds all interfaces, as the container needs.
	Host string
	// Port is the TCP port the service listens on.
	Port int
}

// Load reads the settings through getenv; an invalid PORT is an error naming it, never the default.
func Load(getenv func(string) string) (Config, error) {
	port, err := loadPort(getenv("PORT"))
	if err != nil {
		return Config{}, err
	}
	return Config{Host: strings.TrimSpace(getenv("HOST")), Port: port}, nil
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
