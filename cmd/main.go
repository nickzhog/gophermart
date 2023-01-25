package main

import (
	"github.com/nickzhog/gophermart/internal/config"
	"github.com/nickzhog/gophermart/internal/repositories"
	"github.com/nickzhog/gophermart/internal/web"
	"github.com/nickzhog/gophermart/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	reps := repositories.GetRepositories(logger, cfg)

	web.StartServer(logger, cfg, reps)
}
