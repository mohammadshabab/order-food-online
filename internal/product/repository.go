package product

import "context"

//go:generate mockgen -source=repository.go -destination=mock_repository.go -package=product
type Repository interface {
	List(ctx context.Context) ([]*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
}
