package logger

import (
	"log/slog"
	"os"
)

// Level represents log level (for backward compatibility)
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// Logger wraps slog.Logger for structured logging
type Logger struct {
	slog *slog.Logger
}

// New creates a new logger instance with slog
func New(level Level) *Logger {
	// Convert our Level to slog.Level
	var slogLevel slog.Level
	switch level {
	case DEBUG:
		slogLevel = slog.LevelDebug
	case INFO:
		slogLevel = slog.LevelInfo
	case WARN:
		slogLevel = slog.LevelWarn
	case ERROR:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	// Create handler with JSON output for structured logging
	opts := &slog.HandlerOptions{
		Level: slogLevel,
	}
	
	// Use JSON handler for production, text handler for development
	handler := slog.NewJSONHandler(os.Stdout, opts)
	
	return &Logger{
		slog: slog.New(handler),
	}
}

// NewText creates a logger with human-readable text output (for development)
func NewText(level Level) *Logger {
	var slogLevel slog.Level
	switch level {
	case DEBUG:
		slogLevel = slog.LevelDebug
	case INFO:
		slogLevel = slog.LevelInfo
	case WARN:
		slogLevel = slog.LevelWarn
	case ERROR:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: slogLevel,
	}
	
	handler := slog.NewTextHandler(os.Stdout, opts)
	
	return &Logger{
		slog: slog.New(handler),
	}
}

// Debug logs a debug message with structured fields
func (l *Logger) Debug(message string, fields ...interface{}) {
	l.slog.Debug(message, fields...)
}

// Info logs an info message with structured fields
func (l *Logger) Info(message string, fields ...interface{}) {
	l.slog.Info(message, fields...)
}

// Warn logs a warning message with structured fields
func (l *Logger) Warn(message string, fields ...interface{}) {
	l.slog.Warn(message, fields...)
}

// Error logs an error message with structured fields
func (l *Logger) Error(message string, fields ...interface{}) {
	l.slog.Error(message, fields...)
}

// With returns a new logger with additional context fields
func (l *Logger) With(fields ...interface{}) *Logger {
	return &Logger{
		slog: l.slog.With(fields...),
	}
}

// WithFields is an alias for With (for backward compatibility)
func (l *Logger) WithFields(fields ...interface{}) *Logger {
	return l.With(fields...)
}

// GetSlog returns the underlying slog.Logger for advanced usage
func (l *Logger) GetSlog() *slog.Logger {
	return l.slog
}
