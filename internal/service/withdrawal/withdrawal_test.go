package withdrawal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWithdrawal(t *testing.T) {
	type args struct {
		orderID string
		usrID   string
		sum     float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive case",
			args:    args{orderID: "12321"},
			wantErr: false,
		},
		{
			name:    "wrong order id",
			args:    args{orderID: "abc"},
			wantErr: true,
		},
		{
			name:    "wrong order id 2",
			args:    args{orderID: "0"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			_, err := NewWithdrawal(tt.args.orderID, tt.args.usrID, tt.args.sum)
			assert.Equal(tt.wantErr, err != nil)
		})
	}
}
