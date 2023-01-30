package withdrawal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Withdrawal struct {
	ID          string `json:"order_id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	Sum         string `json:"sum,omitempty"`
	SumFloat    float64
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
	if len(orderID) < 1 {
		return Withdrawal{}, errors.New("empty order_id")
	}
	return w, nil
}
