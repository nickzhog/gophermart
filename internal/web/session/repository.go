package session

import "context"

type Repository interface {
	Create(ctx context.Context, usrID string) (Session, error)
	FindByID(ctx context.Context, id string) (Session, error)
	Disable(ctx context.Context, id string) error
}
