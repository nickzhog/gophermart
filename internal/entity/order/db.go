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
	q := `
	CREATE TABLE IF NOT EXIST public.orders (
		id TEXT PRIMARY KEY,
		user_id UUID NOT NULL,
		status text NOT NULL DEFAULT 'NEW',
		accrual TEXT NOT NULL,
		sum TEXT NOT NULL,
		upload_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		constraint user_id foreign key (user_id) references public.users (id)
	);
	`
	_, err := client.Exec(context.TODO(), q)
	if err != nil {
		logger.Fatal(err)
	}

	return &repository{
		client: client,
		logger: logger,
	}
}
