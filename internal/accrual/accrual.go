package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nickzhog/gophermart/internal/config"
	"github.com/nickzhog/gophermart/internal/repositories"
	"github.com/nickzhog/gophermart/internal/service/order"
	"github.com/nickzhog/gophermart/pkg/logging"
)

func OrdersScanStart(logger *logging.Logger, cfg *config.Config, reps repositories.Repositories) {
	for {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			orders, err := reps.Order.FindForScanner(ctx)
			if err != nil {
				logger.Error(err)
				return
			}
			for _, o := range orders {
				err = getAccrual(ctx, cfg.Settings.AccrualSystemAddress, &o)
				if err != nil {
					logger.Error(err)
					continue
				}
				reps.Order.Update(ctx, &o)
			}
			time.Sleep(time.Millisecond * 150)
		}()
	}
}

type Answer struct {
	Order   string  `json:"order,omitempty"`
	Status  string  `json:"status,omitempty"`
	Accrual float64 `json:"accrual,omitempty"`
}

func getAccrual(ctx context.Context, url string, order *order.Order) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/%s", url, order.ID), nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		return getAccrual(ctx, url, order)
	}

	var ans Answer
	err = json.NewDecoder(res.Body).Decode(&ans)
	if err != nil {
		return err
	}

	order.Accrual = fmt.Sprintf("%g", ans.Accrual)
	order.AccrualFloat = ans.Accrual
	order.Status = ans.Status

	return nil
}
