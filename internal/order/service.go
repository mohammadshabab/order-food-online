package order

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mock_service.go -package=order
type Service interface {
	CreateOrder(ctx context.Context, req *OrderReq) (*Order, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateOrder(ctx context.Context, req *OrderReq) (*Order, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	order := &Order{
		ID:         uuid.New().String(),
		Items:      *req.Items,
		CouponCode: req.CouponCode,
	}

	return s.repo.Create(ctx, order)
}
