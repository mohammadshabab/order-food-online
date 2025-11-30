package product

import (
	"net/http"

	"github.com/google/uuid"
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

func (h *Handler) ListProducts(c echo.Context) error {
	ctx := c.Request().Context()
	logger.Info(ctx, "list products called")

	res, err := h.svc.ListProducts(ctx)
	if err != nil {
		appErr := apperrors.Internal("failed to list products", err)
		logger.Error(ctx, appErr.Message, "error", appErr.Error())
		return c.JSON(appErr.Code, appErr)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetProduct(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("productId")

	// OpenAPI validation: check if ID is provided
	if id == "" {
		appErr := apperrors.BadRequest("invalid ID supplied", nil)
		logger.Warn(ctx, appErr.Message)
		return c.JSON(appErr.Code, appErr)
	}

	// Validate UUID format (OpenAPI requirement)
	if _, err := uuid.Parse(id); err != nil {
		appErr := apperrors.BadRequest("invalid ID supplied", err)
		logger.Warn(ctx, appErr.Message, "id", id)
		return c.JSON(appErr.Code, appErr)
	}

	logger.Info(ctx, "get product called", "id", id)

	res, err := h.svc.GetProduct(ctx, id)
	if err != nil {
		appErr := apperrors.NotFound("product not found", err)
		logger.Warn(ctx, appErr.Message, "id", id)
		return c.JSON(appErr.Code, appErr)
	}

	return c.JSON(http.StatusOK, res)
}
