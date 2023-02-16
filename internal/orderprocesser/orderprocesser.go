package orderprocesser

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

type OrderProcesser interface {
	StartScan(ctx context.Context) error
}

type orderProcesser struct {
	Logger   *logging.Logger
	Cfg      *config.Config
	OrderRep order.Repository
}

func NewProcesser(logger *logging.Logger, cfg *config.Config, orderRep order.Repository) OrderProcesser {
	return &orderProcesser{
		Logger:   logger,
		Cfg:      cfg,
		OrderRep: orderRep,
	}
}

func (p *orderProcesser) StartScan(ctx context.Context) error {
	ticker := time.NewTicker(p.Cfg.Settings.AccrualScanInterval)
	for {
		select {
		case <-ticker.C:
			p.scan(ctx)
		case <-ctx.Done():
			p.Logger.Trace("orders processing exited properly")
			return nil
		}
	}
}

func (p *orderProcesser) scan(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()
	orders, err := p.OrderRep.FindForScanner(ctx)
	if err != nil {
		p.Logger.Error(err)
		return
	}
	for _, o := range orders {
		err = updateAccrual(ctx, p.Cfg.Settings.AccrualSystemAddress, &o)
		if err != nil {
			p.Logger.Error(err)
			continue
		}
		err = p.OrderRep.Update(ctx, &o)
		if err != nil {
			p.Logger.Error(err)
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
