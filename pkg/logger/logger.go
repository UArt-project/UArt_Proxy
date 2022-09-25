// Package logger contains custom logger for the application.
package logger

import (
	"io"
	"log"
)

// Logger is a custom logger implementation for application use.
type Logger struct {
	logger *log.Logger
}

// NewLogger returns new Logger instance.
func NewLogger(output io.Writer, prefix string) *Logger {
	logger := log.New(output, prefix+" ", log.LstdFlags)

	return &Logger{
		logger: logger,
	}
}

// Error logs error in a Printf way.
func (l Logger) Error(format string, args ...any) {
	l.logger.Printf("ERROR "+format, args...)
}

// Info logs information in a Printf way.
func (l Logger) Info(format string, args ...any) {
	l.logger.Printf("INFO "+format, args...)
}

// Fatal logs fatal error in a Panicf way.
func (l Logger) Fatal(format string, args ...any) {
	l.logger.Panicf("FATAL "+format, args...)
}
