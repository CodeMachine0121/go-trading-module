package logger

// Logger defines the logging interface for the application.
type Logger interface {
	// Info logs an informational message.
	Info(msg string, args ...interface{})

	// Error logs an error message.
	Error(msg string, args ...interface{})

	// Warn logs a warning message.
	Warn(msg string, args ...interface{})
}

// SimpleLogger is a basic implementation of Logger.
type SimpleLogger struct{}

// NewSimpleLogger creates a new SimpleLogger instance.
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{}
}

// Info logs an informational message.
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	// Basic implementation - can be enhanced later
}

// Error logs an error message.
func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	// Basic implementation - can be enhanced later
}

// Warn logs a warning message.
func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	// Basic implementation - can be enhanced later
}
