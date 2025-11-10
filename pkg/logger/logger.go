package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents log level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// String returns string representation of log level
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging
type Logger struct {
	level  Level
	logger *log.Logger
}

// New creates a new logger instance
func New(level Level) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...interface{}) {
	if l.level <= DEBUG {
		l.log(DEBUG, message, fields...)
	}
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...interface{}) {
	if l.level <= INFO {
		l.log(INFO, message, fields...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...interface{}) {
	if l.level <= WARN {
		l.log(WARN, message, fields...)
	}
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...interface{}) {
	if l.level <= ERROR {
		l.log(ERROR, message, fields...)
	}
}

// log formats and writes log message
func (l *Logger) log(level Level, message string, fields ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s: %s", timestamp, level.String(), message)
	
	if len(fields) > 0 {
		logMessage += " |"
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				logMessage += fmt.Sprintf(" %v=%v", fields[i], fields[i+1])
			}
		}
	}
	
	l.logger.Println(logMessage)
}

// WithFields returns a new logger with additional fields
func (l *Logger) WithFields(fields ...interface{}) *Logger {
	// For simplicity, we'll just return the same logger
	// In a production system, you might want to use a more sophisticated logging library
	return l
}
