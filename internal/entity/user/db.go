package user

import (
	"context"

	"github.com/nickzhog/gophermart/internal/postgres"
	"github.com/nickzhog/gophermart/pkg/logging"
)

type repository struct {
	client postgres.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, usr User) (User, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) FindByLogin(ctx context.Context, login string) (User, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) FindByID(ctx context.Context, id string) (User, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) Update(ctx context.Context, usr User) error {
	panic("not implemented") // TODO: Implement
}

func NewRepository(client postgres.Client, logger *logging.Logger) Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
