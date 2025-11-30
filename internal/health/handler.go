package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mohammadshabab/order-food-online/internal/db"
	"github.com/mohammadshabab/order-food-online/internal/logger"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Returns 200 if healthy, 503 if unhealthy
func (h *Handler) Check(c echo.Context) error {
	ctx := c.Request().Context()

	if err := db.Pool.DB.PingContext(ctx); err != nil {
		logger.Log().Warn("health check failed: database unreachable", "error", err)
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "unhealthy",
			"error":  "database connection failed",
		})
	}
	logger.Log().Debug("health check passed")
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "order-food-online",
	})
}

// Used by load balancers for frequent checks
func (h *Handler) Ping(c echo.Context) error {
	logger.Log().Info("ping received")
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
