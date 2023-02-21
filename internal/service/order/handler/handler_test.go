package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nickzhog/gophermart/internal/service/order"
	mock_order "github.com/nickzhog/gophermart/internal/service/order/mocks"
	"github.com/nickzhog/gophermart/internal/web/session"
	"github.com/nickzhog/gophermart/pkg/logging"
	"github.com/stretchr/testify/assert"
)

const (
	validOrderID   = "5880182"
	validOrderID2  = "71476808630764"
	validUserID    = "ValidID"
	validSessionID = "SessionID"
)

func prepareNewOrderHandler(ctrl *gomock.Controller) *handler {
	h := &handler{
		logger: logging.GetLogger(),
	}

	orderRep := mock_order.NewMockRepository(ctrl)
	orderRep.EXPECT().FindByID(gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(ctx context.Context, id string) (order.Order, error) {
			if id == validOrderID {
				return order.Order{ID: validOrderID, UserID: validUserID, Status: order.StatusNew, UploadAt: time.Now()}, nil
			}
			return order.Order{}, errors.New("not found")
		})

	orderRep.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(ctx context.Context, order *order.Order) error {
			return nil
		})
	h.Repositories.Order = orderRep

	return h
}
func Test_handler_newOrderHandler(t *testing.T) {
	tests := []struct {
		name        string
		requestBody []byte
		wantStatus  int
	}{
		{
			name:        "already have that order",
			requestBody: []byte(validOrderID),
			wantStatus:  http.StatusOK,
		},
		{
			name:        "new valid order",
			requestBody: []byte(validOrderID2),
			wantStatus:  http.StatusAccepted,
		},
		{
			name:        "wrong order (luhn)",
			requestBody: []byte(`1234`),
			wantStatus:  http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := prepareNewOrderHandler(ctrl)

			request := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBuffer(tt.requestBody))
			request = session.PutSessionDataInRequest(request, validSessionID, validUserID)

			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.newOrderHandler)
			handler.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(tt.wantStatus, res.StatusCode)
		})
	}
}

func prepareOrdersListHandler(ctrl *gomock.Controller) *handler {
	h := &handler{
		logger: logging.GetLogger(),
	}

	orderRep := mock_order.NewMockRepository(ctrl)
	orderRep.EXPECT().FindForUser(gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(ctx context.Context, usrID string) ([]order.Order, error) {
			if usrID == validUserID {
				return []order.Order{
					{ID: validOrderID, UserID: usrID, Status: order.StatusNew, UploadAt: time.Now()},
					{ID: validOrderID2, UserID: usrID, Status: order.StatusNew, UploadAt: time.Now()},
				}, nil
			}

			return []order.Order{}, order.ErrNoRows
		})

	h.Repositories.Order = orderRep

	return h
}

func Test_handler_getOrdersHandler(t *testing.T) {
	tests := []struct {
		name       string
		usrID      string
		wantStatus int
	}{
		{
			name:       "positive case",
			usrID:      validUserID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "user without orders",
			usrID:      "bad_user_id",
			wantStatus: http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := prepareOrdersListHandler(ctrl)

			request := httptest.NewRequest(http.MethodGet, "/api/user/orders", bytes.NewBuffer(nil))
			request = session.PutSessionDataInRequest(request, validSessionID, tt.usrID)

			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.getOrdersHandler)
			handler.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(tt.wantStatus, res.StatusCode)
		})
	}
}
