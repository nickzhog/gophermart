package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccrualSumForProcessedOrders(t *testing.T) {
	tests := []struct {
		name string
		ords []Order
		want float64
	}{
		{
			name: "normal case",
			ords: []Order{
				{Status: StatusProcessed, AccrualFloat: 10},
				{Status: StatusProcessed, AccrualFloat: 10},
			},
			want: 20,
		},
		{
			name: "another status",
			ords: []Order{
				{Status: StatusProcessing, AccrualFloat: 10},
				{Status: StatusNew, AccrualFloat: 10},
			},
			want: 0,
		},
		{
			name: "processed and another status",
			ords: []Order{
				{Status: StatusProcessed, AccrualFloat: 10},
				{Status: StatusProcessed, AccrualFloat: 10},
				{Status: StatusInvalid, AccrualFloat: 10},
				{Status: StatusRegistered, AccrualFloat: 10},
			},
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			sum := AccrualSumForProcessedOrders(tt.ords)
			assert.Equal(tt.want, sum)
		})
	}
}

func TestNewOrder(t *testing.T) {
	type args struct {
		id    string
		usrID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive case",
			args:    args{id: "485542211", usrID: "usrID"},
			wantErr: false,
		},
		{
			name:    "wrong luhn orderID case",
			args:    args{id: "1233211", usrID: "usrID"},
			wantErr: true,
		},
		{
			name:    "empty usrID",
			args:    args{id: "485542211", usrID: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			_, err := NewOrder(tt.args.id, tt.args.usrID)
			assert.Equal(tt.wantErr, err != nil)
		})
	}
}
