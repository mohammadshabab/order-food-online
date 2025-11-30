package health

import (
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

func TestRegisterRoutes(t *testing.T) {
	logger.Init("test-service", "test", slog.LevelInfo)

	// Setup sqlmock for db.Pool
	mockDB, _, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	defer mockDB.Close()
	db.Pool = &db.SQLPool{DB: mockDB}

	e := echo.New()
	Register(e)

	t.Run("GET /health route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Contains(t, []int{http.StatusOK, http.StatusServiceUnavailable}, rec.Code)
	})

	t.Run("GET /health/ping route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health/ping", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"status":"ok"`)
	})
}
