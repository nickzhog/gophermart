package user

import (
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

type UserID string

const ContextKey UserID = "user"

func NewUser(login, password string) (User, error) {
	if len(login) < 1 || len(password) < 1 {
		return User{}, errors.New("login or password is empty")
	}

	phash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	return User{Login: login, PasswordHash: string(phash)}, nil
}

func GetUserIDFromRequest(r *http.Request) string {
	usrID := r.Context().Value(ContextKey).(string)
	if len(usrID) < 1 {
		panic("usrID is empty")
	}
	return usrID
}

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
