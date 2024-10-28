package handler

import (
	"context"
	"errors"
	"seta/pkg/clients/paymentgateway"
	"seta/pkg/logger"
	"seta/pkg/model"

	"github.com/shopspring/decimal"
)

func CreateTransactionFromPaymentGateways(ctx context.Context, paymentGateways []paymentgateway.IPaymentGateway, accountID string, amount decimal.Decimal, transactionType model.TransactionType) (*model.TransactionResponse, error) {
	// loop through the payment gateways to create a transaction
	for _, paymentGateway := range paymentGateways {
		switch transactionType {
		case model.TransactionTypeDeposit:
			transactionResponse, statusCode, err := paymentGateway.Deposit(accountID, amount)
			// if the error is a server error, retry with the next payment gateway by not returning the error or breaking the loop
			if statusCode != nil && (*statusCode > 400 || *statusCode == 0) {
				logger.WithRequestID(ctx).Errorf("payment gateway failed with status code %d", *statusCode)
				transactionResponse = nil
				// if the error is not a server error, then there is no need to retry with the next payment gateway as there is a bad request
			} else if err != nil {
				logger.WithRequestID(ctx).Errorf("payment gateway failed with status code %d . error: %d", *statusCode, err)
				return nil, err
			} else { // success
				logger.WithRequestID(ctx).Infof("payment gateway succeeded with status code %d", *statusCode)
				return transactionResponse, nil
			}
		case model.TransactionTypeWithdraw:
			transactionResponse, statusCode, err := paymentGateway.Withdraw(accountID, amount)
			// if the error is a server error, retry with the next payment gateway
			if statusCode != nil && (*statusCode > 400 || *statusCode == 0) {
				logger.WithRequestID(ctx).Errorf("payment gateway failed with status code %d", *statusCode)
				transactionResponse = nil
				// if the error is not a server error, then there is no need to retry with the next payment gateway as there is a bad request
			} else if err != nil {
				logger.WithRequestID(ctx).Errorf("payment gateway failed with status code %d . error: %d", *statusCode, err)
				return nil, err
			} else { // success
				logger.WithRequestID(ctx).Infof("payment gateway succeeded with status code %d", *statusCode)
				return transactionResponse, nil
			}
		default:
			return nil, errors.New("invalid transaction type")
		}
	}

	return nil, errors.New("all payment gateways failed")
}
