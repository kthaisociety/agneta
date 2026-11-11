// Package server runs the agneta HTTP service.
package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// shutdownLimit is how long requests in progress get to finish once the
// service is told to stop.
const shutdownLimit = 10 * time.Second

// routes returns the routes the service answers: GET /health answers ok.
// Any other method on /health is method not allowed, and any other path is
// not found.
func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", answerHealth)
	return mux
}

// answerHealth tells the caller the process is up.
func answerHealth(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "ok")
}

// Run listens on cfg.Port and serves routes until ctx is cancelled. It logs
// one line naming the port once listening and one line once stopped. It
// returns an error if the port cannot be opened or the server fails; a
// shutdown through ctx returns nil.
func Run(ctx context.Context, cfg Config, logger *slog.Logger) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return err
	}
	return serve(ctx, listener, routes(), logger)
}

// serve answers requests on listener with handler until ctx is cancelled.
// It then stops accepting connections, lets requests in progress finish
// within shutdownLimit, and returns. It closes listener.
func serve(ctx context.Context, listener net.Listener, handler http.Handler, logger *slog.Logger) error {
	server := &http.Server{Handler: handler}
	failed := make(chan error, 1)
	go func() { failed <- server.Serve(listener) }()
	logger.Info("listening", "port", listener.Addr().(*net.TCPAddr).Port)

	select {
	case err := <-failed:
		return err
	case <-ctx.Done():
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownLimit)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	if err := <-failed; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	logger.Info("stopped")
	return nil
}
