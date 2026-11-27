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

// startupFails runs the service with PORT and returns its exit code and stderr, failing if it keeps running.
func startupFails(t *testing.T, port string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var stderr bytes.Buffer
	getenv := func(name string) string {
		switch name {
		case "PORT":
			return port
		case "AUTH_DISABLED":
			return "true" // these tests are about the port, not sign-in
		}
		return ""
	}
	code := run(ctx, getenv, &stderr)
	if ctx.Err() != nil {
		t.Fatalf("service with PORT=%q kept running instead of failing", port)
	}
	return code, stderr.String()
}

// assertOneLineFailure fails unless the exit code is a failure and stderr is one line containing want.
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

// TestInvalidPortStopsStartup checks that a bad PORT stops startup with one line naming it.
func TestInvalidPortStopsStartup(t *testing.T) {
	for _, port := range []string{"abc", "70000"} {
		t.Run(port, func(t *testing.T) {
			code, stderr := startupFails(t, port)
			assertOneLineFailure(t, code, stderr, "PORT")
		})
	}
}

// TestPortInUseStopsStartup checks that a taken port stops startup with one line saying so.
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
