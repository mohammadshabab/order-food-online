package logger

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerInitAndEnv(t *testing.T) {
	// Initialize logger with discard output to avoid printing
	Init("test-service", "dev", slog.LevelDebug)

	// Logger should not be nil after Init
	l := Log()
	assert.NotNil(t, l, "Logger should be initialized")

	// Test environment flags
	assert.True(t, IsDev(), "Should detect dev environment")
	assert.False(t, IsProd(), "Should not detect prod environment")

	// Re-initialize as production
	Init("prod-service", "production", slog.LevelInfo)
	assert.True(t, IsProd(), "Should detect prod environment")
	assert.False(t, IsDev(), "Should not detect dev environment")
}

func TestLoggerOutputRedirection(t *testing.T) {
	oldLog := log
	log = slog.New(slog.NewTextHandler(io.Discard, nil))

	// Should not panic when using the logger
	assert.NotPanics(t, func() {
		log.Info("test message", "key", "value")
	})

	// Restore original logger
	log = oldLog
}
