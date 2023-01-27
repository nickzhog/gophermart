package withdrawal

import (
	"context"
	"time"
)

type Withdrawal struct {
	ID          string `json:"id,omitempty"`
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
