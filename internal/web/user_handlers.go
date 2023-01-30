package web

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/nickzhog/gophermart/internal/entity/order"
	"github.com/nickzhog/gophermart/internal/entity/user"
	"github.com/nickzhog/gophermart/internal/entity/withdrawal"
	"github.com/nickzhog/gophermart/internal/web/session"
	"golang.org/x/crypto/bcrypt"
)

// регистрация пользователя
func (h *HandlerData) registerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	authData, err := user.ParseAuthRequest(body)
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
	writeAnswer(w, "regiteration complete", http.StatusOK)
}

// аутентификация пользователя
func (h *HandlerData) loginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	authData, err := user.ParseAuthRequest(body)
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
	usr := user.GetUserFromRequest(r)
	orders, err := h.Order.FindForUser(r.Context(), usr.ID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(orders) < 1 {
		writeAnswer(w, "no orders", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// получение текущего баланса счёта баллов лояльности пользователя
func (h *HandlerData) balanceHandler(w http.ResponseWriter, r *http.Request) {
	usr := user.GetUserFromRequest(r)
	withdrawals, err := h.Withdrawal.FindForUser(r.Context(), usr.ID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var withdrawn float64
	for _, v := range withdrawals {
		withdrawn += v.SumFloat
	}
	m := make(map[string]interface{})
	m["current"] = usr.BalanceFloat
	m["withdrawn"] = withdrawn

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
func (h *HandlerData) withdrawActionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	wReq, err := withdrawal.ParseWithdrawalRequest(body)
	if err != nil {
		writeError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	usr := user.GetUserFromRequest(r)
	if usr.BalanceFloat < wReq.Sum {
		writeError(w, "not enough balance", http.StatusPaymentRequired)
		return
	}
	wdl, err := withdrawal.NewWithdrawal(wReq.Order, usr.ID, wReq.Sum)
	if err != nil {
		writeError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	_, err = h.Withdrawal.FindByID(r.Context(), wdl.ID)
	if err == nil {
		writeError(w, "order already used", http.StatusConflict)
		return
	}

	err = h.Withdrawal.Create(r.Context(), &wdl)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeAnswer(w, "withdrawal succeeded", http.StatusOK)
}

// получение информации о выводе средств с накопительного счёта пользователем
func (h *HandlerData) withdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	usr := user.GetUserFromRequest(r)
	withdrawals, err := h.Withdrawal.FindForUser(r.Context(), usr.ID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(withdrawals) < 1 {
		writeAnswer(w, "no orders", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(withdrawals)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
