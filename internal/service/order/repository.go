package order

import "context"

type Repository interface {
	Create(ctx context.Context, o *Order) error
	FindByID(ctx context.Context, id string) (Order, error)
	FindForUser(ctx context.Context, usrID string) ([]Order, error)
	FindForScanner(ctx context.Context) ([]Order, error)
	Update(ctx context.Context, o *Order) error
}
