package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mohammadshabab/order-food-online/config"
	"github.com/mohammadshabab/order-food-online/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestAPIKeyMiddleware(t *testing.T) {
	// Initialize logger to avoid nil panics
	logger.Init("test-service", "test", 0)

	cfg := config.Config{
		APIKey: "my-secret-key",
	}

	nextHandler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
	}

	t.Run("health path bypasses middleware", func(t *testing.T) {
		e := echo.New()
		mw := NewAPIKeyMiddleware(cfg)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/health")

		err := mw(nextHandler)(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("valid API key", func(t *testing.T) {
		e := echo.New()
		mw := NewAPIKeyMiddleware(cfg)

		req := httptest.NewRequest(http.MethodGet, "/some-path", nil)
		req.Header.Set("api_key", "my-secret-key")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/some-path")

		err := mw(nextHandler)(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("missing API key", func(t *testing.T) {
		e := echo.New()
		mw := NewAPIKeyMiddleware(cfg)

		req := httptest.NewRequest(http.MethodGet, "/some-path", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/some-path")

		err := mw(nextHandler)(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var resp map[string]string
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "unauthorized", resp["message"])
	})

	t.Run("wrong API key", func(t *testing.T) {
		e := echo.New()
		mw := NewAPIKeyMiddleware(cfg)

		req := httptest.NewRequest(http.MethodGet, "/some-path", nil)
		req.Header.Set("api_key", "wrong-key")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/some-path")

		err := mw(nextHandler)(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var resp map[string]string
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "unauthorized", resp["message"])
	})
}
