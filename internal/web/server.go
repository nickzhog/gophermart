package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nickzhog/gophermart/internal/config"
	"github.com/nickzhog/gophermart/internal/repositories"
	"github.com/nickzhog/gophermart/pkg/logging"
)

type HandlerData struct {
	Logger *logging.Logger
	repositories.Repositories
}

func StartServer(logger *logging.Logger, cfg *config.Config, reps repositories.Repositories) {
	h := HandlerData{
		Logger:       logger,
		Repositories: reps,
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	// r.Use(h.logMiddleware)

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(h.RequireNotAuthMiddleware)

				r.Post("/register", h.registerHandler)
				r.Post("/login", h.loginHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(h.RequireAuthMiddleware)

				//загрузка пользователем номера заказа для расчёта
				r.Post("/orders", h.ordersActionHandler)

				//получение списка загруженных пользователем номеров заказов
				r.Get("/orders", h.ordersHandler)

				r.Get("/balance", h.balanceHandler)
				r.Post("/balance/withdraw", h.withdrawActionHandler)
				r.Get("/withdrawals", h.withdrawalsHandler)
			})
		})
	})

	logger.Fatal(http.ListenAndServe(cfg.Settings.RunAddress, r))
}
