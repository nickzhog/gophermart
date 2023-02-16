package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4"
	"github.com/nickzhog/gophermart/internal/repositories"
	"github.com/nickzhog/gophermart/internal/service/user"
	"github.com/nickzhog/gophermart/internal/web/session"
	"github.com/nickzhog/gophermart/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	logger *logging.Logger
	repositories.Repositories
}

func NewHandler(logger *logging.Logger, reps repositories.Repositories) *handler {
	return &handler{
		logger:       logger,
		Repositories: reps,
	}
}

func (h *handler) GetRouteGroup() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/register", h.registerHandler)
		r.Post("/login", h.loginHandler)
	}
}

func (h *handler) writeError(w http.ResponseWriter, err string, code int) {
	h.logger.Tracef("handler return error: %s, code: %v", err, code)
	http.Error(w, err, code)
}

func (h *handler) writeAnswer(w http.ResponseWriter, ans string, code int) {
	h.logger.Tracef("answer: %s, code: %v", ans, code)

	w.WriteHeader(code)
	w.Write([]byte(ans))
}

// регистрация пользователя
func (h *handler) registerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	authData, err := user.ParseAuthRequest(body)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	usr, err := user.NewUser(authData.Login, authData.Password)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = h.User.FindByLogin(r.Context(), authData.Login)
	if !errors.Is(err, pgx.ErrNoRows) {
		if err == nil {
			h.writeError(w, "login already used", http.StatusConflict)
			return
		}
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.User.Create(r.Context(), &usr)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sID, err := session.GetSessionFromCookie(r)
	if err == nil {
		if err = h.Session.Disable(r.Context(), sID); err != nil {
			h.logger.Errorf("cant disable old sessions: %s", err.Error())
			h.writeError(w, "cant disable old sessions", http.StatusInternalServerError)
			return
		}
	}
	s, err := h.Session.Create(r.Context(), usr.ID)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.PutSessionIDInCookie(w, s.ID)
	h.writeAnswer(w, "regiteration complete", http.StatusOK)
}

// аутентификация пользователя
func (h *handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	authData, err := user.ParseAuthRequest(body)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	usr, err := h.User.FindByLogin(r.Context(), authData.Login)
	if err != nil {
		h.writeError(w, "user not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(authData.Password))
	if err != nil {
		h.writeError(w, "wrong password", http.StatusUnauthorized)
		return
	}

	sID, err := session.GetSessionFromCookie(r)
	if err == nil {
		h.Session.Disable(r.Context(), sID)
	}
	s, err := h.Session.Create(r.Context(), usr.ID)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.PutSessionIDInCookie(w, s.ID)

	h.writeAnswer(w, "authentication complete", http.StatusOK)
}
