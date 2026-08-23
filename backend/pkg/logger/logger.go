package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New builds a slog.Logger that writes to stdout.
// format controls the output shape ("text" for plain lines, anything
// else falls back to JSON), and level sets the minimum severity logged.
func New(level, format string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLevel(level)}

	var handler slog.Handler
	if strings.EqualFold(format, "text") {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// parseLevel turns a level name like "debug" or "warn" into the matching
// slog level, defaulting to info for anything it doesn't recognize.
func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
