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
	q := `
	CREATE TABLE IF NOT EXIST public.users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		login TEXT NOT NULL,
		password_hash TEXT NOT NULL, 
		balance TEXT NOT NULL,
		withdrawn_amount TEXT NOT NULL
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
