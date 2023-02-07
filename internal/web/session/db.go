package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/nickzhog/gophermart/pkg/logging"
	"github.com/nickzhog/gophermart/pkg/postgres"
)

type repository struct {
	client postgres.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, usrID string) (Session, error) {
	q := `
	INSERT INTO public.sessions 
		(user_id) 
	VALUES 
		($1) 
	RETURNING id, user_id, create_at, is_active
	`
	var s Session
	err := r.client.QueryRow(ctx, q, usrID).
		Scan(&s.ID, &s.UserID, &s.CreateAt, &s.IsActive)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error("err:", newErr.Error())
		}
		r.logger.Error("err:", err.Error())
		return Session{}, err
	}
	return s, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (Session, error) {
	q := `
	SELECT
		id, user_id, create_at, is_active
	FROM 
		public.sessions
	WHERE 
		id = $1 and is_active = true
	`
	var s Session
	err := r.client.QueryRow(ctx, q, id).
		Scan(&s.ID, &s.UserID, &s.CreateAt, &s.IsActive)

	if err != nil {
		return Session{}, err
	}

	return s, nil
}

func (r *repository) Disable(ctx context.Context, id string) error {
	q := `
		UPDATE 
			public.sessions 
		SET
			is_active = false
		WHERE 
			id = $1
	`
	_, err := r.client.Exec(ctx, q, id)
	return err
}

func NewRepository(client postgres.Client, logger *logging.Logger) Repository {
	q := `
	CREATE TABLE IF NOT EXISTS public.sessions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL,
		create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		is_active bool NOT NULL DEFAULT true,
		CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES public.users (id)
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
