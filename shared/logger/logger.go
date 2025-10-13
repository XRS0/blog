package logger
package logger

import (
	"log/slog"
	"os"
)

// Logger wraps slog for microservices
type Logger struct {
	*slog.Logger
}

// New creates a new logger
func New(level slog.Leveler, serviceName string) *Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	log := slog.New(handler).With("service", serviceName)
	return &Logger{Logger: log}
}

// WithContext adds context fields to the logger
func (l *Logger) WithContext(key string, value interface{}) *Logger {
	return &Logger{Logger: l.With(key, value)}
}
