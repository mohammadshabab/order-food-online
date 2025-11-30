package product

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mohammadshabab/order-food-online/internal/apperrors"
	"github.com/mohammadshabab/order-food-online/internal/db"
	"github.com/mohammadshabab/order-food-online/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestMariaDBRepository_List(t *testing.T) {
	logger.Init("test-service", "test", 0)
	ctx := context.Background()

	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	db.Pool = db.NewTestPool(sqlDB)

	repo := NewMariaDBRepository()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "price", "category"}).
			AddRow("p1", "Burger", 150, "Food").
			AddRow("p2", "Pizza", 200, "Food")

		mock.ExpectQuery("SELECT id, name, price, category FROM products").
			WillReturnRows(rows)

		products, err := repo.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, products, 2)
		assert.Equal(t, "p1", products[0].ID)
	})

	t.Run("query fails", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, price, category FROM products").
			WillReturnError(errors.New("db failed"))

		products, err := repo.List(ctx)
		assert.Nil(t, products)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Contains(t, appErr.Err.Error(), "db failed")
	})

	t.Run("scan fails", func(t *testing.T) {
		// NULL values force scan error
		rows := sqlmock.NewRows([]string{"id", "name", "price", "category"}).
			AddRow(nil, nil, nil, nil)

		mock.ExpectQuery("SELECT id, name, price, category FROM products").
			WillReturnRows(rows)

		products, err := repo.List(ctx)
		assert.Nil(t, products)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, "failed to scan product row", appErr.Message)

		// Updated: sqlmock returns this new message format
		assert.Contains(t, appErr.Err.Error(), "converting NULL")
	})
}

func TestMariaDBRepository_GetByID(t *testing.T) {
	logger.Init("test-service", "test", 0)
	ctx := context.Background()

	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	db.Pool = db.NewTestPool(sqlDB)

	repo := NewMariaDBRepository()

	t.Run("success", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "name", "price", "category"}).
			AddRow("p1", "Burger", 150, "Food")

		mock.ExpectQuery("SELECT id, name, price, category FROM products WHERE id=?").
			WithArgs("p1").
			WillReturnRows(row)

		product, err := repo.GetByID(ctx, "p1")
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, "p1", product.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, price, category FROM products WHERE id=?").
			WithArgs("p999").
			WillReturnError(sql.ErrNoRows)

		product, err := repo.GetByID(ctx, "p999")
		assert.Nil(t, product)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)

		// message EXACTLY matches your code
		assert.Equal(t, "product not found with id p999", appErr.Message)

		// sql ErrNoRows message
		assert.Contains(t, appErr.Err.Error(), "sql: no rows")
	})
}
