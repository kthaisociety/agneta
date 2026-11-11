package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestHealthAnswersOK checks that GET /health answers success with the body ok.
func TestHealthAnswersOK(t *testing.T) {
	recorder := httptest.NewRecorder()
	routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))
	if recorder.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if body := recorder.Body.String(); body != "ok" {
		t.Errorf("body = %q, want %q", body, "ok")
	}
}

// TestHealthRejectsOtherMethods checks that /health answers any method but GET
// with method not allowed.
func TestHealthRejectsOtherMethods(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			routes().ServeHTTP(recorder, httptest.NewRequest(method, "/health", nil))
			if recorder.Code != http.StatusMethodNotAllowed {
				t.Errorf("status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
			}
		})
	}
}

// TestUnknownPathNotFound checks that a path other than /health is not found.
func TestUnknownPathNotFound(t *testing.T) {
	for _, path := range []string{"/", "/healthz", "/health/x"} {
		t.Run(path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
			if recorder.Code != http.StatusNotFound {
				t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
			}
		})
	}
}

// freePort returns a TCP port nothing listens on right now.
func freePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("find free port: %v", err)
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

// logLines splits captured log output into its lines.
func logLines(log *bytes.Buffer) []string {
	return strings.Split(strings.TrimRight(log.String(), "\n"), "\n")
}

// getHealth calls /health on port until it answers or the deadline passes.
func getHealth(t *testing.T, port int) string {
	t.Helper()
	url := fmt.Sprintf("http://127.0.0.1:%d/health", port)
	deadline := time.Now().Add(5 * time.Second)
	for {
		response, err := http.Get(url)
		if err == nil {
			defer response.Body.Close()
			body, _ := io.ReadAll(response.Body)
			return string(body)
		}
		if time.Now().After(deadline) {
			t.Fatalf("GET %s: %v", url, err)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// TestRunListensLogsAndStops checks that Run opens the configured port, logs
// one line naming it, answers /health, and on cancel logs one line and
// returns without error.
func TestRunListensLogsAndStops(t *testing.T) {
	port := freePort(t)
	var log bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Run(ctx, Config{Port: port}, slog.New(slog.NewTextHandler(&log, nil))) }()

	if body := getHealth(t, port); body != "ok" {
		t.Errorf("GET /health body = %q, want %q", body, "ok")
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after cancel")
	}

	lines := logLines(&log)
	if len(lines) != 2 {
		t.Fatalf("log has %d lines, want 2:\n%s", len(lines), log.String())
	}
	if !strings.Contains(lines[0], fmt.Sprintf("port=%d", port)) {
		t.Errorf("start line %q does not name port %d", lines[0], port)
	}
	if !strings.Contains(lines[1], "stopped") {
		t.Errorf("shutdown line %q does not say the service stopped", lines[1])
	}
}

// TestServeFinishesInFlightRequests checks that a request in progress when
// shutdown starts still gets its full answer, and new connections are refused.
func TestServeFinishesInFlightRequests(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	slow := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		io.WriteString(w, "done")
	})
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := listener.Addr().String()
	ctx, cancel := context.WithCancel(context.Background())
	served := make(chan error, 1)
	go func() { served <- serve(ctx, listener, slow, slog.New(slog.NewTextHandler(io.Discard, nil))) }()

	answer := make(chan string, 1)
	go func() {
		response, err := http.Get("http://" + addr + "/")
		if err != nil {
			answer <- "error: " + err.Error()
			return
		}
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		answer <- string(body)
	}()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("request never reached the handler")
	}

	cancel()
	refusedBy := time.Now().Add(5 * time.Second)
	for {
		connection, err := net.Dial("tcp", addr)
		if err != nil {
			break
		}
		connection.Close()
		if time.Now().After(refusedBy) {
			t.Fatal("service still accepts new connections after shutdown began")
		}
		time.Sleep(20 * time.Millisecond)
	}
	close(release)

	if got := <-answer; got != "done" {
		t.Errorf("in-flight request got %q, want %q", got, "done")
	}
	select {
	case err := <-served:
		if err != nil {
			t.Fatalf("serve: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("serve did not return after shutdown")
	}
}
