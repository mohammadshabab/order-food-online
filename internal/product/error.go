package product

import "github.com/mohammadshabab/order-food-online/internal/apperrors"

var (
	ErrProductNotFound = apperrors.NotFound("product not found", nil)
)
