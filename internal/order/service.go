package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohammadshabab/order-food-online/internal/promo"
)

//go:generate mockgen -source=service.go -destination=mock_service.go -package=order
type Service interface {
	CreateOrder(ctx context.Context, req *OrderReq) (*Order, error)
}

type service struct {
	repo  Repository
	promo *promo.Validator
}

func NewService(repo Repository, promoValidator *promo.Validator) Service {
	return &service{repo: repo, promo: promoValidator}
}

func (s *service) CreateOrder(ctx context.Context, req *OrderReq) (*Order, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Validate coupon if provided
	if req.CouponCode != nil && s.promo != nil {
		if appErr := s.promo.Validate(*req.CouponCode); appErr != nil {
			return nil, appErr
		}
	}

	order := &Order{
		ID:         uuid.New().String(),
		Items:      *req.Items,
		CouponCode: req.CouponCode,
	}

	return s.repo.Create(ctx, order)
}
