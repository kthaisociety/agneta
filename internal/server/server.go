// Package server runs the agneta HTTP service.
package server

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
)

// Handler returns the routes the service answers: GET /health answers ok.
// Any other method on /health is method not allowed, and any other path is
// not found.
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", answerHealth)
	return mux
}

// answerHealth tells the caller the process is up.
func answerHealth(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "ok")
}

// Run listens on cfg.Port and serves Handler until ctx is cancelled.
func Run(ctx context.Context, cfg Config, logger *slog.Logger) error {
	return nil
}

func serve(ctx context.Context, listener net.Listener, handler http.Handler, logger *slog.Logger) error {
	return nil
}
