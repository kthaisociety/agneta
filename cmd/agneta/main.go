// Command agneta runs the agneta MCP server until it receives SIGINT or SIGTERM.
package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/kthaisociety/agneta/internal/config"
	"github.com/kthaisociety/agneta/internal/mcpserver"
	"github.com/kthaisociety/agneta/internal/server"
)

func main() {
	_ = godotenv.Load() // Best effort: a container has no .env.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	os.Exit(run(ctx, os.Getenv, os.Stderr))
}

// run starts the service and returns its exit code, logging one line saying why if startup or serving fails.
func run(ctx context.Context, getenv func(string) string, stderr io.Writer) int {
	logger := slog.New(slog.NewTextHandler(stderr, nil))
	if err := start(ctx, getenv, logger); err != nil {
		logger.Error("cannot run", "err", err)
		return 1
	}
	return 0
}

// start reads the settings and serves until ctx is cancelled.
func start(ctx context.Context, getenv func(string) string, logger *slog.Logger) error {
	cfg, err := config.Load(getenv)
	if err != nil {
		return err
	}
	return server.Run(ctx, cfg, handler(logger), logger)
}

// handler wires every route the service answers.
func handler(logger *slog.Logger) http.Handler {
	return server.Routes(mcpserver.Handler(logger))
}
