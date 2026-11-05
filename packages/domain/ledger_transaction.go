package domain

import "time"

// NOTE - a ledger transaction is a record of a transaction that has been posted to the ledger
type LedgerTransaction struct {
	ID string

	// NOTE - the transaction that originated this ledger transaction, should be unique for each ledger transaction
	TransactionID string
	Transaction   *Transaction

	Status      LedgerTransactionStatus
	Description string
	PostedAt    *time.Time

	ReversalOriginID *string
	ReversalOrigin   *LedgerTransaction
}

type LedgerTransactionStatus string

const (
	LEDGER_TRANSACTION_PENDING LedgerTransactionStatus = "PENDING"
	LEDGER_TRANSACTION_POSTED  LedgerTransactionStatus = "POSTED"
	LEDGER_TRANSACTION_FAILED  LedgerTransactionStatus = "FAILED"
)
