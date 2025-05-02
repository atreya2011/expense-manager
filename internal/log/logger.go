package log

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

// NewLogger creates a new structured logger
func NewLogger() *slog.Logger {
	// Create a JSON handler with default options
	handler := slog.NewJSONHandler(os.Stdout, nil)

	// Create a new logger with the handler
	logger := slog.New(handler)

	return logger
}

// LogWithContext logs with trace context information automatically extracted
func LogWithContext(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, args ...any) {
	// Extract trace IDs if present in context
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		// Add trace information to args
		args = append(args,
			"trace_id", spanCtx.TraceID().String(),
			"span_id", spanCtx.SpanID().String(),
		)
	}

	// Log with the context information
	logger.Log(ctx, level, msg, args...)
}

// InfoContext logs at info level with trace context
func InfoContext(ctx context.Context, logger *slog.Logger, msg string, args ...any) {
	LogWithContext(ctx, logger, slog.LevelInfo, msg, args...)
}

// ErrorContext logs at error level with trace context
func ErrorContext(ctx context.Context, logger *slog.Logger, msg string, args ...any) {
	LogWithContext(ctx, logger, slog.LevelError, msg, args...)
}

// WarnContext logs at warn level with trace context
func WarnContext(ctx context.Context, logger *slog.Logger, msg string, args ...any) {
	LogWithContext(ctx, logger, slog.LevelWarn, msg, args...)
}

// DebugContext logs at debug level with trace context
func DebugContext(ctx context.Context, logger *slog.Logger, msg string, args ...any) {
	LogWithContext(ctx, logger, slog.LevelDebug, msg, args...)
}
