package order

import "context"

//go:generate mockgen -source=repository.go -destination=mock_repository.go -package=order
type Repository interface {
	Create(ctx context.Context, order *Order) (*Order, error)
}
