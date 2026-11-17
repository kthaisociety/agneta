// Package server runs the agneta HTTP service.
package server

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/kthaisociety/agneta/internal/config"
)

// shutdownLimit is how long in-flight requests get to finish once the service is told to stop.
const shutdownLimit = 10 * time.Second

// readHeaderLimit is how long a client gets to send its request headers.
const readHeaderLimit = 10 * time.Second

// Routes serves GET /health, and mcp on /mcp unless it is nil; any other path is not found.
func Routes(mcp http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", answerHealth)
	if mcp != nil {
		mux.Handle("/mcp", mcp)
	}
	return mux
}

// answerHealth tells the caller the process is up.
func answerHealth(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "ok")
}

// Run serves handler on cfg.Host:cfg.Port until ctx is cancelled, logging one line on start and one on stop.
func Run(ctx context.Context, cfg config.Config, handler http.Handler, logger *slog.Logger) error {
	listener, err := net.Listen("tcp", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)))
	if err != nil {
		return err
	}
	return serve(ctx, listener, handler, logger)
}

// serve answers on listener until ctx is cancelled, then lets in-flight requests finish within shutdownLimit.
func serve(ctx context.Context, listener net.Listener, handler http.Handler, logger *slog.Logger) error {
	server := &http.Server{Handler: handler, ReadHeaderTimeout: readHeaderLimit}
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
