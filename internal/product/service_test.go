package product

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_ListProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := []*Product{
			{ID: "p1", Name: "Prod 1"},
			{ID: "p2", Name: "Prod 2"},
		}

		mockRepo.EXPECT().
			List(ctx).
			Return(expected, nil)

		res, err := svc.ListProducts(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().
			List(ctx).
			Return(nil, errors.New("db error"))

		res, err := svc.ListProducts(ctx)
		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestService_GetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := &Product{ID: "p1", Name: "Prod 1"}

		mockRepo.EXPECT().
			GetByID(ctx, "p1").
			Return(expected, nil)

		res, err := svc.GetProduct(ctx, "p1")
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByID(ctx, "x99").
			Return(nil, errors.New("not found"))

		res, err := svc.GetProduct(ctx, "x99")
		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, "not found", err.Error())
	})
}
