package main

import (
	"bytes"
	"context"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

// startupFails runs the service with the given PORT and returns its exit code
// and what it wrote to stderr. It fails the test if the service keeps running.
func startupFails(t *testing.T, port string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var stderr bytes.Buffer
	getenv := func(name string) string {
		if name == "PORT" {
			return port
		}
		return ""
	}
	code := run(ctx, getenv, &stderr)
	if ctx.Err() != nil {
		t.Fatalf("service with PORT=%q kept running instead of failing", port)
	}
	return code, stderr.String()
}

// assertOneLineFailure fails the test unless the exit code is a failure and
// stderr is exactly one line containing want.
func assertOneLineFailure(t *testing.T, code int, stderr, want string) {
	t.Helper()
	if code == 0 {
		t.Errorf("exit code = 0, want a failure")
	}
	lines := strings.Split(strings.TrimRight(stderr, "\n"), "\n")
	if len(lines) != 1 || lines[0] == "" {
		t.Fatalf("stderr has %d lines, want 1:\n%s", len(lines), stderr)
	}
	if !strings.Contains(lines[0], want) {
		t.Errorf("stderr line %q does not contain %q", lines[0], want)
	}
}

// TestInvalidPortStopsStartup checks that a PORT that is not a number or is
// out of range stops startup with one line naming PORT and a failure code.
func TestInvalidPortStopsStartup(t *testing.T) {
	for _, port := range []string{"abc", "70000"} {
		t.Run(port, func(t *testing.T) {
			code, stderr := startupFails(t, port)
			assertOneLineFailure(t, code, stderr, "PORT")
		})
	}
}

// TestPortInUseStopsStartup checks that a port something else holds stops
// startup with one line saying so and a failure code.
func TestPortInUseStopsStartup(t *testing.T) {
	holder, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("hold a port: %v", err)
	}
	defer holder.Close()
	port := strconv.Itoa(holder.Addr().(*net.TCPAddr).Port)

	code, stderr := startupFails(t, port)
	assertOneLineFailure(t, code, stderr, "address already in use")
}
