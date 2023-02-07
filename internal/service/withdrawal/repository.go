package withdrawal

import "context"

type Repository interface {
	Create(ctx context.Context, w *Withdrawal) error
	FindForUser(ctx context.Context, usrID string) ([]Withdrawal, error)
	FindByID(ctx context.Context, id string) (Withdrawal, error)
}
