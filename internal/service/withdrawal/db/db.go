package db

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgconn"
	"github.com/nickzhog/gophermart/internal/service/withdrawal"
	"github.com/nickzhog/gophermart/pkg/logging"
	"github.com/nickzhog/gophermart/pkg/postgres"
)

type repository struct {
	client postgres.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, w *withdrawal.Withdrawal) error {

	w.Sum = fmt.Sprintf("%g", w.SumFloat)
	q := `
		INSERT INTO public.withdrawals 
		    (id, user_id, sum) 
		VALUES 
		    ($1, $2, $3) 
		RETURNING processed_at
	`
	err := r.client.QueryRow(ctx, q, w.ID, w.UserID, w.Sum).Scan(&w.ProcessedAt)
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

func (r *repository) FindForUser(ctx context.Context, usrID string) ([]withdrawal.Withdrawal, error) {
	q := `
		SELECT 
			id, user_id, sum, processed_at
		FROM 
			public.withdrawals 
		WHERE 
			user_id = $1;
	`

	rows, err := r.client.Query(ctx, q, usrID)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	wdls := make([]withdrawal.Withdrawal, rows.CommandTag().RowsAffected())

	for rows.Next() {
		var w withdrawal.Withdrawal

		err = rows.Scan(&w.ID, &w.UserID, &w.Sum, &w.ProcessedAt)

		if err != nil {
			r.logger.Error(err)
			return nil, err
		}

		if w.SumFloat, err = strconv.ParseFloat(w.Sum, 64); err != nil {
			r.logger.Error(err)
			return nil, err
		}

		wdls = append(wdls, w)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return wdls, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (withdrawal.Withdrawal, error) {
	q := `
	SELECT
	 id, user_id, sum, processed_at
	FROM 
		public.withdrawals 
	WHERE 
		id = $1
	`

	var w withdrawal.Withdrawal
	err := r.client.QueryRow(ctx, q, id).
		Scan(&w.ID, &w.UserID, &w.Sum, &w.ProcessedAt)
	if err != nil {
		return withdrawal.Withdrawal{}, err
	}

	w.SumFloat, err = strconv.ParseFloat(w.Sum, 64)
	if err != nil {
		return withdrawal.Withdrawal{}, err
	}

	return w, nil
}

func NewRepository(client postgres.Client, logger *logging.Logger) withdrawal.Repository {
	q := `
	CREATE TABLE IF NOT EXISTS public.withdrawals (
		id TEXT PRIMARY KEY,
		user_id UUID NOT NULL,
		sum TEXT NOT NULL,
		processed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
