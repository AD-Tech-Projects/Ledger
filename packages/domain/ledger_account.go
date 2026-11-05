package domain

// NOTE ledger account represents an account in the accounting chart of accounts.
// It defines how financial transactions are classified and recorded in the ledger,
// including their accounting nature, normal balance direction, and associated asset.
type LedgerAccount struct {
	ID string

	Type                   LedgerAccountType
	Status                 LedgerAccountStatus
	NormalBalanceDirection LedgerBalanceDirection

	Name string
	Code string

	AssetID string
	Asset   *Asset
}

type LedgerAccountType string

const (
	LEDGER_ACCOUNT_ASSET     LedgerAccountType = "ASSET"
	LEDGER_ACCOUNT_LIABILITY LedgerAccountType = "LIABILITY"
	LEDGER_ACCOUNT_EQUITY    LedgerAccountType = "EQUITY"
	LEDGER_ACCOUNT_REVENUE   LedgerAccountType = "REVENUE"
	LEDGER_ACCOUNT_EXPENSE   LedgerAccountType = "EXPENSE"
)

type LedgerBalanceDirection string

const (
	LEDGER_ACCOUNT_DEBIT  LedgerBalanceDirection = "DEBIT"
	LEDGER_ACCOUNT_CREDIT LedgerBalanceDirection = "CREDIT"
)

type LedgerAccountStatus string

const (
	LEDGER_ACCOUNT_ACTIVE   LedgerAccountStatus = "ACTIVE"
	LEDGER_ACCOUNT_INACTIVE LedgerAccountStatus = "INACTIVE"
)
