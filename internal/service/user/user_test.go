package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive case",
			args:    args{login: "Login", password: "Password1234"},
			wantErr: false,
		},
		{
			name:    "empty login",
			args:    args{login: "", password: "Password1234"},
			wantErr: true,
		},
		{
			name:    "empty password",
			args:    args{login: "Login", password: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			usr, err := NewUser(tt.args.login, tt.args.password)

			assert.EqualValues(tt.wantErr, err != nil)
			if !tt.wantErr {
				err = bcrypt.CompareHashAndPassword(
					[]byte(usr.PasswordHash), []byte(tt.args.password))
				assert.NoError(err)
			}
		})
	}
}
