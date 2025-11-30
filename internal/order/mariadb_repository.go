package order

import (
	"context"
	"database/sql"
	"time"

	"github.com/mohammadshabab/order-food-online/internal/apperrors"
	"github.com/mohammadshabab/order-food-online/internal/db"
	"github.com/mohammadshabab/order-food-online/internal/logger"
)

type MariaDBRepository struct{}

func NewMariaDBRepository() Repository {
	return &MariaDBRepository{}
}

func (r *MariaDBRepository) Create(ctx context.Context, order *Order) (*Order, error) {

	// Insert order (UUID provided from service)
	query := `INSERT INTO orders (id, coupon_code, created_at) VALUES (?, ?, ?)`
	_, err := db.Pool.Exec(ctx, query, order.ID, order.CouponCode, time.Now())
	if err != nil {
		appErr := apperrors.Internal("failed to create order", err)
		logger.Error(ctx, appErr.Message, "error", err.Error())
		return nil, appErr
	}

	// Insert items + collect productRefs
	productRefs := make([]ProductRef, 0)

	for _, item := range order.Items {
		// Fetch product details for response
		var p ProductRef
		productQuery := `SELECT id, name, category, price FROM products WHERE id = ?`
		err := db.Pool.QueryRow(ctx, productQuery, item.ProductID).
			Scan(&p.ID, &p.Name, &p.Category, &p.Price)

		if err != nil {
			if err == sql.ErrNoRows {
				// Return 404 if product not found
				appErr := apperrors.NotFound("product not found", nil)
				logger.Warn(ctx, appErr.Message, "productId", item.ProductID)
				return nil, appErr
			}
			appErr := apperrors.Internal("failed to fetch product details", err)
			logger.Error(ctx, appErr.Message, "error", err.Error())
			return nil, appErr
		}

		productRefs = append(productRefs, p)

		// Insert order item
		itemQuery := `INSERT INTO order_items (id, order_id, product_id, quantity, created_at)
		              VALUES (UUID(), ?, ?, ?, NOW())`

		_, err = db.Pool.Exec(ctx, itemQuery, order.ID, item.ProductID, item.Quantity)
		if err != nil {
			appErr := apperrors.Internal("failed to insert order item", err)
			logger.Error(ctx, appErr.Message, "error", err.Error())
			return nil, appErr
		}
	}

	order.Products = productRefs
	return order, nil
}
