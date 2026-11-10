package server

import (
	"fmt"
	"strconv"
)

// portVariable is the environment variable that holds the listening port.
const portVariable = "PORT"

// defaultPort is the port used when portVariable is unset. It matches the
// port the container image exposes.
const defaultPort = 8000

// highestPort is the largest valid TCP port number.
const highestPort = 65535

// Config holds the settings the service reads from its environment.
type Config struct {
	// Port is the TCP port the service listens on.
	Port int
}

// LoadConfig reads the service settings through getenv, which looks up one
// environment variable by name and returns "" when it is unset. An unset PORT
// gives port 8000. A PORT that is set but not a whole number from 1 to 65535
// is an error naming the variable; it never falls back to the default.
func LoadConfig(getenv func(string) string) (Config, error) {
	raw := getenv(portVariable)
	if raw == "" {
		return Config{Port: defaultPort}, nil
	}
	port, err := strconv.Atoi(raw)
	if err != nil {
		return Config{}, fmt.Errorf("%s %q is not a number", portVariable, raw)
	}
	if port < 1 || port > highestPort {
		return Config{}, fmt.Errorf("%s %d is out of range 1-%d", portVariable, port, highestPort)
	}
	return Config{Port: port}, nil
}
