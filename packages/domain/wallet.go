package domain

type Wallet struct {
	ID string

	Type    WalletType   // Type is the purpose of the wallet (what it is used for), not the asset it stores.
	Status  WalletStatus // Status defines whether the wallet can take part in new transactions.
	Version int64        // Version is reserved for optimistic locking; it is not enforced yet.

	// NOTE - store the balance in the smallest unit of the currency (e.g., cents for USD)
	FreeBalance     int64 // FreeBalance is the amount available to be moved.
	ReservedBalance int64 // ReservedBalance is the amount committed to an operation in progress, such as a scheduled transaction.
	// NOTE - balance that is on hold, such as for pending transactions or disputes
	// It is also the intermediate state of a transfer between FreeBalance of the source and of the destination.
	HoldBalance int64

	OwnerID string
	Owner   *Owner

	AssetID string // AssetID is the single asset held by the wallet.
	Asset   *Asset

	AccountID *string // AccountID is the optional external account used for deposits and withdrawals.
	Account   *Account
}

type WalletType string

const (
	WALLET_CHECKING_CUSTODY  WalletType = "CHECKING_CUSTODY"  // Deposits, withdrawals and internal transfers, without yield (like a checking account); may have an external account to receive funds
	WALLET_CHECKING_EXTERNAL WalletType = "CHECKING_EXTERNAL" // Payments to bank accounts outside the platform; requires an external account
	WALLET_SAVINGS           WalletType = "SAVINGS"           // Funds that may earn yield, such as the CDI (like a savings account); has no external account
	WALLET_INVESTMENT        WalletType = "INVESTMENT"        // Funds allocated to investments, locked until their redemption or maturity
	WALLET_CRYPTO_CUSTODY    WalletType = "CRYPTO_CUSTODY"    // Custody of a crypto asset inside the platform; has no external account and no yield
	WALLET_CRYPTO_EXTERNAL   WalletType = "CRYPTO_EXTERNAL"   // Transfers to a crypto wallet outside the platform; requires an external account (crypto address)
)

type WalletStatus string

const (
	WALLET_ACTIVE   WalletStatus = "ACTIVE"   // Can send and receive funds
	WALLET_INACTIVE WalletStatus = "INACTIVE" // Deactivated (logically deleted); cannot take part in new transactions
)
