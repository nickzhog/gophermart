package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nickzhog/gophermart/internal/service/order"
	"github.com/nickzhog/gophermart/internal/service/withdrawal"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string `json:"id,omitempty"`
	Login        string `json:"login,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
}

type UserID string

const ContextKey UserID = "user"

var ErrNoRows = errors.New("user not found")

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

func (u *User) CalculateWithdrawn(
	ctx context.Context,
	withdrawalRep withdrawal.Repository,
) (float64, error) {
	withdrawals, err := withdrawalRep.FindForUser(ctx, u.ID)
	if err != nil && err != withdrawal.ErrNoRows {
		return 0, err
	}
	withdrawn := withdrawal.SumForWithdrawals(withdrawals)
	return withdrawn, nil
}

func (u *User) CalculateBalance(
	ctx context.Context,
	orderRep order.Repository,
	withdrawalRep withdrawal.Repository) (float64, error) {

	orders, err := orderRep.FindForUser(ctx, u.ID)
	if err != nil && err != order.ErrNoRows {
		return 0, err
	}
	withdrawn, err := u.CalculateWithdrawn(ctx, withdrawalRep)
	if err != nil {
		return 0, err
	}

	balance := order.AccrualSumForProcessedOrders(orders) - withdrawn
	return balance, nil
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
