package model

import "github.com/shopspring/decimal"

type DepositRequest struct {
	AccountID string          `json:"account_id" xml:"AccountID"`
	Amount    decimal.Decimal `json:"amount" xml:"Amount"`
}
