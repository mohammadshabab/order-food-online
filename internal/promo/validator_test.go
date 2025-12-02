package promo

import (
	"log/slog"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mohammadshabab/order-food-online/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestValidator_Validate(t *testing.T) {
	logger.Init("test-service", "test", slog.LevelInfo)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := NewMockCache(ctrl)

	// Create a validator using mocked cache
	validator := &Validator{cache: mockCache}

	t.Run("invalid code length", func(t *testing.T) {
		err := validator.Validate("SHORT")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid coupon code format")
	})

	t.Run("code not found in cache", func(t *testing.T) {
		code := "VALID123"
		mockCache.EXPECT().Get(code).Return(Coupon{}, false)

		err := validator.Validate(code)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid coupon code")
	})

	t.Run("code found but insufficient file count", func(t *testing.T) {
		code := "VALID123"
		mockCache.EXPECT().Get(code).Return(Coupon{
			Code:      code,
			FileCount: 1, // less than 2
		}, true)

		err := validator.Validate(code)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not in enough files")
	})

	t.Run("code found with sufficient file count", func(t *testing.T) {
		code := "VALID123"
		mockCache.EXPECT().Get(code).Return(Coupon{
			Code:      code,
			FileCount: 2,
		}, true)

		err := validator.Validate(code)
		require.NoError(t, err)
	})
}
