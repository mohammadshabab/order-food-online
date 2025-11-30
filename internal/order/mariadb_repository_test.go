package order

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mohammadshabab/order-food-online/internal/apperrors"
	"github.com/mohammadshabab/order-food-online/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestMariaDBRepository_Create(t *testing.T) {
	ctx := context.Background()

	// Setup sqlmock
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	db.Pool = db.NewTestPool(sqlDB) // keep db.Pool usage

	repo := NewMariaDBRepository()

	t.Run("success", func(t *testing.T) {
		orderObj := &Order{
			ID:         "order123",
			CouponCode: nil,
			Items: []OrderItem{
				{ProductID: "p1", Quantity: 2},
			},
		}

		// Insert order
		mock.ExpectExec("INSERT INTO orders").
			WithArgs(orderObj.ID, orderObj.CouponCode, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Fetch product
		rows := sqlmock.NewRows([]string{"id", "name", "category", "price"}).
			AddRow("p1", "Burger", "Food", 150)
		mock.ExpectQuery("SELECT id, name, category, price FROM products").
			WithArgs("p1").
			WillReturnRows(rows)

		// Insert order_items
		mock.ExpectExec("INSERT INTO order_items").
			WithArgs(orderObj.ID, "p1", 2).
			WillReturnResult(sqlmock.NewResult(1, 1))

		res, err := repo.Create(ctx, orderObj)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "p1", res.Products[0].ID)
	})

	t.Run("order insert fails", func(t *testing.T) {
		orderObj := &Order{ID: "order123"}

		mock.ExpectExec("INSERT INTO orders").
			WillReturnError(errors.New("db error"))

		res, err := repo.Create(ctx, orderObj)
		assert.Nil(t, res)
		assert.Error(t, err)
		// match your current implementation error
		assert.Contains(t, err.Error(), "DB Exec failed")
	})

	t.Run("product not found", func(t *testing.T) {
		orderObj := &Order{
			ID: "order123",
			Items: []OrderItem{
				{ProductID: "p999", Quantity: 1},
			},
		}

		mock.ExpectExec("INSERT INTO orders").
			WithArgs(orderObj.ID, orderObj.CouponCode, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT id, name, category, price FROM products").
			WithArgs("p999").
			WillReturnError(sql.ErrNoRows)

		res, err := repo.Create(ctx, orderObj)
		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, 404, err.(*apperrors.AppError).Code)
	})

	t.Run("product scan fails", func(t *testing.T) {
		orderObj := &Order{
			ID: "order123",
			Items: []OrderItem{
				{ProductID: "p1", Quantity: 1},
			},
		}

		mock.ExpectExec("INSERT INTO orders").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// return invalid rows to cause Scan error
		rows := sqlmock.NewRows([]string{"id", "name", "category", "price"}).
			AddRow(nil, nil, nil, nil)
		mock.ExpectQuery("SELECT id, name, category, price FROM products").
			WithArgs("p1").
			WillReturnRows(rows)

		res, err := repo.Create(ctx, orderObj)
		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch product details")
	})

	t.Run("insert order item fails", func(t *testing.T) {
		orderObj := &Order{
			ID: "order123",
			Items: []OrderItem{
				{ProductID: "p1", Quantity: 2},
			},
		}

		mock.ExpectExec("INSERT INTO orders").
			WillReturnResult(sqlmock.NewResult(1, 1))

		rows := sqlmock.NewRows([]string{"id", "name", "category", "price"}).
			AddRow("p1", "Burger", "Food", 150)
		mock.ExpectQuery("SELECT id, name, category, price FROM products").
			WithArgs("p1").
			WillReturnRows(rows)

		mock.ExpectExec("INSERT INTO order_items").
			WillReturnError(errors.New("item insert error"))

		res, err := repo.Create(ctx, orderObj)
		assert.Nil(t, res)
		assert.Error(t, err)
		// match actual error returned by your repo
		assert.Contains(t, err.Error(), "DB Exec failed")
	})
}
