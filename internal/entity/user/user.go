package user

import (
	"context"
	"strconv"
)

type User struct {
	ID                   string  `json:"id,omitempty"`
	Login                string  `json:"login,omitempty"`
	PasswordHash         string  `json:"password_hash,omitempty"`
	WithdrawnAmount      string  `json:"withdrawn_amount,omitempty"` // сумма использованных за весь период регистрации баллов
	WithdrawnAmountFloat float64 `json:"withdrawn_amount_float,omitempty"`
	Balance              string  `json:"balance,omitempty"`
	BalanceFloat         float64 `json:"balance_float,omitempty"`
}

type Repository interface {
	Create(ctx context.Context, usr User) (User, error)
	FindByLogin(ctx context.Context, login string) (User, error)
	FindByID(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, usr User) error
}

func (u *User) ParseFloats() (err error) {
	if u.WithdrawnAmountFloat, err = strconv.ParseFloat(u.WithdrawnAmount, 64); err != nil {
		return
	}
	u.BalanceFloat, err = strconv.ParseFloat(u.Balance, 64)
	return
}
