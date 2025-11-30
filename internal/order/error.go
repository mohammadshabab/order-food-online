package order

import "github.com/mohammadshabab/order-food-online/internal/apperrors"

var (
	ErrOrderNotFound   = apperrors.NotFound("order not found", nil)
	ErrOrderInvalid    = apperrors.BadRequest("invalid order request", nil)
	ErrOrderValidation = apperrors.Wrap(422, "validation failed", apperrors.LevelWarn, nil)
)

func (or *OrderReq) Validate() *apperrors.AppError {
	if or.Items == nil {
		return apperrors.Wrap(400, "items is required", apperrors.LevelWarn, nil)
	}

	if len(*or.Items) == 0 {
		return apperrors.Wrap(422, "order must have at least one item", apperrors.LevelWarn, nil)
	}

	// Validate each item
	for _, i := range *or.Items {

		// productId missing or empty
		if i.ProductID == "" {
			return apperrors.Wrap(400, "product ID is required", apperrors.LevelWarn, nil)
		}

		// quantity missing or invalid
		if i.Quantity <= 0 {
			return apperrors.Wrap(400, "quantity must be greater than 0", apperrors.LevelWarn, nil)
		}
	}

	return nil
}
