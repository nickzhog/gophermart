package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nickzhog/gophermart/internal/service/order"
	mock_order "github.com/nickzhog/gophermart/internal/service/order/mocks"
	"github.com/nickzhog/gophermart/internal/service/user"
	mock_user "github.com/nickzhog/gophermart/internal/service/user/mocks"
	"github.com/nickzhog/gophermart/internal/service/withdrawal"
	mock_withdrawal "github.com/nickzhog/gophermart/internal/service/withdrawal/mocks"
	"github.com/nickzhog/gophermart/internal/web/session"
	"github.com/nickzhog/gophermart/pkg/logging"
	"github.com/stretchr/testify/assert"
)

const (
	validUsrID = "ValidUsrID"

	alreadyUsedID = "used_id"

	balance   float64 = 22
	withdrawn float64 = 8
)

var (
	order1 = order.Order{ID: "5880182", UserID: validUsrID,
		Status: order.StatusProcessed, Accrual: "10", AccrualFloat: 10, UploadAt: time.Now()}

	order2 = order.Order{ID: "71476808630764", UserID: validUsrID,
		Status: order.StatusProcessed, Accrual: "20", AccrualFloat: 20, UploadAt: time.Now()}

	//

	withdrawal1 = withdrawal.Withdrawal{ID: "123", UserID: validUsrID, Sum: "5", SumFloat: 5, ProcessedAt: time.Now()}
	withdrawal2 = withdrawal.Withdrawal{ID: "321", UserID: validUsrID, Sum: "3", SumFloat: 3, ProcessedAt: time.Now()}
)

func prepareHandler(ctrl *gomock.Controller) *handler {
	h := &handler{
		logger: logging.GetLogger(),
	}

	usrRep := mock_user.NewMockRepository(ctrl)
	usrRep.EXPECT().FindByID(gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(ctx context.Context, id string) (user.User, error) {
			if id == validUsrID {
				return user.User{ID: validUsrID}, nil
			}
			return user.User{}, user.ErrNoRows
		})
	h.Repositories.User = usrRep

	orderRep := mock_order.NewMockRepository(ctrl)
	orderRep.EXPECT().FindForUser(gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(ctx context.Context, usrID string) ([]order.Order, error) {
			if usrID == validUsrID {
				return []order.Order{order1, order2}, nil
			}

			return []order.Order{}, order.ErrNoRows
		})

	h.Repositories.Order = orderRep

	withdrawalRep := mock_withdrawal.NewMockRepository(ctrl)
	withdrawalRep.EXPECT().FindForUser(gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(ctx context.Context, usrID string) ([]withdrawal.Withdrawal, error) {
			if usrID == validUsrID {
				return []withdrawal.Withdrawal{withdrawal1, withdrawal2}, nil
			}

			return []withdrawal.Withdrawal{}, withdrawal.ErrNoRows
		})

	withdrawalRep.EXPECT().FindByID(gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(ctx context.Context, id string) (withdrawal.Withdrawal, error) {
			if id == alreadyUsedID {
				return withdrawal.Withdrawal{}, nil
			}

			return withdrawal.Withdrawal{}, withdrawal.ErrNoRows
		})

	withdrawalRep.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(ctx context.Context, wdl *withdrawal.Withdrawal) error {
			return nil
		})

	h.Repositories.Withdrawal = withdrawalRep

	return h
}

type Balance struct {
	Balance   float64 `json:"current,omitempty"`
	Withdrawn float64 `json:"withdrawn,omitempty"`
}

func Test_handler_balanceHandler(t *testing.T) {
	tests := []struct {
		name       string
		usrID      string
		wantStatus int
	}{
		{
			name:       "positive case",
			usrID:      validUsrID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "wrong user",
			usrID:      "bad_user_id",
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := prepareHandler(ctrl)

			request := httptest.NewRequest(http.MethodGet, "/api/user/balance", bytes.NewBuffer(nil))
			request = session.PutSessionDataInRequest(request, "session", tt.usrID)

			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.balanceHandler)
			handler.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(tt.wantStatus, res.StatusCode)
			if res.StatusCode == http.StatusOK {
				var ans Balance
				err := json.NewDecoder(res.Body).Decode(&ans)
				assert.NoError(err)
				assert.Equal(ans.Balance, balance)
				assert.Equal(ans.Withdrawn, withdrawn)
			}
		})
	}
}

func Test_handler_withdrawActionHandler(t *testing.T) {
	tests := []struct {
		name        string
		requestBody []byte
		wantStatus  int
	}{
		{
			name:        "positive case",
			requestBody: []byte(`{"order":"123321", "sum":20}`),
			wantStatus:  http.StatusOK,
		},
		{
			name:        "wrong json",
			requestBody: []byte(`{"orde`),
			wantStatus:  http.StatusUnprocessableEntity,
		},
		{
			name:        "sum more than balance",
			requestBody: []byte(`{"order":"123321", "sum":120}`),
			wantStatus:  http.StatusPaymentRequired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := prepareHandler(ctrl)

			request := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBuffer(tt.requestBody))
			request = session.PutSessionDataInRequest(request, "session", validUsrID)

			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.withdrawActionHandler)
			handler.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(tt.wantStatus, res.StatusCode)
		})
	}
}

func Test_handler_withdrawalsHandler(t *testing.T) {
	tests := []struct {
		name       string
		usrID      string
		wantStatus int
	}{
		{
			name:       "positive case",
			usrID:      validUsrID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "wrong user",
			usrID:      "bad_user_id",
			wantStatus: http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := prepareHandler(ctrl)

			request := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", bytes.NewBuffer(nil))
			request = session.PutSessionDataInRequest(request, "session", tt.usrID)

			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.withdrawalsHandler)
			handler.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(tt.wantStatus, res.StatusCode)
			if res.StatusCode == http.StatusOK {
				var withdrawals []withdrawal.Withdrawal
				err := json.NewDecoder(res.Body).Decode(&withdrawals)
				assert.NoError(err)
			}
		})
	}
}
