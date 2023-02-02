package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/nickzhog/gophermart/internal/postgres"
	"github.com/nickzhog/gophermart/internal/service/user"
	"github.com/nickzhog/gophermart/pkg/logging"
)

type repository struct {
	client postgres.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, usr *user.User) error {
	q := `
		INSERT INTO public.users 
		    (login, password_hash) 
		VALUES 
		    ($1, $2) 
		RETURNING id
	`
	err := r.client.QueryRow(ctx, q, usr.Login, usr.PasswordHash).
		Scan(&usr.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error("err:", newErr.Error())
		}
		r.logger.Error("err:", err.Error())
	}
	return err
}

func (r *repository) FindByLogin(ctx context.Context, login string) (user.User, error) {
	q := `
	SELECT
	 id, password_hash
	FROM public.users WHERE login = $1
	`

	var usr user.User
	err := r.client.QueryRow(ctx, q, login).
		Scan(&usr.ID, &usr.PasswordHash)
	if err != nil {
		return user.User{}, err
	}

	return usr, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (user.User, error) {
	q := `
	SELECT
	 login, password_hash
	FROM public.users WHERE id = $1
	`

	var usr user.User
	err := r.client.QueryRow(ctx, q, id).
		Scan(&usr.Login, &usr.PasswordHash)
	if err != nil {
		return user.User{}, err
	}

	return usr, nil
}

func (r *repository) Update(ctx context.Context, usr *user.User) error {
	q := `
		UPDATE public.users 
		SET
		 password_hash = $1
		WHERE id = $2
	`

	_, err := r.client.Exec(ctx, q, usr.PasswordHash, usr.ID)

	return err
}

func NewRepository(client postgres.Client, logger *logging.Logger) user.Repository {
	q := `
	CREATE TABLE IF NOT EXISTS public.users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		login TEXT NOT NULL,
		password_hash TEXT NOT NULL
	);
	`
	_, err := client.Exec(context.TODO(), q)
	if err != nil {
		logger.Fatal(err)
	}
	return &repository{
		client: client,
		logger: logger,
	}
}
