package domain

type LedgerEntry struct {
	ID string

	LedgerTransactionID string
	LedgerTransaction   *LedgerTransaction

	LedgerAccountID string
	LedgerAccount   *LedgerAccount

	AssetID string // AssetID must match the asset of the ledger account.
	Asset   *Asset

	Amount    int64                  // Amount is always positive and expressed in the smallest unit of the asset.
	Direction LedgerBalanceDirection // Direction is the side of the entry: DEBIT or CREDIT.
	//  NOTE - the sequence number of the entry within the transaction, starting from 1
	Sequence int64
}
