package session

import (
	"context"
	"net/http"
	"time"
)

type Session struct {
	ID       string
	UserID   string
	CreateAt time.Time
	IsActive bool
}

type Repository interface {
	Create(ctx context.Context, usrID string) (Session, error)
	FindByID(ctx context.Context, id string) (Session, error)
	Disable(ctx context.Context, id string) error
}

type SessionID string

const (
	CookieKey            = "session"
	ContextKey SessionID = "session"
)

func GetSessionFromCookie(r *http.Request) (string, error) {
	s, err := r.Cookie(CookieKey)
	if err != nil {
		return "", err
	}
	return s.Value, nil
}

func PutSessionIDInCookie(w http.ResponseWriter, sID string) {
	cookie := &http.Cookie{
		Name:     CookieKey,
		Value:    sID,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func GetSessionIDFromRequest(r *http.Request) string {
	s, exist := r.Context().Value(ContextKey).(string)
	if !exist {
		return ""
	}
	return s
}

func PutSessionIDInRequest(r *http.Request, sID string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), ContextKey, sID))
}
