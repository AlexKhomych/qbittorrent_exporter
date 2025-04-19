package log

import (
	"log/slog"
	"os"
	"strings"
)

var (
	logger *slog.Logger
	lvl    *slog.LevelVar = new(slog.LevelVar)
)

// Debug calls [slog.Debug].
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

// Info calls [slog.Info].
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

// Warn calls [slog.Warn].
func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

// Error calls [slog.Error].
func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

// Fatal calls [slog.Error] and [os.Exit(1)].
func Fatal(msg string, args ...any) {
	logger.Error(msg, args...)
	os.Exit(1)
}

func Set(logLevel, logFormat string) *slog.Logger {
	logLevel = strings.ToLower(logLevel)
	switch logLevel {
	case "debug":
		lvl.Set(slog.LevelDebug)
	case "info":
		lvl.Set(slog.LevelInfo)
	case "warn":
		lvl.Set(slog.LevelWarn)
	case "error":
		lvl.Set(slog.LevelError)
	default:
		lvl.Set(slog.LevelInfo)
	}

	opts := &slog.HandlerOptions{
		Level: lvl,
	}

	switch logFormat {
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	default:
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	}

	slog.SetDefault(logger)
	Info("Updated logger with LogFormat: " + logFormat + ", LogLevel: " + logLevel)
	return logger
}
