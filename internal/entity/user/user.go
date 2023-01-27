package user

import (
	"context"
	"net/http"
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

type UserKey string

const ContextKey UserKey = "user"

func GetUserFromRequest(r *http.Request) (User, bool) {
	s, exist := r.Context().Value(ContextKey).(User)
	if !exist {
		return User{}, false
	}
	return s, true
}

func PutUserInRequest(r *http.Request, usr User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), ContextKey, usr))
}
