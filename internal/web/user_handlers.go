package web

import "net/http"

// регистрация пользователя
func (h *HandlerData) registerHandler(w http.ResponseWriter, r *http.Request) {

}

// аутентификация пользователя
func (h *HandlerData) loginHandler(w http.ResponseWriter, r *http.Request) {

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
