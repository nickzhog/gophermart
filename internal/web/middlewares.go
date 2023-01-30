package web

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/nickzhog/gophermart/internal/entity/user"
	"github.com/nickzhog/gophermart/internal/web/session"
)

func (h *HandlerData) createSession(w http.ResponseWriter, r *http.Request) (session.Session, error) {
	s, err := h.Session.Create(r.Context(), r.UserAgent(), r.RemoteAddr)
	if err != nil {
		return session.Session{}, err
	}
	session.PutSessionIDInCookie(w, s.ID)

	return s, nil
}

func (h *HandlerData) HandleSession(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var s session.Session
		var err error

		s, err = session.GetSessionFromCookie(r, h.Session)
		if err != nil {
			s, err = h.createSession(w, r)
			if err != nil {
				writeError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		r = session.PutSessionInRequest(r, s)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (h *HandlerData) HandleUserFromSession(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		s := session.GetSessionFromRequest(r)

		usrID, err := h.SessionAccount.FindUserForSession(r.Context(), s.ID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		usr, err := h.User.FindByID(r.Context(), usrID)
		if err != nil {
			h.SessionAccount.Disable(r.Context(), s.ID)
			next.ServeHTTP(w, r)
			return
		}

		r = user.PutUserInRequest(r, usr)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (h *HandlerData) RequireAuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !user.IsAuthenticated(r) {
			writeError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (h *HandlerData) RequireNotAuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if user.IsAuthenticated(r) {
			writeError(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

///

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
