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
	PostedAt    *time.Time // PostedAt is filled when the entries are posted to the ledger.

	ReversalOriginID *string // ReversalOriginID is the ledger transaction reversed by this one, in case of a reversal.
	ReversalOrigin   *LedgerTransaction
}

type LedgerTransactionStatus string

const (
	LEDGER_TRANSACTION_PENDING LedgerTransactionStatus = "PENDING" // Entries are being built/validated, not yet part of the books
	LEDGER_TRANSACTION_POSTED  LedgerTransactionStatus = "POSTED"  // Entries were posted and are now immutable
	LEDGER_TRANSACTION_FAILED  LedgerTransactionStatus = "FAILED"  // Posting failed (e.g. unbalanced entries); no entry affects balances
)
