package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
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

// MaskSensitive masks sensitive values in logging fields
// It looks for common patterns like "key", "token", "secret", "password"
// and masks their values
func MaskSensitive(fields ...any) []any {
	if len(fields) == 0 {
		return fields
	}

	masked := make([]any, len(fields))
	copy(masked, fields)

	for i := 0; i < len(masked)-1; i += 2 {
		key, ok := masked[i].(string)
		if !ok {
			continue
		}

		// Check if the field name suggests it's sensitive
		lowerKey := strings.ToLower(key)
		if isSensitiveField(lowerKey) {
			if strVal, ok := masked[i+1].(string); ok {
				masked[i+1] = maskValue(strVal)
			}
		}
	}

	return masked
}

// isSensitiveField checks if a field name suggests it contains sensitive data
func isSensitiveField(name string) bool {
	sensitivePatterns := []string{
		"key",
		"token",
		"secret",
		"password",
		"credential",
		"auth",
		"api_key",
		"apikey",
	}

	for _, pattern := range sensitivePatterns {
		if strings.Contains(name, pattern) {
			return true
		}
	}

	return false
}

// maskValue masks a string value for logging
// Shows only first 2 and last 2 characters for security
func maskValue(value string) string {
	if value == "" {
		return "<empty>"
	}

	length := len(value)
	if length <= 4 {
		return "****"
	}

	return fmt.Sprintf("%s...%s", value[:2], value[length-2:])
}

// InfoWithMasking logs an info message with automatic masking of sensitive fields
func InfoWithMasking(ctx context.Context, msg string, fields ...any) {
	Logger.InfoContext(ctx, msg, MaskSensitive(fields...)...)
}

// DebugWithMasking logs a debug message with automatic masking of sensitive fields
func DebugWithMasking(ctx context.Context, msg string, fields ...any) {
	Logger.DebugContext(ctx, msg, MaskSensitive(fields...)...)
}

// WarnWithMasking logs a warning message with automatic masking of sensitive fields
func WarnWithMasking(ctx context.Context, msg string, fields ...any) {
	Logger.WarnContext(ctx, msg, MaskSensitive(fields...)...)
}

// ErrorWithMasking logs an error message with automatic masking of sensitive fields
func ErrorWithMasking(ctx context.Context, msg string, fields ...any) {
	Logger.ErrorContext(ctx, msg, MaskSensitive(fields...)...)
}
