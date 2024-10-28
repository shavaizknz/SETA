package paymentgateway

import (
	"seta/pkg/model"

	"github.com/shopspring/decimal"
)

type IPaymentGateway interface {
	Deposit(AccountID string, amount decimal.Decimal) (*model.TransactionResponse, *int, error)
	Withdraw(AccountID string, amount decimal.Decimal) (*model.TransactionResponse, *int, error)
}
