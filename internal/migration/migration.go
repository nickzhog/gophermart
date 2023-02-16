package migration

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/nickzhog/gophermart/pkg/logging"
)

//go:embed migrations/*
var migrations embed.FS

func Migrate(logger *logging.Logger, connString string) {
	src, err := iofs.New(migrations, "migrations")
	if err != nil {
		logger.Fatal(err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", src, connString)
	if err != nil {
		logger.Fatal(err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal(err)
	}
}
