package withdrawal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Withdrawal struct {
	ID          string    `json:"order,omitempty"`
	UserID      string    `json:"-"`
	Sum         string    `json:"-"`
	SumFloat    float64   `json:"sum,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}

type Repository interface {
	Create(ctx context.Context, w *Withdrawal) error
	FindForUser(ctx context.Context, usrID string) ([]Withdrawal, error)
	FindByID(ctx context.Context, id string) (Withdrawal, error)
}

type WithdrawalRequest struct {
	Order string  `json:"order,omitempty"`
	Sum   float64 `json:"sum,omitempty"`
}

func ParseWithdrawalRequest(data []byte) (WithdrawalRequest, error) {
	var wr WithdrawalRequest
	err := json.Unmarshal(data, &wr)
	if err != nil {
		return WithdrawalRequest{}, err
	}
	return wr, nil
}

func NewWithdrawal(orderID, usrID string, sum float64) (Withdrawal, error) {
	w := Withdrawal{
		ID:     orderID,
		UserID: usrID,
		Sum:    fmt.Sprintf("%g", sum),
	}
	idInt, err := strconv.Atoi(orderID)
	if err != nil {
		return Withdrawal{}, err
	}

	if idInt < 1 {
		return Withdrawal{}, errors.New("wrong order")
	}

	return w, nil
}
