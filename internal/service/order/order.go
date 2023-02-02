package order

import (
	"context"
	"errors"
	"strconv"
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
	ID           string    `json:"number"`
	UserID       string    `json:"-"`
	Status       string    `json:"status"`
	Accrual      string    `json:"-"`
	AccrualFloat float64   `json:"accrual"`
	UploadAt     time.Time `json:"uploaded_at"`
}

type Repository interface {
	Create(ctx context.Context, o *Order) error
	FindByID(ctx context.Context, id string) (Order, error)
	FindForUser(ctx context.Context, usrID string) ([]Order, error)
	Update(ctx context.Context, o *Order) error
}

func NewOrder(id, usrID string) (Order, error) {
	if len(id) < 1 || len(usrID) < 1 {
		return Order{}, errors.New("wrong data")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return Order{}, err
	}
	if !ValidLuhn(idInt) {
		return Order{}, errors.New("luhn check fail")
	}

	o := Order{
		ID:      id,
		UserID:  usrID,
		Accrual: "0.0",
	}
	return o, nil
}

func AccrualSumForOrders(ords []Order) float64 {
	ans := 0.0
	for _, o := range ords {
		if o.Status != StatusProcessed {
			continue
		}
		ans += o.AccrualFloat
	}
	return ans
}