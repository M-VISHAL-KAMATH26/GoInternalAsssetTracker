package logger

import (
	"log/slog"
	"os"
)

// New returns a structured JSON logger tagged with the given service
// name, so log lines from different services are distinguishable when
// aggregated (e.g. piped into a file or log collector later).
func New(serviceName string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(handler).With("service", serviceName)
}