package user

import "context"

type Repository interface {
	Create(ctx context.Context, usr *User) error
	FindByLogin(ctx context.Context, login string) (User, error)
	FindByID(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, usr *User) error
}
