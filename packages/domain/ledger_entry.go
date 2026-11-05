package domain

type LedgerEntry struct {
	ID string

	LedgerTransactionID string
	LedgerTransaction   *LedgerTransaction

	LedgerAccountID string
	LedgerAccount   *LedgerAccount

	AssetID string
	Asset   *Asset

	Amount    int64
	Direction LedgerBalanceDirection
	//  NOTE - the sequence number of the entry within the transaction, starting from 1
	Sequence int64
}
