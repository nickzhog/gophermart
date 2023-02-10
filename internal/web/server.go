package web

import (
	"context"
	"log"
	"net/http"
	"time"

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

func PrepareServer(logger *logging.Logger, cfg *config.Config, reps repositories.Repositories) *http.Server {
	h := HandlerData{
		Logger:       logger,
		Cfg:          cfg,
		Repositories: reps,
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Use(gzipMiddleWare)

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

	return &http.Server{
		Addr:    cfg.Settings.RunAddress,
		Handler: r,
	}
}

func Serve(ctx context.Context, logger *logging.Logger, srv *http.Server) (err error) {
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	logger.Tracef("server started")

	<-ctx.Done()

	logger.Tracef("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		logger.Fatalf("server Shutdown Failed:%+s", err)
	}

	logger.Tracef("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
