package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nickzhog/gophermart/internal/repositories"
	"github.com/nickzhog/gophermart/internal/service/user"
	"github.com/nickzhog/gophermart/internal/service/withdrawal"
	"github.com/nickzhog/gophermart/pkg/logging"
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
		r.Get("/balance", h.balanceHandler)
		r.Post("/balance/withdraw", h.withdrawActionHandler)
		r.Get("/withdrawals", h.withdrawalsHandler)
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

// получение текущего баланса счёта баллов лояльности пользователя
func (h *handler) balanceHandler(w http.ResponseWriter, r *http.Request) {
	usrID := user.GetUserIDFromRequest(r)

	usr, err := h.User.FindByID(r.Context(), usrID)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balance, err := usr.CalculateBalance(r.Context(), h.Order, h.Withdrawal)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	withdrawn, err := usr.CalculateWithdrawn(r.Context(), h.Withdrawal)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := make(map[string]interface{})
	m["current"] = balance
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
func (h *handler) withdrawActionHandler(w http.ResponseWriter, r *http.Request) {
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
	h.logger.Tracef("withdrawal request: %+v", wReq)

	usrID := user.GetUserIDFromRequest(r)
	usr, err := h.User.FindByID(r.Context(), usrID)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balance, err := usr.CalculateBalance(r.Context(), h.Order, h.Withdrawal)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	h.logger.Tracef("new withdrawal: %+v", wdl)

	h.writeAnswer(w, "withdrawal succeeded", http.StatusOK)
}

// получение информации о выводе средств с накопительного счёта пользователем
func (h *handler) withdrawalsHandler(w http.ResponseWriter, r *http.Request) {
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
