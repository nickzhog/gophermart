package db

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgconn"
	"github.com/nickzhog/gophermart/internal/postgres"
	"github.com/nickzhog/gophermart/internal/service/order"
	"github.com/nickzhog/gophermart/pkg/logging"
)

type repository struct {
	client postgres.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, o *order.Order) error {
	q := `
		INSERT INTO public.orders 
		    (id, user_id, accrual) 
		VALUES 
		    ($1, $2, $3) 
		RETURNING status, upload_at
	`
	err := r.client.QueryRow(ctx, q, o.ID, o.UserID, o.Accrual).
		Scan(&o.Status, &o.UploadAt)
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

func (r *repository) FindByID(ctx context.Context, id string) (order.Order, error) {
	q := `
	SELECT
		id, user_id, status, 
		accrual, upload_at
	FROM
		public.orders 
	WHERE 
		id = $1
	`

	var o order.Order
	err := r.client.QueryRow(ctx, q, id).
		Scan(&o.ID, &o.UserID, &o.Status,
			&o.Accrual, &o.UploadAt)
	if err != nil {
		return order.Order{}, err
	}

	if o.AccrualFloat, err = strconv.ParseFloat(o.Accrual, 64); err != nil {
		return order.Order{}, err
	}

	return o, nil
}

func (r *repository) FindForUser(ctx context.Context, usrID string) ([]order.Order, error) {
	q := `
		SELECT 
			id, user_id, status, 
			accrual, upload_at
		FROM public.orders 
		WHERE user_id = $1;
	`

	rows, err := r.client.Query(ctx, q, usrID)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	orders := make([]order.Order, 0)

	for rows.Next() {
		var o order.Order

		err = rows.Scan(&o.ID, &o.UserID, &o.Status,
			&o.Accrual, &o.UploadAt)

		if err != nil {
			r.logger.Error(err)
			return nil, err
		}
		if o.AccrualFloat, err = strconv.ParseFloat(o.Accrual, 64); err != nil {
			r.logger.Error(err)
			return nil, err
		}

		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return orders, nil
}

func (r *repository) Update(ctx context.Context, o *order.Order) error {

	o.Accrual = fmt.Sprintf("%g", o.AccrualFloat)

	q := `
		UPDATE public.orders 
		SET
		 status = $1,
		 accrual = $2
		WHERE id = $3
	`

	_, err := r.client.Exec(ctx, q,
		o.Status, o.Accrual, o.ID)

	return err

}

func (r *repository) FindForScanner(ctx context.Context) ([]order.Order, error) {
	q := `
		SELECT 
			id, user_id, status, 
			accrual, upload_at
		FROM 
			public.orders 
		WHERE 
			status = $1
			or status = $2 
			or status = $3
	`

	rows, err := r.client.Query(ctx, q,
		order.StatusNew, order.StatusProccessing, order.StatusProcessed)
	if err != nil {
		return nil, err
	}

	orders := make([]order.Order, 0)

	for rows.Next() {
		var o order.Order

		err = rows.Scan(&o.ID, &o.UserID, &o.Status,
			&o.Accrual, &o.UploadAt)

		if err != nil {
			return nil, err
		}
		if o.AccrualFloat, err = strconv.ParseFloat(o.Accrual, 64); err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func NewRepository(client postgres.Client, logger *logging.Logger) order.Repository {
	q := `
	CREATE TABLE IF NOT EXISTS public.orders (
		id TEXT PRIMARY KEY,
		user_id UUID NOT NULL,
		status TEXT NOT NULL DEFAULT 'NEW',
		accrual TEXT NOT NULL,
		upload_at TIMESTAMP NOT NULL  DEFAULT CURRENT_TIMESTAMP,
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
