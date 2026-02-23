package utils

import (
	"log/slog"
	"os"
	"strings"
)

var Logger *slog.Logger

func init() {
	level := slog.LevelInfo
	if s := os.Getenv("LOG_LEVEL"); s != "" {
		switch strings.ToLower(s) {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn", "warning":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
	}

	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}
