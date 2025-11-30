package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSensitiveAndFilterSensitive(t *testing.T) {
	// Non-prod environment: values should remain unmasked
	env = "dev"

	kvs := []any{
		"user", "alice",
		Sensitive("password", "secret123"),
	}

	filtered := filterSensitive(kvs)
	assert.Len(t, filtered, 3) // corrected length

	// First two values remain the same
	assert.Equal(t, "user", filtered[0])
	assert.Equal(t, "alice", filtered[1])

	// Sensitive value should remain visible in dev
	attr, ok := filtered[2].(slog.Attr)
	assert.True(t, ok, "Should be slog.Attr")
	assert.Equal(t, "password", attr.Key)
	assert.Equal(t, "secret123", attr.Value.Any())

	// Prod environment: sensitive value should be masked
	env = "prod"

	kvsProd := []any{
		"user", "bob",
		Sensitive("password", "secret123"),
	}

	filteredProd := filterSensitive(kvsProd)
	assert.Len(t, filteredProd, 3) // corrected length

	// Sensitive value should be masked in prod
	attrProd, ok := filteredProd[2].(slog.Attr)
	assert.True(t, ok)
	assert.Equal(t, "password", attrProd.Key)
	assert.Equal(t, "****", attrProd.Value.Any())

	// Other values remain unchanged
	assert.Equal(t, "user", filteredProd[0])
	assert.Equal(t, "bob", filteredProd[1])
}
