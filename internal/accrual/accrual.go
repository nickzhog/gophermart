package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nickzhog/gophermart/internal/config"
	"github.com/nickzhog/gophermart/internal/repositories"
	"github.com/nickzhog/gophermart/internal/service/order"
	"github.com/nickzhog/gophermart/pkg/logging"
)

func StartOrdersScan(logger *logging.Logger, cfg *config.Config, reps repositories.Repositories) {
	go func() {
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
					err = updateAccrual(ctx, cfg.Settings.AccrualSystemAddress, &o)
					if err != nil {
						logger.Error(err)
						continue
					}
					err = reps.Order.Update(ctx, &o)
					if err != nil {
						logger.Error(err)
					}
				}
			}()

			time.Sleep(cfg.Settings.AccrualScanInterval)
		}
	}()
}

type Answer struct {
	Order   string  `json:"order,omitempty"`
	Status  string  `json:"status,omitempty"`
	Accrual float64 `json:"accrual,omitempty"`
}

func updateAccrual(ctx context.Context, url string, order *order.Order) error {
	fullURL := fmt.Sprintf("%s/api/orders/%s", url, order.ID)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		time.Sleep(time.Second)
		return updateAccrual(ctx, url, order)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var ans Answer
	err = json.Unmarshal(body, &ans)
	if err != nil {
		return fmt.Errorf("url:%s,body:%s, err:%s",
			fullURL, string(body), err.Error())
	}

	order.Accrual = fmt.Sprintf("%g", ans.Accrual)
	order.AccrualFloat = ans.Accrual
	order.Status = ans.Status

	return nil
}
