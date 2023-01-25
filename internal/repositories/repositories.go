package repositories

import (
	"context"
	"time"

	"github.com/nickzhog/gophermart/internal/config"
	"github.com/nickzhog/gophermart/internal/postgres"
	"github.com/nickzhog/gophermart/pkg/logging"
)

const (
	maxAttempts      = 3
	dbConnectTimeOut = time.Second * 5
)

type Repositories struct {
}

func GetRepositories(logger *logging.Logger, cfg *config.Config) Repositories {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeOut)
	defer cancel()

	pool, err := postgres.NewConnection(ctx, maxAttempts, cfg.Settings.DatabaseURI)
	if err != nil {
		logger.Fatal(err)
	}
	if err = pool.Ping(ctx); err != nil {
		logger.Fatal(err)
	}
	return Repositories{}
}
