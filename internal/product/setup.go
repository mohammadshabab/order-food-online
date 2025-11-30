package product

import "github.com/labstack/echo/v4"

func Setup(e *echo.Echo, repo Repository) {
	// Create service
	svc := NewService(repo)

	// Create handler
	h := NewHandler(svc)

	// Register routes
	e.GET("/product", h.ListProducts)
	e.GET("/product/:productId", h.GetProduct)
}
