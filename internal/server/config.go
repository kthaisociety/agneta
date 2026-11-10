package server

// Config holds the settings the service reads from its environment.
type Config struct {
	// Port is the TCP port the service listens on.
	Port int
}

// LoadConfig reads the service settings through getenv, which looks up one
// environment variable by name and returns "" when it is unset.
func LoadConfig(getenv func(string) string) (Config, error) {
	return Config{}, nil
}
