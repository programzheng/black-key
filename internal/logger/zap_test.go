package logger_test

import (
	"strings"
	"testing"

	"github.com/programzheng/black-key/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestDebug(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := logger.NewWithLogger(zap.New(core))

	logger.Debug("test message", zap.String("key", "value"))

	if logs.Len() != 1 {
		t.Errorf("expected 1 log message, got %d", logs.Len())
	}

	loggedMsg := logs.All()[0].Message
	if !strings.Contains(loggedMsg, "test message") {
		t.Errorf("expected log message to contain %q, got %q", "test message", loggedMsg)
	}

	loggedField := logs.All()[0].ContextMap()["key"]
	if loggedField != "value" {
		t.Errorf("expected logged field value to be %q, got %q", "value", loggedField)
	}
}

func TestInfo(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := logger.NewWithLogger(zap.New(core))

	logger.Info("test message", zap.String("key", "value"))

	if logs.Len() != 1 {
		t.Errorf("expected 1 log message, got %d", logs.Len())
	}

	loggedMsg := logs.All()[0].Message
	if !strings.Contains(loggedMsg, "test message") {
		t.Errorf("expected log message to contain %q, got %q", "test message", loggedMsg)
	}

	loggedField := logs.All()[0].ContextMap()["key"]
	if loggedField != "value" {
		t.Errorf("expected logged field value to be %q, got %q", "value", loggedField)
	}
}

func TestError(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := logger.NewWithLogger(zap.New(core))

	logger.Error("test message", zap.String("key", "value"))

	if logs.Len() != 1 {
		t.Errorf("expected 1 log message, got %d", logs.Len())
	}

	loggedMsg := logs.All()[0].Message
	if !strings.Contains(loggedMsg, "test message") {
		t.Errorf("expected log message to contain %q, got %q", "test message", loggedMsg)
	}

	loggedField := logs.All()[0].ContextMap()["key"]
	if loggedField != "value" {
		t.Errorf("expected logged field value to be %q, got %q", "value", loggedField)
	}
}

func TestSync(t *testing.T) {
	core, _ := observer.New(zap.DebugLevel)
	logger := logger.NewWithLogger(zap.New(core))

	err := logger.Sync()
	if err != nil {
		t.Errorf("expected Sync method to return nil, got %v", err)
	}
}
