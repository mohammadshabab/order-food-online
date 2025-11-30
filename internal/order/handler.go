package order

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mohammadshabab/order-food-online/internal/apperrors"
	"github.com/mohammadshabab/order-food-online/internal/logger"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateOrder(c echo.Context) error {
	ctx := c.Request().Context()

	var req OrderReq
	if err := c.Bind(&req); err != nil {
		logger.Warn(ctx, ErrOrderInvalid.Message, "error", err.Error())
		return c.JSON(ErrOrderInvalid.Code, ErrOrderInvalid)
	}

	if appErr := req.Validate(); appErr != nil {
		return c.JSON(appErr.Code, appErr)
	}

	order, err := h.svc.CreateOrder(ctx, &req)
	if err != nil {
		appErr := apperrors.Internal("failed to create order", err)
		return c.JSON(appErr.Code, appErr)
	}

	return c.JSON(http.StatusOK, order)
}
