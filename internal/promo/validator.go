package promo

import (
	"context"
	"fmt"
	"time"

	"github.com/mohammadshabab/order-food-online/internal/apperrors"
	"github.com/mohammadshabab/order-food-online/internal/logger"
)

type Validator struct {
	cache Cache
}

// New creates a new Validator and loads coupons with context timeout
func New(dir string) (*Validator, error) {
	cache := NewCache()
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second) // 30s timeout
	defer cancel()

	errCh := make(chan error)
	go func() {
		errCh <- LoadCouponsWithContext(ctx, LoaderConfig{
			Dir:         dir,
			WorkerCount: 6,
		}, cache)
	}()

	select {
	case <-ctx.Done():
		err := apperrors.Internal("timeout waiting for coupons to load", nil)
		logger.Error(context.Background(), "Promo cache not ready", "error", err)
		return &Validator{cache: cache}, err
	case err := <-errCh:
		if err != nil {
			logger.Error(context.Background(), "failed to load promo coupons", "error", err)
			return &Validator{cache: cache}, err
		}
	}

	return &Validator{cache: cache}, nil
}

func (v *Validator) Validate(code string) error {
	if len(code) < 8 || len(code) > 10 {
		err := apperrors.BadRequest(fmt.Sprintf("invalid coupon code format: %s", code), nil)
		logger.Warn(context.Background(), "Coupon validation failed: invalid length", "code", code, "error", err)
		return err
	}

	cp, ok := v.cache.Get(code)
	if !ok {
		err := apperrors.BadRequest(fmt.Sprintf("invalid coupon code: %s", code), nil)
		logger.Warn(context.Background(), "Coupon validation failed: not found in cache", "code", code, "error", err)
		return err
	}

	if cp.FileCount < 2 {
		err := apperrors.BadRequest(fmt.Sprintf("invalid coupon code (not in enough files): %s", code), nil)
		logger.Warn(context.Background(), "Coupon validation failed: insufficient file count", "code", code, "error", err)
		return err
	}

	logger.Debug(context.Background(), "Coupon validated successfully", "code", code)
	return nil
}
