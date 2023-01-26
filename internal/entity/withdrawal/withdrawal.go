package withdrawal

import (
	"context"
	"os/user"
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
	Create(ctx context.Context, id string, sum int, usr user.User) error
}
