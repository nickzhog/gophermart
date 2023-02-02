package sessionaccount

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

func (r *repository) Create(ctx context.Context, usrID, sID string) error {
	q := `
	INSERT INTO public.session_user 
		(session_id, user_id) 
	VALUES 
		($1, $2) 
	`
	_, err := r.client.Exec(ctx, q, usrID, sID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error("err:", newErr.Error())
		}
		r.logger.Error("err:", err.Error())
		return err
	}
	return nil
}

func (r *repository) FindUserForSession(ctx context.Context, sessionID string) (usrID string, err error) {
	q := `
	SELECT
		user_id
	FROM 
		public.session_user 
	WHERE 
	session_id = $1 and is_active = true
	`

	err = r.client.QueryRow(ctx, q, sessionID).Scan(&usrID)
	if err != nil {
		return "", err
	}

	return
}

func (r *repository) Disable(ctx context.Context, sessionID string) {
	q := `
		UPDATE 
			public.session_user 
		SET
			is_active = false
		WHERE 
			session_id = $1 and is_active = true
	`

	_, err := r.client.Exec(ctx, q, sessionID)
	if err != nil {
		r.logger.Error(err)
	}
}

func NewRepository(client postgres.Client, logger *logging.Logger) Repository {
	q := `
	CREATE TABLE IF NOT EXISTS public.session_user (
		session_id UUID NOT NULL,
		user_id UUID NOT NULL,
		loginned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		is_active bool NOT NULL DEFAULT true,
		constraint session_id FOREIGN KEY (session_id) REFERENCES public.sessions (id),
		constraint user_id FOREIGN KEY (user_id) REFERENCES public.users (id)
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
