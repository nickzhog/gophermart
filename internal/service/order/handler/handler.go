package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4"
	"github.com/nickzhog/gophermart/internal/repositories"
	"github.com/nickzhog/gophermart/internal/service/order"
	"github.com/nickzhog/gophermart/internal/service/user"
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
		r.Post("/orders", h.newOrderHandler)
		r.Get("/orders", h.getOrdersHandler)
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

// загрузка пользователем номера заказа для расчёта
func (h *handler) newOrderHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	usrID := user.GetUserIDFromRequest(r)
	order, err := order.NewOrder(string(body), usrID)
	if err != nil {
		h.logger.Errorf("bad order: %s, %s", string(body), err.Error())
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

	h.logger.Tracef("new order: %+v", order)
	h.writeAnswer(w, "success", http.StatusAccepted)
}

// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
func (h *handler) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	usrID := user.GetUserIDFromRequest(r)

	orders, err := h.Order.FindForUser(r.Context(), usrID)
	if err != nil && err != pgx.ErrNoRows {
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
