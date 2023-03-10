package withdrawal

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Withdrawal struct {
	ID          string    `json:"order,omitempty"`
	UserID      string    `json:"-"`
	Sum         string    `json:"-"`
	SumFloat    float64   `json:"sum,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}

type WithdrawalRequest struct {
	Order string  `json:"order,omitempty"`
	Sum   float64 `json:"sum,omitempty"`
}

var ErrNoRows = errors.New("withdrawal not found")

func ParseWithdrawalRequest(data []byte) (WithdrawalRequest, error) {
	var wr WithdrawalRequest
	err := json.Unmarshal(data, &wr)
	if err != nil {
		return WithdrawalRequest{}, err
	}
	return wr, nil
}

func NewWithdrawal(orderID, usrID string, sum float64) (Withdrawal, error) {
	w := Withdrawal{
		ID:       orderID,
		UserID:   usrID,
		Sum:      fmt.Sprintf("%g", sum),
		SumFloat: sum,
	}
	idInt, err := strconv.Atoi(orderID)
	if err != nil {
		return Withdrawal{}, err
	}

	if idInt < 1 {
		return Withdrawal{}, errors.New("wrong order")
	}

	return w, nil
}

func SumForWithdrawals(wdls []Withdrawal) float64 {
	answer := 0.0
	for _, v := range wdls {
		answer += v.SumFloat
	}
	return answer
}
