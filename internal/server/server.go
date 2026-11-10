// Package server runs the agneta HTTP service.
package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
)

// Handler returns the routes the service answers.
func Handler() http.Handler {
	return http.NewServeMux()
}

// Run listens on cfg.Port and serves Handler until ctx is cancelled.
func Run(ctx context.Context, cfg Config, logger *slog.Logger) error {
	return nil
}

func serve(ctx context.Context, listener net.Listener, handler http.Handler, logger *slog.Logger) error {
	return nil
}
