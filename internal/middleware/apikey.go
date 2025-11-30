package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mohammadshabab/order-food-online/config"
	"github.com/mohammadshabab/order-food-online/internal/logger"
)

func NewAPIKeyMiddleware(cfg config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/health" || c.Path() == "/ping" {
				return next(c)
			}

			key := c.Request().Header.Get("api_key")
			if key == "" || key != cfg.APIKey {
				logger.Log().Warn("unauthorized request", "path", c.Path(), "method", c.Request().Method)
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
			}
			return next(c)
		}
	}
}
