package service

import (
	"context"
	"fmt"
	"seta/pkg/clients/paymentgateway"
	"seta/pkg/handler"
	"seta/pkg/logger"
	"seta/pkg/model"
	"seta/pkg/repository"

	"github.com/jackc/pgx/v4"
	"github.com/shopspring/decimal"
)

type ITransactionService interface {
	CreateTransaction(ctx context.Context, accountID string, amount decimal.Decimal, transactionType model.TransactionType) (*model.TransactionResponse, error)
	GetTransaction(ctx context.Context, transactionID string) (*model.TransactionResponse, error)
	UpdateTransaction(ctx context.Context, accountID string, transactionID string, status model.TransactionStatus) error
}

type TransactionService struct {
	TransactionRepository repository.ITransactionRepository
	PaymentGateways       []paymentgateway.IPaymentGateway
}

func TransactionServiceProvider(transactionRepository repository.ITransactionRepository, paymentGateways []paymentgateway.IPaymentGateway) ITransactionService {
	return &TransactionService{
		TransactionRepository: transactionRepository,
		PaymentGateways:       paymentGateways,
	}
}

func (ts *TransactionService) CreateTransaction(ctx context.Context, accountID string, amount decimal.Decimal, transactionType model.TransactionType) (*model.TransactionResponse, error) {
	transactionResponse, err := handler.CreateTransactionFromPaymentGateways(ctx, ts.PaymentGateways, accountID, amount, transactionType)
	if err != nil {
		return nil, err
	}

	transactionDAO := model.MapTransactionResponseToTransactionDAO(transactionResponse)
	err = ts.TransactionRepository.CreateTransaction(ctx, transactionDAO)
	if err != nil {
		logger.WithRequestID(ctx).Errorf("failed to create transaction in database: %v", err)
		return nil, err
	}

	return transactionResponse, nil
}

func (ts *TransactionService) GetTransaction(ctx context.Context, transactionID string) (*model.TransactionResponse, error) {
	transactionDAO, err := ts.TransactionRepository.GetTransaction(context.Background(), transactionID)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}

	logger.WithRequestID(ctx).Infof("transaction found: %v", transactionDAO)

	transactionResponse := model.MapTransactionDAOToTransactionResponse(&transactionDAO)
	return &transactionResponse, nil
}

func (ts *TransactionService) UpdateTransaction(ctx context.Context, accountID string, transactionID string, status model.TransactionStatus) error {
	transactionDAO, err := ts.TransactionRepository.GetTransaction(ctx, transactionID)
	if err != nil {
		return err
	}

	if transactionDAO.AccountID != accountID {
		return fmt.Errorf("transaction does not belong to account")
	}

	transactionDAO.Status = model.TransactionStatusDAO(status)
	err = ts.TransactionRepository.UpdateTransaction(ctx, transactionDAO)
	if err != nil {
		return err
	}

	return nil
}
