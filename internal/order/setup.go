package order

import "github.com/labstack/echo/v4"

func Setup(e *echo.Echo, repo Repository) {
	svc := NewService(repo)
	h := NewHandler(svc)

	e.POST("/order", h.CreateOrder)
}
