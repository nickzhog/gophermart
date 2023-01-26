package order

import (
	"context"
	"time"
)

const (
	StatusNew         = "NEW"
	StatusProccessing = "PROCESSING"
	StatusInvalid     = "INVALID"
	StatusProcessed   = "PROCESSED"
)

type Order struct {
	ID       string `json:"id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	Status   string `json:"status,omitempty"`
	Accrual  string `json:"accrual,omitempty"`
	Sum      string `json:"sum,omitempty"`
	SumFloat float64
	UploadAt time.Time `json:"upload_at,omitempty"`
}

type Repository interface {
	Create(ctx context.Context) (Order, error)
	FindForUser(ctx context.Context, usrID string) ([]Order, error)
	Update(ctx context.Context, o Order) error
}
