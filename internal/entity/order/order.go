package order

import (
	"context"
	"strconv"
	"time"
)

const (
	StatusNew         = "NEW"
	StatusProccessing = "PROCESSING"
	StatusInvalid     = "INVALID"
	StatusProcessed   = "PROCESSED"
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
	Create(ctx context.Context) (Order, error)
	FindForUser(ctx context.Context, usrID string) ([]Order, error)
	Update(ctx context.Context, o Order) error
}

func (o *Order) ParseFloats() (err error) {
	if o.AccrualFloat, err = strconv.ParseFloat(o.Accrual, 64); err != nil {
		return
	}
	o.SumFloat, err = strconv.ParseFloat(o.Sum, 64)
	return
}
