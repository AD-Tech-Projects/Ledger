package transaction

import (
	account "github.com/AD-Tech-Projects/Ledger/packages/domain/account"
)

type Transaction struct {
	Id                 string
	Amount             int64
	Status             TransactionStatus
	TransactionType    TransactionType
	SourceAccount      *account.Account
	DestinationAccount *account.Account
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
	operation TransactionType
	message   string
}
