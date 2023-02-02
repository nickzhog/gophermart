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
	Cfg    *config.Config
	repositories.Repositories
}

func StartServer(logger *logging.Logger, cfg *config.Config, reps repositories.Repositories) {
	h := HandlerData{
		Logger:       logger,
		Repositories: reps,
		Cfg:          cfg,
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Use(GzipMiddleWare)

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Post("/register", h.registerHandler)
				r.Post("/login", h.loginHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(h.SessionMiddleware)

				r.Post("/orders", h.newOrderHandler)
				r.Get("/orders", h.getOrdersHandler)

				r.Get("/balance", h.balanceHandler)
				r.Post("/balance/withdraw", h.withdrawActionHandler)
				r.Get("/withdrawals", h.withdrawalsHandler)
			})
		})
	})

	logger.Trace("starting web server")
	logger.Fatal(http.ListenAndServe(cfg.Settings.RunAddress, r))
}
