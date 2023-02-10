package repositories

import (
	"context"
	"time"

	"github.com/nickzhog/gophermart/internal/config"
	"github.com/nickzhog/gophermart/internal/service/order"
	orderdb "github.com/nickzhog/gophermart/internal/service/order/db"
	"github.com/nickzhog/gophermart/internal/service/user"
	userdb "github.com/nickzhog/gophermart/internal/service/user/db"
	"github.com/nickzhog/gophermart/internal/service/withdrawal"
	withdrawaldb "github.com/nickzhog/gophermart/internal/service/withdrawal/db"
	"github.com/nickzhog/gophermart/internal/web/session"
	"github.com/nickzhog/gophermart/pkg/logging"
	"github.com/nickzhog/gophermart/pkg/postgres"
)

const (
	maxAttempts      = 3
	dbConnectTimeOut = time.Second * 5
)

type Repositories struct {
	User       user.Repository
	Order      order.Repository
	Withdrawal withdrawal.Repository
	Session    session.Repository
}

func GetRepositories(ctx context.Context, logger *logging.Logger, cfg *config.Config) Repositories {
	ctx, cancel := context.WithTimeout(ctx, dbConnectTimeOut)
	defer cancel()

	pool, err := postgres.NewConnection(ctx, maxAttempts, cfg.Settings.DatabaseURI)
	if err != nil {
		logger.Fatal(err)
	}
	if err = pool.Ping(ctx); err != nil {
		logger.Fatal(err)
	}
	return Repositories{
		User:       userdb.NewRepository(pool, logger),
		Order:      orderdb.NewRepository(pool, logger),
		Withdrawal: withdrawaldb.NewRepository(pool, logger),
		Session:    session.NewRepository(pool, logger),
	}
}
