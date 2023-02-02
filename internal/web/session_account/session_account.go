package sessionaccount

import (
	"context"
	"time"
)

type UserForSession struct {
	SessionID  string
	UserID     string
	LoginnedAt time.Time
	IsActive   bool
}

type Repository interface {
	Create(ctx context.Context, usrID, sessionID string) error
	FindUserForSession(ctx context.Context, sessionID string) (usrID string, err error)
	Disable(ctx context.Context, sessionID string)
}
