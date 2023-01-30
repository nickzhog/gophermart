package web

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/nickzhog/gophermart/internal/entity/order"
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
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	var authData LoginPassword
	err = json.Unmarshal(body, &authData)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	usr, err := user.NewUser(authData.Login, authData.Password)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = h.User.FindByLogin(r.Context(), authData.Login)
	if err == nil {
		writeError(w, "login already used", http.StatusConflict)
		return
	}

	err = h.User.Create(r.Context(), &usr)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s := session.GetSessionFromRequest(r)

	err = h.SessionAccount.Create(r.Context(), usr.ID, s.ID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("regiteration complete"))
}

// аутентификация пользователя
func (h *HandlerData) loginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	var authData LoginPassword
	err = json.Unmarshal(body, &authData)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	usr, err := h.User.FindByLogin(r.Context(), authData.Login)
	if err != nil {
		writeError(w, "user not found", http.StatusUnauthorized)
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(authData.Password))
	if err != nil {
		writeError(w, "wrong password", http.StatusUnauthorized)
		return
	}

	s := session.GetSessionFromRequest(r)

	err = h.SessionAccount.Create(r.Context(), usr.ID, s.ID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeAnswer(w, "authentication complete", http.StatusAccepted)
}

// загрузка пользователем номера заказа для расчёта
func (h *HandlerData) newOrderHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	usr := user.GetUserFromRequest(r)
	order, err := order.NewOrder(string(body), usr.ID)
	if err != nil {
		writeError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if o, err := h.Order.FindByID(r.Context(), string(body)); err == nil {
		if o.UserID == usr.ID {
			writeAnswer(w, "already have that order", http.StatusOK)
			return
		} else {
			writeError(w, "order for another user", http.StatusConflict)
			return
		}
	}

	err = h.Order.Create(r.Context(), &order)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeAnswer(w, "success", http.StatusAccepted)
}

// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
func (h *HandlerData) getOrdersHandler(w http.ResponseWriter, r *http.Request) {

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
