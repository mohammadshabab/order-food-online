package order

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mohammadshabab/order-food-online/internal/apperrors"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	t.Run("validation fails", func(t *testing.T) {
		req := &OrderReq{Items: &[]OrderItem{}} // empty items triggers validation error

		order, err := svc.CreateOrder(context.Background(), req)
		assert.Nil(t, order)
		assert.NotNil(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok, "error should be of type *apperrors.AppError")
		assert.Equal(t, 422, appErr.Code)
		assert.Equal(t, "order must have at least one item", appErr.Message)
	})

	t.Run("repository returns error", func(t *testing.T) {
		req := &OrderReq{Items: &[]OrderItem{{ProductID: "p1", Quantity: 2}}}

		// Expect Create to be called and return an error
		mockRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("db error"))

		order, err := svc.CreateOrder(context.Background(), req)
		assert.Nil(t, order)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "db error")
	})

	t.Run("success", func(t *testing.T) {
		req := &OrderReq{Items: &[]OrderItem{{ProductID: "p1", Quantity: 2}}}
		expectedOrder := &Order{ID: "order1", Items: *req.Items}

		// Mock repository to return expected order
		mockRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, o *Order) (*Order, error) {
				o.ID = expectedOrder.ID // assign deterministic ID
				return o, nil
			})

		order, err := svc.CreateOrder(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, expectedOrder.ID, order.ID)
		assert.Equal(t, expectedOrder.Items, order.Items)
	})
}
