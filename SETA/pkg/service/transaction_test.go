// internal/service/transaction_service_test.go
package service

import (
	"context"

	"seta/pkg/clients/paymentgateway"
	"seta/pkg/model"
	"seta/pkg/repository"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction_Success(t *testing.T) {
	// Initialize mock repository and service
	amount := decimal.NewFromFloat(100.0)
	transactionExpected := model.TransactionResponse{
		Data: model.TransactionData{
			TransactionID: "txn123",
			AccountID:     "acc123",
			Amount:        amount,
			Status:        model.TransactionStatusSuccess,
			Type:          model.TransactionTypeWithdraw,
		},
	}

	transactionDAO := model.MapTransactionResponseToTransactionDAO(&transactionExpected)

	mockRepo := repository.MockTransactionRepositoryProvider(&transactionDAO, false, nil)
	mockPaymentGatewayClient := paymentgateway.MockClientProvider(&transactionExpected, 200, nil, false)
	service := TransactionServiceProvider(mockRepo, []paymentgateway.IPaymentGateway{mockPaymentGatewayClient})

	// Test CreateTransaction
	transactionActual, err := service.CreateTransaction(context.Background(), transactionExpected.Data.AccountID, transactionExpected.Data.Amount, transactionExpected.Data.Type)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, transactionActual)
	assert.Equal(t, transactionExpected, *transactionActual)
}

func TestCreateTransaction_TwoGateways_OneActive_Success(t *testing.T) {
	// Initialize mock repository and service
	amount := decimal.NewFromFloat(100.0)
	transactionExpected := model.TransactionResponse{
		Data: model.TransactionData{
			TransactionID: "txn123",
			AccountID:     "acc123",
			Amount:        amount,
			Status:        model.TransactionStatusSuccess,
			Type:          model.TransactionTypeWithdraw,
		},
	}

	transactionDAO := model.MapTransactionResponseToTransactionDAO(&transactionExpected)

	mockRepo := repository.MockTransactionRepositoryProvider(&transactionDAO, false, nil)
	mockPaymentGatewayClient1 := paymentgateway.MockClientProvider(nil, 500, nil, false)
	mockPaymentGatewayClient2 := paymentgateway.MockClientProvider(&transactionExpected, 200, nil, false)
	service := TransactionServiceProvider(mockRepo, []paymentgateway.IPaymentGateway{mockPaymentGatewayClient1, mockPaymentGatewayClient2})

	// Test CreateTransaction
	transactionActual, err := service.CreateTransaction(context.Background(), transactionExpected.Data.AccountID, transactionExpected.Data.Amount, transactionExpected.Data.Type)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, transactionActual)
	assert.Equal(t, transactionExpected, *transactionActual)
}

func TestCreateTransaction_TwoGateways_OneActive_Failure(t *testing.T) {
	// Initialize mock repository and service
	amount := decimal.NewFromFloat(100.0)
	transactionExpected := model.TransactionResponse{
		Data: model.TransactionData{
			TransactionID: "txn123",
			AccountID:     "acc123",
			Amount:        amount,
			Status:        model.TransactionStatusSuccess,
			Type:          model.TransactionTypeWithdraw,
		},
	}

	transactionDAO := model.MapTransactionResponseToTransactionDAO(&transactionExpected)

	mockRepo := repository.MockTransactionRepositoryProvider(&transactionDAO, false, nil)
	mockPaymentGatewayClient1 := paymentgateway.MockClientProvider(nil, 500, nil, false)
	mockPaymentGatewayClient2 := paymentgateway.MockClientProvider(nil, 500, nil, false)
	service := TransactionServiceProvider(mockRepo, []paymentgateway.IPaymentGateway{mockPaymentGatewayClient1, mockPaymentGatewayClient2})

	// Test CreateTransaction
	transactionActual, err := service.CreateTransaction(context.Background(), transactionExpected.Data.AccountID, transactionExpected.Data.Amount, transactionExpected.Data.Type)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, transactionActual)
}
