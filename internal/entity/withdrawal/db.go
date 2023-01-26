package withdrawal

import (
	"context"
	"os/user"

	"github.com/nickzhog/gophermart/internal/postgres"
	"github.com/nickzhog/gophermart/pkg/logging"
)

type repository struct {
	client postgres.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, id string, sum int, usr user.User) error {
	panic("not implemented") // TODO: Implement
}

func NewRepository(client postgres.Client, logger *logging.Logger) Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
