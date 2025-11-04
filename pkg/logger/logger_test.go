package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSimpleLoggerCreation tests that logger can be created
func TestSimpleLoggerCreation(t *testing.T) {
	log := NewSimpleLogger()
	assert.NotNil(t, log)
}

// TestLoggerInfo tests Info method
func TestLoggerInfo(t *testing.T) {
	log := NewSimpleLogger()
	assert.NotPanics(t, func() {
		log.Info("test message")
		msg := "test message"
		log.Info(msg)
	})
}

// TestLoggerError tests Error method
func TestLoggerError(t *testing.T) {
	log := NewSimpleLogger()
	assert.NotPanics(t, func() {
		log.Error("error message")
		msg := "error message"
		log.Error(msg)
	})
}

// TestLoggerWarn tests Warn method
func TestLoggerWarn(t *testing.T) {
	log := NewSimpleLogger()
	assert.NotPanics(t, func() {
		log.Warn("warning message")
		msg := "warning message"
		log.Warn(msg)
	})
}

// TestLoggerInterface tests that SimpleLogger implements Logger interface
func TestLoggerInterface(t *testing.T) {
	var log Logger = NewSimpleLogger()
	assert.NotNil(t, log)

	assert.NotPanics(t, func() {
		log.Info("test")
		log.Error("test")
		log.Warn("test")
	})
}

// TestFormatMessage tests the formatMessage helper function
func TestFormatMessage(t *testing.T) {
	tests := []struct {
		name   string
		msg    string
		args   []interface{}
		expect string
	}{
		{
			name:   "message without args",
			msg:    "hello world",
			args:   []interface{}{},
			expect: "hello world",
		},
		{
			name:   "message with key-value pairs",
			msg:    "user created",
			args:   []interface{}{"id", "123", "name", "John"},
			expect: "user created id=123 name=John",
		},
		{
			name:   "message with printf-style args",
			msg:    "value is %d",
			args:   []interface{}{42},
			expect: "value is 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatMessage(tt.msg, tt.args...)
			assert.Equal(t, tt.expect, result)
		})
	}
}
