package session

import (
	"context"
	"net/http"
	"time"
)

type Session struct {
	ID         string
	CreateAt   time.Time
	UserAgent  string
	IP         string
	IsLoginned bool
	IsActive   bool
}

type Repository interface {
	Create(ctx context.Context, userAgent, ip string) (Session, error)
	FindByID(ctx context.Context, id string) (Session, error)
}

type SessionKey string

const (
	CookieKey             = "session"
	ContextKey SessionKey = "session"
)

func GetSessionFromCookie(r *http.Request, rep Repository) (Session, error) {
	sCookie, err := r.Cookie(CookieKey)
	if err != nil {
		return Session{}, err
	}
	s, err := rep.FindByID(r.Context(), sCookie.Value)
	if err != nil {
		return Session{}, err
	}
	return s, nil
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

func GetSessionFromRequest(r *http.Request) (Session, bool) {
	s, exist := r.Context().Value(ContextKey).(Session)
	if !exist {
		return Session{}, false
	}
	return s, true
}

func PutSessionInRequest(r *http.Request, s Session) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), ContextKey, s))
}
