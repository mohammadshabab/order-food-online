package health

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mohammadshabab/order-food-online/internal/db"
	"github.com/mohammadshabab/order-food-online/internal/logger"
)

func TestHandler_Check(t *testing.T) {
	// Initialize logger to avoid nil panic
	logger.Init("test-service", "test", slog.LevelInfo)

	h := NewHandler()
	e := echo.New()

	t.Run("healthy database", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		defer mockDB.Close()

		// Ping will succeed
		mock.ExpectPing().WillReturnError(nil)

		db.Pool = &db.SQLPool{DB: mockDB}

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = h.Check(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"status":"ok"`)
	})

	t.Run("unhealthy database", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		defer mockDB.Close()

		// Ping will fail
		mock.ExpectPing().WillReturnError(errors.New("db down"))

		db.Pool = &db.SQLPool{DB: mockDB}

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = h.Check(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
		assert.Contains(t, rec.Body.String(), `"status":"unhealthy"`)
		assert.Contains(t, rec.Body.String(), `"error":"database connection failed"`)
	})
}

func TestHandler_Ping(t *testing.T) {
	// Initialize logger
	logger.Init("test-service", "test", slog.LevelInfo)

	h := NewHandler()
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Ping(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"ok"`)
}
