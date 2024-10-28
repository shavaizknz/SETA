package repository

const (
	InsertTransactionQuery = `INSERT INTO transactions (account_id, transaction_id, amount, status, type)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (account_id, transaction_id) DO UPDATE SET status = $4`
	GetTransactionQuery    = "SELECT account_id, transaction_id, amount, status, type FROM transactions WHERE transaction_id = $1"
	UpdateTransactionQuery = "UPDATE transactions SET status = $3 WHERE account_id = $1 AND transaction_id = $2"
)
