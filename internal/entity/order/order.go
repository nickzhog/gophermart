package order

import (
	"context"
	"errors"
	"time"
)

const (
	StatusNew         = "NEW"        // заказ создан
	StatusInvalid     = "INVALID"    // заказ не принят к расчёту, и вознаграждение не будет начислено
	StatusRegistered  = "REGISTERED" // заказ зарегистрирован, но не начисление не рассчитано
	StatusProccessing = "PROCESSING" // расчёт начисления в процессе
	StatusProcessed   = "PROCESSED"  // расчёт начисления окончен
)

type Order struct {
	ID           string `json:"id,omitempty"`
	UserID       string `json:"user_id,omitempty"`
	Status       string `json:"status,omitempty"`
	Accrual      string `json:"accrual,omitempty"`
	AccrualFloat float64
	Sum          string `json:"sum,omitempty"`
	SumFloat     float64
	UploadAt     time.Time `json:"upload_at,omitempty"`
}

type Repository interface {
	Create(ctx context.Context, o *Order) error
	FindByID(ctx context.Context, id string) (Order, error)
	FindForUser(ctx context.Context, usrID string) ([]Order, error)
	Update(ctx context.Context, o *Order) error
}

func IsValidID(id int64) bool {
	if id < 1 {
		return false
	}

	return true
}

func NewOrder(id, usrID string) (Order, error) {

	if len(id) < 1 || len(usrID) < 1 {
		return Order{}, errors.New("login or password is empty")
	}

	o := Order{
		ID:      id,
		UserID:  usrID,
		Sum:     "0.0",
		Accrual: "0.0",
	}
	return o, nil
}
