// Command agneta runs the agneta upload service. It reads its settings from
// the environment and runs until it receives SIGINT or SIGTERM.
package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kthaisociety/agneta/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	os.Exit(run(ctx, os.Getenv, os.Stderr))
}

// run starts the service with settings read through getenv and logs to
// stderr. It returns once ctx is cancelled, with exit code 0, or once startup
// or serving fails, with exit code 1 after logging one line saying why.
func run(ctx context.Context, getenv func(string) string, stderr io.Writer) int {
	logger := slog.New(slog.NewTextHandler(stderr, nil))
	if err := start(ctx, getenv, logger); err != nil {
		logger.Error("cannot run", "err", err)
		return 1
	}
	return 0
}

// start reads the settings through getenv and runs the service until ctx is
// cancelled. It returns the error that stopped startup or serving, or nil.
func start(ctx context.Context, getenv func(string) string, logger *slog.Logger) error {
	cfg, err := server.LoadConfig(getenv)
	if err != nil {
		return err
	}
	return server.Run(ctx, cfg, logger)
}
