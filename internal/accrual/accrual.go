package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nickzhog/gophermart/internal/config"
	"github.com/nickzhog/gophermart/internal/service/order"
	"github.com/nickzhog/gophermart/pkg/logging"
)

type Scanner interface {
	StartScan(ctx context.Context) error
}

type scanner struct {
	Logger   *logging.Logger
	Cfg      *config.Config
	OrderRep order.Repository
}

func NewScanner(logger *logging.Logger, cfg *config.Config, orderRep order.Repository) Scanner {
	return &scanner{
		Logger:   logger,
		Cfg:      cfg,
		OrderRep: orderRep,
	}
}

func (s *scanner) StartScan(ctx context.Context) error {
	ticker := time.NewTicker(s.Cfg.Settings.AccrualScanInterval)
	for {
		select {
		case <-ticker.C:
			s.scan(ctx)
		case <-ctx.Done():
			s.Logger.Trace("orders scan stopped")
			return errors.New("context deadline exceeded")
		}
	}
}

func (s *scanner) scan(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()
	orders, err := s.OrderRep.FindForScanner(ctx)
	if err != nil {
		s.Logger.Error(err)
		return
	}
	for _, o := range orders {
		err = updateAccrual(ctx, s.Cfg.Settings.AccrualSystemAddress, &o)
		if err != nil {
			s.Logger.Error(err)
			continue
		}
		err = s.OrderRep.Update(ctx, &o)
		if err != nil {
			s.Logger.Error(err)
		}
	}
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
		return errors.New("too many requests")
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

	if order.AccrualFloat == ans.Accrual &&
		order.Status == ans.Status {
		return errors.New("nothing changed")
	}
	if order.ID != ans.Order {
		return errors.New("order ID changed")
	}

	order.Accrual = fmt.Sprintf("%g", ans.Accrual)
	order.AccrualFloat = ans.Accrual
	order.Status = ans.Status

	return nil
}
