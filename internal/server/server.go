// Package server wires the HTTP server, routes and graceful shutdown.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/your-org/go_webserver/internal/config"
)

// Server wraps http.Server with configuration and logging.
type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	http   *http.Server
}

// New builds a Server with the given config and logger. It never returns nil.
func New(cfg *config.Config, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	mux := http.NewServeMux()
	registerRoutes(mux, logger)

	return &Server{
		cfg:    cfg,
		logger: logger,
		http: &http.Server{
			Addr:         cfg.Addr,
			Handler:      mux,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

// Start listens and serves until the server is shut down.
// It returns http.ErrServerClosed on graceful shutdown.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return fmt.Errorf("server: listen %s: %w", s.cfg.Addr, err)
	}
	if err := s.http.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server: serve: %w", err)
	}
	return http.ErrServerClosed
}

// Shutdown gracefully stops the server, waiting at most cfg.ShutdownTimeout.
func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("server: shutdown: %w", err)
	}
	s.logger.Info("server shutdown completed", "grace", s.cfg.ShutdownTimeout.String())
	return nil
}

// Addr returns the configured listen address.
func (s *Server) Addr() string {
	return s.cfg.Addr
}
