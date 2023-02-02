package main

import (
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

	reps := repositories.GetRepositories(logger, cfg)

	go accrual.OrdersScanStart(logger, cfg, reps)
	web.StartServer(logger, cfg, reps)
}
