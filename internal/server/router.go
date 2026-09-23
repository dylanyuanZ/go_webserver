package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// registerRoutes attaches all HTTP routes to mux.
func registerRoutes(mux *http.ServeMux, logger *slog.Logger) {
	mux.HandleFunc("GET /healthz", handleHealthz(logger))
	mux.HandleFunc("GET /", handleIndex(logger))
}

// handleHealthz reports service liveness.
func handleHealthz(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		logger.Debug("healthz checked", "remote", r.RemoteAddr)
	}
}

// handleIndex returns basic service information.
func handleIndex(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"service": "go_webserver",
			"docs":    "https://github.com/your-org/go_webserver",
		})
		logger.Debug("index served", "remote", r.RemoteAddr)
	}
}

// writeJSON serializes v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("write response failed", "error", err)
	}
}
