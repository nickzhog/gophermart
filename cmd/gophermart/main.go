package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/nickzhog/gophermart/internal/accrual"
	"github.com/nickzhog/gophermart/internal/config"
	"github.com/nickzhog/gophermart/internal/repositories"
	"github.com/nickzhog/gophermart/internal/web"
	"github.com/nickzhog/gophermart/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()
	logger.Tracef("%+v", cfg.Settings)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		oscall := <-c
		logger.Tracef("system call:%+v", oscall)
		cancel()
	}()

	reps := repositories.GetRepositories(ctx, logger, cfg)

	srv := web.PrepareServer(logger, cfg, reps)
	go func() {
		if err := web.Serve(ctx, logger, srv); err != nil {
			logger.Errorf("failed to serve: %s", err.Error())
		}
	}()

	err := accrual.NewScanner(logger, cfg, reps.Order).StartScan(ctx)
	if err != nil {
		logger.Errorf("scanner error: %s", err.Error())
	}
}
