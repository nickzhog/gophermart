package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string `json:"id,omitempty"`
	Login        string `json:"login,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
}

type Repository interface {
	Create(ctx context.Context, usr *User) error
	FindByLogin(ctx context.Context, login string) (User, error)
	FindByID(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, usr *User) error
}

type UserID string

const ContextKey UserID = "user"

type AuthRequest struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func ParseAuthRequest(data []byte) (AuthRequest, error) {
	var authData AuthRequest
	err := json.Unmarshal(data, &authData)
	if err != nil {
		return AuthRequest{}, err
	}
	return authData, nil
}

func NewUser(login, password string) (User, error) {
	if len(login) < 1 || len(password) < 1 {
		return User{}, errors.New("login or password is empty")
	}

	phash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	var usr User
	usr.Login = login
	usr.PasswordHash = string(phash)

	return usr, nil
}

func IsAuthenticated(r *http.Request) bool {
	_, exist := r.Context().Value(ContextKey).(string)

	return exist
}

func GetUserIDFromRequest(r *http.Request) string {
	s, exist := r.Context().Value(ContextKey).(string)
	if !exist {
		return ""
	}
	return s
}

func PutUserIDInRequest(r *http.Request, usrID string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), ContextKey, usrID))
}
