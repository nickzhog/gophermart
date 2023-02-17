package session

import (
	"context"
	"net/http"
	"time"

	"github.com/nickzhog/gophermart/internal/service/user"
)

type Session struct {
	ID       string
	UserID   string
	CreateAt time.Time
	IsActive bool
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
	cookie := http.Cookie{
		Name:     CookieKey,
		Value:    sID,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func GetSessionIDFromRequest(r *http.Request) string {
	sID := r.Context().Value(ContextKey).(string)
	if len(sID) < 1 {
		panic("sessionID is empty")
	}
	return sID
}

func PutSessionDataInRequest(r *http.Request, sID, usrID string) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), ContextKey, sID))
	return r.WithContext(context.WithValue(r.Context(), user.ContextKey, usrID))
}
