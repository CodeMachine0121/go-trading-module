package logger

import (
	"fmt"
	"log"
	"os"
)

// Logger defines the logging interface for the application.
type Logger interface {
	// Info logs an informational message.
	Info(msg string, args ...interface{})

	// Error logs an error message.
	Error(msg string, args ...interface{})

	// Warn logs a warning message.
	Warn(msg string, args ...interface{})
}

// SimpleLogger is a basic implementation of Logger using standard library log.
type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
}

// NewSimpleLogger creates a new SimpleLogger instance.
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		warnLogger:  log.New(os.Stdout, "[WARN] ", log.LstdFlags),
	}
}

// Info logs an informational message.
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	l.infoLogger.Println(formatMessage(msg, args...))
}

// Error logs an error message.
func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	l.errorLogger.Println(formatMessage(msg, args...))
}

// Warn logs a warning message.
func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	l.warnLogger.Println(formatMessage(msg, args...))
}

// formatMessage formats a message with key-value pairs or printf-style args
func formatMessage(msg string, args ...interface{}) string {
	// If no args, just return the message
	if len(args) == 0 {
		return msg
	}

	// Check if args are key-value pairs (even number of args)
	if len(args)%2 == 0 {
		// Format as key-value pairs
		output := msg
		for i := 0; i < len(args); i += 2 {
			output += fmt.Sprintf(" %v=%v", args[i], args[i+1])
		}
		return output
	}

	// Otherwise, treat as printf-style arguments
	return fmt.Sprintf(msg, args...)
}
