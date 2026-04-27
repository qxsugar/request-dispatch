package request_dispatch

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected LogLevel
	}{
		{"uppercase DEBUG", "DEBUG", LogLevelDebug},
		{"lowercase debug", "debug", LogLevelDebug},
		{"mixed case Debug", "Debug", LogLevelDebug},
		{"uppercase INFO", "INFO", LogLevelInfo},
		{"lowercase info", "info", LogLevelInfo},
		{"uppercase ERROR", "ERROR", LogLevelError},
		{"lowercase error", "error", LogLevelError},
		{"invalid level", "INVALID", LogLevelError},
		{"empty string", "", LogLevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLogLevel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoggerDebug(t *testing.T) {
	t.Run("debug logs at DEBUG level", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &Logger{
			logLevel:    LogLevelDebug,
			loggerDebug: log.New(&buf, "DEBUG: ", 0),
		}

		logger.Debug("test message")
		assert.Contains(t, buf.String(), "test message")
	})

	t.Run("debug logs not shown at INFO level", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &Logger{
			logLevel:    LogLevelInfo,
			loggerDebug: log.New(&buf, "DEBUG: ", 0),
		}

		logger.Debug("test message")
		assert.Empty(t, buf.String())
	})
}

func TestLoggerInfo(t *testing.T) {
	t.Run("info logs at INFO level", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &Logger{
			logLevel:   LogLevelInfo,
			loggerInfo: log.New(&buf, "INFO: ", 0),
		}

		logger.Info("test message")
		assert.Contains(t, buf.String(), "test message")
	})

	t.Run("info logs at DEBUG level", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &Logger{
			logLevel:   LogLevelDebug,
			loggerInfo: log.New(&buf, "INFO: ", 0),
		}

		logger.Info("test message")
		assert.Contains(t, buf.String(), "test message")
	})

	t.Run("info logs not shown at ERROR level", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &Logger{
			logLevel:   LogLevelError,
			loggerInfo: log.New(&buf, "INFO: ", 0),
		}

		logger.Info("test message")
		assert.Empty(t, buf.String())
	})
}

func TestLoggerError(t *testing.T) {
	t.Run("error logs at all levels", func(t *testing.T) {
		levels := []LogLevel{LogLevelDebug, LogLevelInfo, LogLevelError}
		for _, level := range levels {
			var buf bytes.Buffer
			logger := &Logger{
				logLevel:    level,
				loggerError: log.New(&buf, "ERROR: ", 0),
			}

			logger.Error("test message")
			assert.Contains(t, buf.String(), "test message")
		}
	})
}

func TestNewLogger(t *testing.T) {
	t.Run("creates logger with DEBUG level", func(t *testing.T) {
		logger := NewLogger("DEBUG")
		assert.Equal(t, LogLevelDebug, logger.logLevel)
		assert.NotNil(t, logger.loggerDebug)
		assert.NotNil(t, logger.loggerInfo)
		assert.NotNil(t, logger.loggerError)
	})

	t.Run("creates logger with lowercase debug", func(t *testing.T) {
		logger := NewLogger("debug")
		assert.Equal(t, LogLevelDebug, logger.logLevel)
	})

	t.Run("creates logger with INFO level", func(t *testing.T) {
		logger := NewLogger("INFO")
		assert.Equal(t, LogLevelInfo, logger.logLevel)
	})

	t.Run("creates logger with ERROR level", func(t *testing.T) {
		logger := NewLogger("ERROR")
		assert.Equal(t, LogLevelError, logger.logLevel)
	})

	t.Run("defaults to ERROR level for invalid input", func(t *testing.T) {
		logger := NewLogger("INVALID")
		assert.Equal(t, LogLevelError, logger.logLevel)
	})

	t.Run("loggers write to stdout", func(t *testing.T) {
		logger := NewLogger("DEBUG")
		assert.NotNil(t, logger.loggerDebug)
		assert.NotNil(t, logger.loggerInfo)
		assert.NotNil(t, logger.loggerError)
	})
}

func TestLoggerMultipleArgs(t *testing.T) {
	t.Run("logs multiple arguments", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &Logger{
			logLevel:    LogLevelDebug,
			loggerDebug: log.New(&buf, "", 0),
		}

		logger.Debug("arg1", "arg2", "arg3")
		output := buf.String()
		assert.Contains(t, output, "arg1")
		assert.Contains(t, output, "arg2")
		assert.Contains(t, output, "arg3")
	})
}
