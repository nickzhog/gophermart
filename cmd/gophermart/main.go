package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/nickzhog/gophermart/internal/config"
	"github.com/nickzhog/gophermart/internal/migration"
	"github.com/nickzhog/gophermart/internal/orderprocesser"
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

	err := migration.Migrate(cfg.Settings.DatabaseURI)
	if err != nil {
		logger.Fatal(err)
	}

	reps := repositories.GetRepositories(ctx, logger, cfg)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		err := orderprocesser.NewProcesser(logger, cfg, reps.Order).StartScan(ctx)
		if err != nil {
			logger.Errorf("order processer error: %s", err.Error())
		}
		wg.Done()
	}()

	go func() {
		srv := web.PrepareServer(logger, cfg, reps)
		if err := web.Serve(ctx, logger, srv); err != nil {
			logger.Errorf("failed to serve: %s", err.Error())
		}
		wg.Done()
	}()
	wg.Wait()
}
