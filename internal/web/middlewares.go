package web

import (
	"net/http"

	"github.com/nickzhog/gophermart/internal/web/session"
)

func (h *HandlerData) createSession(r *http.Request) (session.Session, error) {
	s, err := h.Session.Create(r.Context(), r.UserAgent(), r.RemoteAddr)
	if err != nil {
		return session.Session{}, err
	}
	return s, nil
}

func (h *HandlerData) HandleSession(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var s session.Session
		var err error

		sID, err := session.GetSessionIDFromCookie(r)
		if err != nil {
			s, err = h.createSession(r)
			if err != nil {
				showError(w, err.Error(), http.StatusBadGateway)
				return
			}

			session.PutSessionIDInCookie(w, s.ID)
		} else {
			s, err = h.Session.FindByID(r.Context(), sID)
			if err != nil {
				s, err = h.createSession(r)
				if err != nil {
					showError(w, err.Error(), http.StatusBadGateway)
					return
				}
				session.PutSessionIDInCookie(w, s.ID)
			}
		}

		session.PutSessionInRequest(r, s)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (h *HandlerData) RequireAuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// todo
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (h *HandlerData) RequireNotAuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// todo
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
