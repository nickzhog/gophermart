package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/nickzhog/gophermart/internal/postgres"
	"github.com/nickzhog/gophermart/pkg/logging"
)

type repository struct {
	client postgres.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, userAgent, ip string) (Session, error) {
	q := `
	INSERT INTO public.sessions 
		(user_agent, ip) 
	VALUES 
		($1, $2) 
	RETURNING id, create_at, is_active
	`
	s := Session{
		UserAgent: userAgent,
		IP:        ip,
	}
	err := r.client.QueryRow(ctx, q, s.UserAgent, s.IP).
		Scan(&s.ID, &s.CreateAt, &s.IsActive)
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
		id, create_at, useragent, ip, is_active
	FROM 
		public.session 
	WHERE 
		id = $1 and is_active = true
	`
	var s Session
	err := r.client.QueryRow(ctx, q, id).
		Scan(&s.ID, &s.CreateAt, &s.UserAgent, &s.IP, &s.IsActive)

	if err != nil {
		return Session{}, err
	}

	return s, nil
}

func NewRepository(client postgres.Client, logger *logging.Logger) Repository {
	q := `
	CREATE TABLE IF NOT EXISTS public.sessions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		user_agent TEXT NOT NULL,
		ip TEXT NOT NULL,
		is_active bool NOT NULL DEFAULT true
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
