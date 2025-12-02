package order

import (
	"github.com/labstack/echo/v4"
	"github.com/mohammadshabab/order-food-online/internal/promo"
)

func Setup(e *echo.Echo, repo Repository, promoValidator *promo.Validator) {
	svc := NewService(repo, promoValidator)
	h := NewHandler(svc)

	e.POST("/order", h.CreateOrder)
}
