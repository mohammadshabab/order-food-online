package health

import "github.com/labstack/echo/v4"

// Register registers health-check endpoints on the provided Echo router.
func Register(e *echo.Echo) {
	h := NewHandler()
	e.GET("/health", h.Check)
	e.GET("/health/ping", h.Ping)
}
