package repository

import (
	"context"
	"seta/pkg/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ITransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction model.TransactionDAO) error
	GetTransaction(ctx context.Context, transactionID string) (model.TransactionDAO, error)
	UpdateTransaction(ctx context.Context, transaction model.TransactionDAO) error
}

type TransactionRepository struct {
	DB *pgxpool.Pool
}

func TransactionRepositoryProvider(db *pgxpool.Pool) ITransactionRepository {
	return &TransactionRepository{DB: db}
}

func (tr *TransactionRepository) CreateTransaction(ctx context.Context, transaction model.TransactionDAO) error {
	//account_id transaction_id amount status transaction_type
	_, err := tr.DB.Exec(ctx, InsertTransactionQuery, transaction.AccountID, transaction.TransactionID, transaction.Amount, transaction.Status, transaction.Type)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TransactionRepository) GetTransaction(ctx context.Context, transactionID string) (model.TransactionDAO, error) {
	var transaction model.TransactionDAO
	err := tr.DB.QueryRow(ctx, GetTransactionQuery, transactionID).Scan(&transaction.AccountID, &transaction.TransactionID, &transaction.Amount, &transaction.Status, &transaction.Type)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (tr *TransactionRepository) UpdateTransaction(ctx context.Context, transaction model.TransactionDAO) error {
	_, err := tr.DB.Exec(ctx, UpdateTransactionQuery, transaction.AccountID, transaction.TransactionID, transaction.Status)
	if err != nil {
		return err
	}
	return nil
}
