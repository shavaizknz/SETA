package model

import "github.com/shopspring/decimal"

type WithdrawRequest struct {
	AccountID string          `json:"account_id" xml:"AccountID"`
	Amount    decimal.Decimal `json:"amount" xml:"Amount"`
}
