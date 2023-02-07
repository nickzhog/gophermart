package web

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/nickzhog/gophermart/internal/service/order"
	"github.com/nickzhog/gophermart/internal/service/user"
	"github.com/nickzhog/gophermart/internal/service/withdrawal"
	"github.com/nickzhog/gophermart/internal/web/session"
	"golang.org/x/crypto/bcrypt"
)

// регистрация пользователя
func (h *HandlerData) registerHandler(w http.ResponseWriter, r *http.Request) {
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
	if err == nil {
		h.writeError(w, "login already used", http.StatusConflict)
		return
	}

	err = h.User.Create(r.Context(), &usr)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
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
	h.writeAnswer(w, "regiteration complete", http.StatusOK)
}

// аутентификация пользователя
func (h *HandlerData) loginHandler(w http.ResponseWriter, r *http.Request) {
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

// загрузка пользователем номера заказа для расчёта
func (h *HandlerData) newOrderHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	usrID := user.GetUserIDFromRequest(r)
	order, err := order.NewOrder(string(body), usrID)
	if err != nil {
		h.Logger.Errorf("bad order: %s, %s", string(body), err.Error())
		h.writeError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if o, err := h.Order.FindByID(r.Context(), string(body)); err == nil {
		if o.UserID == usrID {
			h.writeAnswer(w, "already have that order", http.StatusOK)
			return
		} else {
			h.writeError(w, "order for another user", http.StatusConflict)
			return
		}
	}

	err = h.Order.Create(r.Context(), &order)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Tracef("new order: %+v", order)
	h.writeAnswer(w, "success", http.StatusAccepted)
}

// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
func (h *HandlerData) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	usrID := user.GetUserIDFromRequest(r)
	orders, err := h.Order.FindForUser(r.Context(), usrID)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(orders) < 1 {
		h.writeAnswer(w, "no orders", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(orders)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeAnswer(w, string(data), http.StatusOK)
}

// получение текущего баланса счёта баллов лояльности пользователя
func (h *HandlerData) balanceHandler(w http.ResponseWriter, r *http.Request) {
	usrID := user.GetUserIDFromRequest(r)
	withdrawals, err := h.Withdrawal.FindForUser(r.Context(), usrID)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	withdrawn := withdrawal.SumForWithdrawals(withdrawals)

	orders, _ := h.Order.FindForUser(r.Context(), usrID)

	m := make(map[string]interface{})
	m["current"] = order.AccrualSumForOrders(orders) - withdrawn
	m["withdrawn"] = withdrawn

	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(m)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeAnswer(w, string(data), http.StatusOK)
}

// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
func (h *HandlerData) withdrawActionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	wReq, err := withdrawal.ParseWithdrawalRequest(body)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	h.Logger.Tracef("withdrawal request: %+v", wReq)
	usrID := user.GetUserIDFromRequest(r)

	withdrawals, _ := h.Withdrawal.FindForUser(r.Context(), usrID)
	orders, _ := h.Order.FindForUser(r.Context(), usrID)

	withdrawn := withdrawal.SumForWithdrawals(withdrawals)
	balance := order.AccrualSumForOrders(orders) - withdrawn
	if balance < wReq.Sum {
		h.writeError(w, "not enough balance", http.StatusPaymentRequired)
		return
	}

	wdl, err := withdrawal.NewWithdrawal(wReq.Order, usrID, wReq.Sum)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	_, err = h.Withdrawal.FindByID(r.Context(), wdl.ID)
	if err == nil {
		h.writeError(w, "order already used", http.StatusConflict)
		return
	}

	err = h.Withdrawal.Create(r.Context(), &wdl)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Tracef("new withdrawal: %+v", wdl)

	h.writeAnswer(w, "withdrawal succeeded", http.StatusOK)
}

// получение информации о выводе средств с накопительного счёта пользователем
func (h *HandlerData) withdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	usrID := user.GetUserIDFromRequest(r)
	withdrawals, err := h.Withdrawal.FindForUser(r.Context(), usrID)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(withdrawals) < 1 {
		h.writeAnswer(w, "no orders", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(withdrawals)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeAnswer(w, string(data), http.StatusOK)
}
