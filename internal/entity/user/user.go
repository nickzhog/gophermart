package user

import (
	"context"
)

type User struct {
	ID           string  `json:"id,omitempty"`
	Login        string  `json:"login,omitempty"`
	PasswordHash string  `json:"password_hash,omitempty"`
	Balance      string  `json:"balance,omitempty"`
	BalanceFloat float64 `json:"balance_float,omitempty"`
}

type Repository interface {
	Create(ctx context.Context, usr *User) error
	FindByLogin(ctx context.Context, login string) (User, error)
	FindByID(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, usr *User) error
}
