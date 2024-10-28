package model

import "github.com/shopspring/decimal"

//---------------- API Data models ---------------- //

// Response model for gateway A transaction, used to simulate different response structures from different gateways
type GatewayATransactionResponse struct {
	Data TransactionDataA `json:"data"`
}

// Response model for gateway B transaction, used to simulate different response structures from different gateways
type GatewayBTransactionResponse TransactionData

type TransactionDataA TransactionData
type TransactionDataB TransactionData

type TransactionResponse struct { // Common TransactionModel
	Data TransactionData `json:"data" xml:"Data"`
}

type TransactionData struct { // Data model for transaction
	AccountID     string            `json:"account_id" xml:"AccountID"`
	TransactionID string            `json:"transaction_id" xml:"TransactionID"`
	Status        TransactionStatus `json:"status" xml:"Status"`
	Type          TransactionType   `json:"type" xml:"Type"`
	Amount        decimal.Decimal   `json:"amount" xml:"Amount"`
}

type TransactionStatus string

const (
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
	TransactionStatusPending TransactionStatus = "pending"
)

const TransactionStatuses = "success,failed,pending"

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"
)

//---------------- Database models ---------------- //

type TransactionDAO struct { // Data Access Object, used to interact with the database
	TransactionID string
	AccountID     string
	Amount        string
	Status        TransactionStatusDAO
	Type          TransactionTypeDAO
}

type TransactionStatusDAO string

const (
	TransactionStatusSuccessDAO TransactionStatusDAO = "success"
	TransactionStatusFailedDAO  TransactionStatusDAO = "failed"
)

type TransactionTypeDAO string

const (
	TransactionTypeDepositDAO  TransactionTypeDAO = "deposit"
	TransactionTypeWithdrawDAO TransactionTypeDAO = "withdraw"
)

//---------------- Mapping functions ---------------- //

func MapGatewayATransactionResponse(transactionResponse *GatewayATransactionResponse) TransactionResponse {
	return TransactionResponse{
		Data: TransactionData{
			AccountID:     transactionResponse.Data.AccountID,
			TransactionID: transactionResponse.Data.TransactionID,
			Status:        transactionResponse.Data.Status,
			Type:          transactionResponse.Data.Type,
			Amount:        transactionResponse.Data.Amount,
		},
	}
}

func MapGatewayBTransactionResponse(transactionResponse *GatewayBTransactionResponse) TransactionResponse {
	return TransactionResponse{
		Data: TransactionData{
			AccountID:     transactionResponse.AccountID,
			TransactionID: transactionResponse.TransactionID,
			Status:        transactionResponse.Status,
			Type:          transactionResponse.Type,
			Amount:        transactionResponse.Amount,
		},
	}
}

func MapTransactionResponseToTransactionDAO(transactionResponse *TransactionResponse) TransactionDAO {
	return TransactionDAO{
		AccountID:     transactionResponse.Data.AccountID,
		Amount:        transactionResponse.Data.Amount.String(),
		TransactionID: transactionResponse.Data.TransactionID,
		Status:        TransactionStatusDAO(transactionResponse.Data.Status),
		Type:          TransactionTypeDAO(transactionResponse.Data.Type),
	}
}

func MapTransactionDAOToTransactionResponse(transactionDAO *TransactionDAO) TransactionResponse {
	return TransactionResponse{
		Data: TransactionData{
			AccountID:     transactionDAO.AccountID,
			TransactionID: transactionDAO.TransactionID,
			Status:        TransactionStatus(transactionDAO.Status),
			Type:          TransactionType(transactionDAO.Type),
			Amount:        decimal.RequireFromString(transactionDAO.Amount),
		},
	}
}
