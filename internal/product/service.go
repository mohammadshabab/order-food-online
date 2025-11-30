package product

import "context"

//go:generate mockgen -source=service.go -destination=mock_service.go -package=product
type Service interface {
	ListProducts(ctx context.Context) ([]*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListProducts(ctx context.Context) ([]*Product, error) {
	return s.repo.List(ctx)
}

func (s *service) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}
