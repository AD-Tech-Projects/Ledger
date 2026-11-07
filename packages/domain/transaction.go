package domain

import "time"

type Transaction struct {
	ID        string
	CreatedAt time.Time

	Amount      int64
	Description string

	Type        TransactionType
	Status      TransactionStatus
	CompletedAt *time.Time

	AssetID string
	Asset   *Asset

	SourceWalletID *string
	SourceWallet   *Wallet

	DestinationWalletID *string
	DestinationWallet   *Wallet

	// NOTE - identifier for the operation that originated this transaction
	ReferenceID *string

	// NOTE - ID of the transaction that originated this transaction, in case of a reversal
	ReversalOriginID *string
	ReversalOrigin   *Transaction
}

type TransactionType string

const (
	TRANSACTION_DEPOSIT    TransactionType = "DEPOSIT"    // Deposit of funds into a wallet, typically from an external source
	TRANSACTION_WITHDRAWAL TransactionType = "WITHDRAWAL" // Withdrawal of funds from a wallet, typically to an external destination
	TRANSACTION_TRANSFER   TransactionType = "TRANSFER"   // Transfer of funds between two wallets within the system
	TRANSACTION_REVERSAL   TransactionType = "REVERSAL"   // Reversal of a previous transaction, typically due to an error or dispute

	TRANSACTION_YIELD    TransactionType = "YIELD"    // Yield generated from investments or interest-bearing accounts
	TRANSACTION_FEE      TransactionType = "FEE"      // Fee charged for a service or transaction
	TRANSACTION_PURCHASE TransactionType = "PURCHASE" // Purchase of an asset or service, typically from an external source
	TRANSACTION_SALE     TransactionType = "SALE"     // Sale of an asset or service, typically to an external destination

	TRANSACTION_RESERVE TransactionType = "RESERVE" // Reserve of funds in a wallet, typically for a scheduled transaction or future obligation
	TRANSACTION_RELEASE TransactionType = "RELEASE" // Release of reserved funds in a wallet, typically after a scheduled transaction or future obligation
	TRANSACTION_HOLD    TransactionType = "HOLD"    // Hold of funds in a wallet, typically for a pending transaction or dispute
	TRANSACTION_UNHOLD  TransactionType = "UNHOLD"  // Release of held funds in a wallet, typically after a pending transaction or dispute is resolved
)

type TransactionStatus string

const (
	TRANSACTION_PENDING   TransactionStatus = "PENDING"
	TRANSACTION_COMPLETED TransactionStatus = "COMPLETED"
	TRANSACTION_FAILED    TransactionStatus = "FAILED"
)
