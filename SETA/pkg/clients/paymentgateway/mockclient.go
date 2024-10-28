package paymentgateway

import (
	"fmt"
	"seta/pkg/model"

	"github.com/shopspring/decimal"
)

type MockClient struct {
	TransactionResponse *model.TransactionResponse
	StatusCode          int
	Err                 error
	IsServiceDown       bool
}

func MockClientProvider(transactionResponse *model.TransactionResponse, statusCode int, err error, isServiceDown bool) IPaymentGateway {
	return &MockClient{
		TransactionResponse: transactionResponse,
		StatusCode:          statusCode,
		Err:                 err,
		IsServiceDown:       isServiceDown,
	}
}

func (c *MockClient) Deposit(AccountID string, amount decimal.Decimal) (*model.TransactionResponse, *int, error) {
	if c.IsServiceDown {
		return nil, &c.StatusCode, fmt.Errorf("service is down")
	}

	return c.TransactionResponse, &c.StatusCode, c.Err
}

func (c *MockClient) Withdraw(AccountID string, amount decimal.Decimal) (*model.TransactionResponse, *int, error) {
	if c.IsServiceDown {
		return nil, &c.StatusCode, fmt.Errorf("service is down")
	}

	return c.TransactionResponse, &c.StatusCode, c.Err
}
