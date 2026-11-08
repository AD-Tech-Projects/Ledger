package domain

// NOTE ledger account represents an account in the accounting chart of accounts.
// It defines how financial transactions are classified and recorded in the ledger,
// including their accounting nature, normal balance direction, and associated asset.
type LedgerAccount struct {
	ID string

	Type                   LedgerAccountType
	Status                 LedgerAccountStatus
	NormalBalanceDirection LedgerBalanceDirection // NormalBalanceDirection is the side (DEBIT or CREDIT) that increases the account balance.

	Name string // Name is the human readable description (e.g. "Users available balance").
	Code string // Code is the unique code of the account in the chart of accounts (e.g. "2001").

	AssetID string // AssetID is the asset recorded by the account.
	Asset   *Asset
}

type LedgerAccountType string

const (
	LEDGER_ACCOUNT_ASSET     LedgerAccountType = "ASSET"     // What the platform owns (e.g. bank reserves); normal balance DEBIT
	LEDGER_ACCOUNT_LIABILITY LedgerAccountType = "LIABILITY" // What the platform owes (e.g. users' balances); normal balance CREDIT
	LEDGER_ACCOUNT_EQUITY    LedgerAccountType = "EQUITY"    // Owners' capital; normal balance CREDIT
	LEDGER_ACCOUNT_REVENUE   LedgerAccountType = "REVENUE"   // Income (e.g. fees charged); normal balance CREDIT
	LEDGER_ACCOUNT_EXPENSE   LedgerAccountType = "EXPENSE"   // Costs (e.g. yield paid to users); normal balance DEBIT
)

type LedgerBalanceDirection string

const (
	LEDGER_ACCOUNT_DEBIT  LedgerBalanceDirection = "DEBIT"  // Debit side: increases ASSET/EXPENSE, decreases LIABILITY/EQUITY/REVENUE
	LEDGER_ACCOUNT_CREDIT LedgerBalanceDirection = "CREDIT" // Credit side: increases LIABILITY/EQUITY/REVENUE, decreases ASSET/EXPENSE
)

type LedgerAccountStatus string

const (
	LEDGER_ACCOUNT_ACTIVE   LedgerAccountStatus = "ACTIVE"   // Can receive new entries
	LEDGER_ACCOUNT_INACTIVE LedgerAccountStatus = "INACTIVE" // Cannot receive new entries; history is preserved
)
