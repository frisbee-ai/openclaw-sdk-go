package types

import (
	"bytes"
	"context"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewDefaultLoggerWithWriter(buf)

	logger.Info("test message %s", "world")
	logger.Debug("debug message")
	logger.Warn("warning message")
	logger.Error("error message")

	output := buf.String()
	if output == "" {
		t.Error("expected logger output")
	}
}

func TestNewDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger()
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	// Verify it implements Logger interface
	var _ Logger = logger

	// Verify internal loggers are not nil
	if logger.debug == nil {
		t.Error("expected debug logger to be non-nil")
	}
	if logger.info == nil {
		t.Error("expected info logger to be non-nil")
	}
	if logger.warn == nil {
		t.Error("expected warn logger to be non-nil")
	}
	if logger.error == nil {
		t.Error("expected error logger to be non-nil")
	}

	// Methods should not panic
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	// With formatting args
	logger.Debug("debug %s %d", "arg", 123)
	logger.Info("info %s %d", "arg", 123)
	logger.Warn("warn %s %d", "arg", 123)
	logger.Error("error %s %d", "arg", 123)
}

func TestNopLogger(t *testing.T) {
	logger := &NopLogger{}
	// Should not panic
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
}

func TestLoggerInterface(t *testing.T) {
	// Verify DefaultLogger implements Logger
	var _ Logger = &DefaultLogger{}
	// Verify NopLogger implements Logger
	var _ Logger = &NopLogger{}
}

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	logger := &NopLogger{}

	ctx = WithContext(ctx, logger)
	retrieved, ok := FromContext(ctx)

	if !ok {
		t.Error("expected to retrieve logger from context")
	}
	if retrieved != logger {
		t.Error("expected to retrieve same logger")
	}
}
