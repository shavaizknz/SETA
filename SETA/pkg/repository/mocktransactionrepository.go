package repository

import (
	"context"
	"errors"
	"seta/pkg/model"
)

// MockTransactionRepository simulates a TransactionRepository for testing purposes
type MockTransactionRepository struct {
	Transaction   *model.TransactionDAO
	ShouldFail    bool
	ExpectedError error
}

// NewMockTransactionRepository initializes the mock with an empty transactions map
func MockTransactionRepositoryProvider(transaction *model.TransactionDAO, shouldFail bool, expectedError error) *MockTransactionRepository {
	return &MockTransactionRepository{
		Transaction:   transaction,
		ShouldFail:    shouldFail,
		ExpectedError: expectedError,
	}
}

// CreateTransaction simulates creating a transaction, optionally failing based on configuration
func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, transaction model.TransactionDAO) error {
	if m.ShouldFail {
		return m.ExpectedError
	}
	m.Transaction = &transaction
	return nil
}

// GetTransaction simulates retrieving a transaction, returning an error if ShouldFail is set
func (m *MockTransactionRepository) GetTransaction(ctx context.Context, transactionID string) (model.TransactionDAO, error) {
	if m.ShouldFail {
		return model.TransactionDAO{}, m.ExpectedError
	}
	if m.Transaction == nil || m.Transaction.TransactionID != transactionID {
		return model.TransactionDAO{}, errors.New("transaction not found")
	}

	return *m.Transaction, nil
}

// UpdateTransaction simulates updating a transaction, returning an error if ShouldFail is set
func (m *MockTransactionRepository) UpdateTransaction(ctx context.Context, transaction model.TransactionDAO) error {
	if m.ShouldFail {
		return m.ExpectedError
	}

	m.Transaction = &transaction
	return nil
}
