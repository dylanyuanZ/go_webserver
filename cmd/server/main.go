// Package main is the entrypoint of the go_webserver service.
package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/your-org/go_webserver/internal/config"
	"github.com/your-org/go_webserver/internal/server"
)

func main() {
	configPath := flag.String("config", config.DefaultFile, "path to the YAML config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("load config failed", "error", err, "path", *configPath)
		os.Exit(1)
	}

	logger := newLogger(cfg)
	slog.SetDefault(logger)

	srv := server.New(cfg, logger)

	errCh := make(chan error, 1)
	go func() {
		logger.Info("server starting", "addr", cfg.Server.Addr, "config", *configPath)
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil {
			logger.Error("server exited with error", "error", err)
			os.Exit(1)
		}
	case sig := <-quit:
		logger.Info("shutdown signal received", "signal", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("graceful shutdown failed", "error", err)
			os.Exit(1)
		}
	}

	logger.Info("server stopped")
}

// newLogger builds a slog.Logger from the logging configuration.
func newLogger(cfg *config.Config) *slog.Logger {
	level, err := config.ParseLogLevel(cfg.Log.Level)
	if err != nil {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}
	if cfg.Log.Format == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}
