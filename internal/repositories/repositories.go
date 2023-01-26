package repositories

import (
	"context"
	"time"

	"github.com/nickzhog/gophermart/internal/config"
	"github.com/nickzhog/gophermart/internal/entity/order"
	"github.com/nickzhog/gophermart/internal/entity/user"
	"github.com/nickzhog/gophermart/internal/entity/withdrawal"
	"github.com/nickzhog/gophermart/internal/postgres"
	"github.com/nickzhog/gophermart/pkg/logging"
)

const (
	maxAttempts      = 3
	dbConnectTimeOut = time.Second * 5
)

type Repositories struct {
	Order      order.Repository
	User       user.Repository
	Withdrawal withdrawal.Repository
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
	return Repositories{
		Order:      order.NewRepository(pool, logger),
		User:       user.NewRepository(pool, logger),
		Withdrawal: withdrawal.NewRepository(pool, logger),
	}
}
