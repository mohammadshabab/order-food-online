package logger

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndExtractAttrs(t *testing.T) {
	ctx := context.Background()

	ctx = AddAttrs(ctx, Attributes{"user": "john", "id": 10})
	ctx = AddAttrs(ctx, Attributes{"role": "admin"}) // merge with previous

	attrs := extractAttrs(ctx)
	assert.Len(t, attrs, 3)

	// Verify that all keys exist
	keys := map[string]any{}
	for _, a := range attrs {
		keys[a.Key] = a.Value.Any()
	}

	assert.Equal(t, "john", keys["user"])
	assert.Equal(t, int64(10), keys["id"].(int64)) // <- cast to int64
	assert.Equal(t, "admin", keys["role"])
}

func TestAttrsToAny(t *testing.T) {
	ctx := context.Background()
	ctx = AddAttrs(ctx, Attributes{"x": 1, "y": 2})

	attrs := extractAttrs(ctx)
	anySlice := attrsToAny(attrs)

	assert.Len(t, anySlice, 2)

	m := map[string]any{}
	for _, a := range anySlice {
		if attr, ok := a.(slog.Attr); ok {
			m[attr.Key] = attr.Value.Any()
		}
	}
	assert.Equal(t, int64(1), m["x"])
	assert.Equal(t, int64(2), m["y"])
}

func TestLoggingFunctions(t *testing.T) {
	// Redirect logs to discard to avoid clutter
	Init("test-service", "test", slog.LevelInfo)

	ctx := context.Background()
	ctx = AddAttrs(ctx, Attributes{"user": "alice"})

	// Ensure all logging functions do not panic
	assert.NotPanics(t, func() { Info(ctx, "info message", "key", "val") })
	assert.NotPanics(t, func() { Debug(ctx, "debug message", "dummy", 0) })
	assert.NotPanics(t, func() { Warn(ctx, "warn message", "secret", "123") })
	assert.NotPanics(t, func() { Error(ctx, "error message", "err", nil) })
}
