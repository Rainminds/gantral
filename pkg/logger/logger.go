package logger

import (
	"log/slog"
	"os"
)

// Init initializes the default global logger.
// It configures logical defaults like JSON format for production.
func Init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// Use JSON handler for structured logging
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)
}
