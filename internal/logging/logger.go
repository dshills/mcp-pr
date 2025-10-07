package logging

import (
	"context"
	"log/slog"
	"os"
)

// Logger is the structured logger instance
var Logger *slog.Logger

// Init initializes the structured logger with JSON format
func Init(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	Logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

// WithFields returns a logger with additional fields
func WithFields(fields ...any) *slog.Logger {
	return Logger.With(fields...)
}

// Info logs an info message
func Info(ctx context.Context, msg string, fields ...any) {
	Logger.InfoContext(ctx, msg, fields...)
}

// Debug logs a debug message
func Debug(ctx context.Context, msg string, fields ...any) {
	Logger.DebugContext(ctx, msg, fields...)
}

// Warn logs a warning message
func Warn(ctx context.Context, msg string, fields ...any) {
	Logger.WarnContext(ctx, msg, fields...)
}

// Error logs an error message
func Error(ctx context.Context, msg string, fields ...any) {
	Logger.ErrorContext(ctx, msg, fields...)
}
