package order

import (
	"testing"

	"github.com/mohammadshabab/order-food-online/internal/apperrors"
	"github.com/stretchr/testify/assert"
)

func TestOrderReq_Validate(t *testing.T) {
	t.Run("nil items", func(t *testing.T) {
		order := &OrderReq{Items: nil}
		err := order.Validate()
		assert.NotNil(t, err)
		assert.Equal(t, 400, err.Code)
		assert.Equal(t, "items is required", err.Message)
		assert.Equal(t, apperrors.LevelWarn, err.Level)
	})

	t.Run("empty items slice", func(t *testing.T) {
		order := &OrderReq{Items: &[]OrderItem{}}
		err := order.Validate()
		assert.NotNil(t, err)
		assert.Equal(t, 422, err.Code)
		assert.Equal(t, "order must have at least one item", err.Message)
		assert.Equal(t, apperrors.LevelWarn, err.Level)
	})

	t.Run("missing product ID", func(t *testing.T) {
		order := &OrderReq{Items: &[]OrderItem{
			{ProductID: "", Quantity: 1},
		}}
		err := order.Validate()
		assert.NotNil(t, err)
		assert.Equal(t, 400, err.Code)
		assert.Equal(t, "product ID is required", err.Message)
		assert.Equal(t, apperrors.LevelWarn, err.Level)
	})

	t.Run("quantity <= 0", func(t *testing.T) {
		order := &OrderReq{Items: &[]OrderItem{
			{ProductID: "prod1", Quantity: 0},
		}}
		err := order.Validate()
		assert.NotNil(t, err)
		assert.Equal(t, 400, err.Code)
		assert.Equal(t, "quantity must be greater than 0", err.Message)
		assert.Equal(t, apperrors.LevelWarn, err.Level)
	})

	t.Run("valid order", func(t *testing.T) {
		order := &OrderReq{Items: &[]OrderItem{
			{ProductID: "prod1", Quantity: 2},
			{ProductID: "prod2", Quantity: 1},
		}}
		err := order.Validate()
		assert.Nil(t, err)
	})
}
