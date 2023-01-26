package order

import (
	"context"

	"github.com/nickzhog/gophermart/internal/postgres"
	"github.com/nickzhog/gophermart/pkg/logging"
)

type repository struct {
	client postgres.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context) (Order, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) FindForUser(ctx context.Context, usrID string) ([]Order, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) Update(ctx context.Context, o Order) error {
	panic("not implemented") // TODO: Implement
}

func NewRepository(client postgres.Client, logger *logging.Logger) Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
