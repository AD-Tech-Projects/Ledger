package transaction

import (
	account "github.com/AD-Tech-Projects/Ledger/packages/domain/account"
)

type Transaction struct {
	id                 string
	amount             int64
	status             TransactionStatus
	transactionType    TransactionType
	sourceAccount      *account.Account
	destinationAccount *account.Account
}

type TransactionType string

const (
	DEPOSIT    TransactionType = "DEPOSIT"
	WITHDRAWAL TransactionType = "WITHDRAWAL"
)

type TransactionStatus string

const (
	PENDING   TransactionStatus = "PENDING"
	COMPLETED TransactionStatus = "COMPLETED"
	FAILED    TransactionStatus = "FAILED"
)

type TransactionError struct {
	message string
}
