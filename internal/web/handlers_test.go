package web

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nickzhog/gophermart/internal/config"
	"github.com/nickzhog/gophermart/internal/service/user"
	mock_user "github.com/nickzhog/gophermart/internal/service/user/mocks"
	"github.com/nickzhog/gophermart/internal/web/session"
	mock_session "github.com/nickzhog/gophermart/internal/web/session/mocks"
	"github.com/nickzhog/gophermart/pkg/logging"
	"github.com/stretchr/testify/assert"
)

const (
	validUsrID           = "ValidID"
	validUsrLogin        = "ValidLogin"
	validUsrPassword     = "Password1234"
	validUsrPasswordHash = "$2a$10$BlVn15UubSt8R7/w99TMceHREqd6PEwk1d42zrnKJQ6fs5XUp3Wqa"

	validSessionID = "ValidSessionID"
)

func prepareHandlerData(ctrl *gomock.Controller) *HandlerData {
	cfg := config.Config{}
	h := &HandlerData{
		Logger: logging.GetLogger(),
		Cfg:    &cfg,
	}

	usrRep := mock_user.NewMockRepository(ctrl)
	usrRep.EXPECT().FindByLogin(gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(ctx context.Context, login string) (user.User, error) {
			if login == validUsrLogin {
				return user.User{ID: validUsrID, Login: validUsrLogin, PasswordHash: validUsrPasswordHash}, nil
			}
			return user.User{}, errors.New("not found")
		})
	h.Repositories.User = usrRep

	sessionRep := mock_session.NewMockRepository(ctrl)
	sessionRep.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(ctx context.Context, usrID string) (session.Session, error) {
			if usrID == validUsrID {
				return session.Session{
					ID:       validSessionID,
					UserID:   usrID,
					CreateAt: time.Now(),
					IsActive: true,
				}, nil
			}
			return session.Session{}, errors.New("not found")
		})
	h.Repositories.Session = sessionRep

	return h
}

func TestHandlerData_loginHandler(t *testing.T) {

	tests := []struct {
		name        string
		requestBody []byte
		wantStatus  int
	}{
		{
			name:        "positive case",
			requestBody: []byte(`{"login":"ValidLogin","password":"Password1234"}`),
			wantStatus:  http.StatusOK,
		},
		{
			name:        "wrong json",
			requestBody: []byte(`{"l`),
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "wrong password",
			requestBody: []byte(`{"login":"ValidLogin","password":"wrong_password"}`),
			wantStatus:  http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := prepareHandlerData(ctrl)

			request := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(tt.requestBody))
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.loginHandler)
			handler.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(tt.wantStatus, res.StatusCode)
		})
	}
}
