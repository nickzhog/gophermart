package web

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/nickzhog/gophermart/internal/entity/user"
	"github.com/nickzhog/gophermart/internal/web/session"
	"golang.org/x/crypto/bcrypt"
)

type LoginPassword struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

// регистрация пользователя
func (h *HandlerData) registerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		showError(w, err.Error(), http.StatusBadRequest)
		return
	}
	var authData LoginPassword
	err = json.Unmarshal(body, &authData)
	if err != nil {
		showError(w, err.Error(), http.StatusBadRequest)
		return
	}
	usr, err := user.NewUser(authData.Login, authData.Password)
	if err != nil {
		showError(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = h.User.FindByLogin(r.Context(), authData.Login)
	if err == nil {
		showError(w, "login already used", http.StatusConflict)
		return
	}

	err = h.User.Create(r.Context(), &usr)
	if err != nil {
		showError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s, _ := session.GetSessionFromRequest(r)

	err = h.SessionAccount.Create(r.Context(), usr.ID, s.ID)
	if err != nil {
		showError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("regiteration complete"))
}

// аутентификация пользователя
func (h *HandlerData) loginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		showError(w, err.Error(), http.StatusBadRequest)
		return
	}
	var authData LoginPassword
	err = json.Unmarshal(body, &authData)
	if err != nil {
		showError(w, err.Error(), http.StatusBadRequest)
		return
	}
	usr, err := h.User.FindByLogin(r.Context(), authData.Login)
	if err != nil {
		showError(w, "user not found", http.StatusUnauthorized)
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(authData.Password))
	if err != nil {
		showError(w, "wrong password", http.StatusUnauthorized)
		return
	}

	s, _ := session.GetSessionFromRequest(r)

	err = h.SessionAccount.Create(r.Context(), usr.ID, s.ID)
	if err != nil {
		showError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("authentication complete"))
}

// загрузка пользователем номера заказа для расчёта
func (h *HandlerData) ordersActionHandler(w http.ResponseWriter, r *http.Request) {

}

// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
func (h *HandlerData) ordersHandler(w http.ResponseWriter, r *http.Request) {

}

// получение текущего баланса счёта баллов лояльности пользователя
func (h *HandlerData) balanceHandler(w http.ResponseWriter, r *http.Request) {

}

// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
func (h *HandlerData) withdrawActionHandler(w http.ResponseWriter, r *http.Request) {

}

// получение информации о выводе средств с накопительного счёта пользователем
func (h *HandlerData) withdrawalsHandler(w http.ResponseWriter, r *http.Request) {

}
