package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mohammadshabab/order-food-online/internal/apperrors"
	"github.com/mohammadshabab/order-food-online/internal/db"
	"github.com/mohammadshabab/order-food-online/internal/logger"
)

type MariaDBRepository struct{}

func NewMariaDBRepository() Repository {
	return &MariaDBRepository{}
}

func (r *MariaDBRepository) List(ctx context.Context) ([]*Product, error) {
	query := `SELECT id, name, price, category FROM products`

	logger.Info(ctx, "DB Query start", "query", query)

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		appErr := apperrors.Internal("failed to list products", err)
		logger.Error(ctx, appErr.Message, "query", query, "error", err.Error())
		return nil, appErr
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category); err != nil {
			appErr := apperrors.Internal("failed to scan product row", err)
			logger.Error(ctx, appErr.Message, "query", query, "error", err.Error())
			return nil, appErr
		}
		products = append(products, &p)
	}

	logger.Info(ctx, "DB Query completed", "query", query, "rows", len(products))
	return products, nil
}

func (r *MariaDBRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	query := `SELECT id, name, price, category FROM products WHERE id=?`
	args := []any{id}

	logger.Info(ctx, "DB QueryRow start", "query", query, "args", args)

	row := db.Pool.QueryRow(ctx, query, args...)

	var p Product
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			appErr := apperrors.NotFound(fmt.Sprintf("product not found with id %s", id), err)
			logger.Warn(ctx, appErr.Message, "id", id)
			return nil, appErr
		}

		appErr := apperrors.Internal("failed to fetch product", err)
		logger.Error(ctx, appErr.Message, "id", id, "error", err.Error())
		return nil, appErr
	}

	logger.Info(ctx, "DB QueryRow completed", "id", id)
	return &p, nil
}
