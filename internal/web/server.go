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
	orderHandler "github.com/nickzhog/gophermart/internal/service/order/handler"
	userHandler "github.com/nickzhog/gophermart/internal/service/user/handler"
	withdrawalHandler "github.com/nickzhog/gophermart/internal/service/withdrawal/handler"
	"github.com/nickzhog/gophermart/pkg/logging"
)

func PrepareServer(logger *logging.Logger, cfg *config.Config, reps repositories.Repositories) *http.Server {
	orderHandler := orderHandler.NewHandler(logger, reps)
	userHandler := userHandler.NewHandler(logger, reps)
	withdrawalHander := withdrawalHandler.NewHandler(logger, reps)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Use(gzipCompress)
	r.Use(gzipDecompress)

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Group(userHandler.GetRouteGroup())

			r.Group(func(r chi.Router) {
				r.Use(SessionMiddleware(logger, reps))

				r.Group(orderHandler.GetRouteGroup())
				r.Group(withdrawalHander.GetRouteGroup())
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
