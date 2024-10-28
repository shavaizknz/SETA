package controller

import (
	"seta/pkg/model"

	"github.com/shopspring/decimal"
)

type DepositRequest struct {
	AccountID string          `json:"account_id"`
	Amount    decimal.Decimal `json:"amount"`
}

type UpdateTransactionRequest struct {
	AccountID     string                  `json:"account_id"`
	TransactionID string                  `json:"transaction_id"`
	Status        model.TransactionStatus `json:"status"`
}
