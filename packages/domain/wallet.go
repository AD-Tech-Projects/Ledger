package domain

import "errors"

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
	WALLET_CHECKING   WalletType = "CHECKING"   // Payments and transfers; requires an external account
	WALLET_SAVINGS    WalletType = "SAVINGS"    // Storage of funds; requires an external account
	WALLET_INVESTMENT WalletType = "INVESTMENT" // Investments that may generate yield
	WALLET_CRYPTO     WalletType = "CRYPTO"     // Custody of crypto assets
)

type WalletStatus string

const (
	WALLET_ACTIVE   WalletStatus = "ACTIVE"   // Can send and receive funds
	WALLET_INACTIVE WalletStatus = "INACTIVE" // Deactivated (logically deleted); cannot take part in new transactions
)

func (t WalletType) RequiresAccount() bool {
	return t == WALLET_CHECKING || t == WALLET_SAVINGS
}

func (w *Wallet) MoveFreeBalanceToHold(amount int64) error {
	if err := w.validateFreeBalance(amount); err != nil {
		return err
	}

	w.FreeBalance -= amount
	w.HoldBalance += amount
	return nil
}

func (w *Wallet) IncreaseHoldBalance(amount int64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}

	w.HoldBalance += amount
	return nil
}

func (w *Wallet) DecreaseHoldBalance(amount int64) error {
	if err := w.validateHoldBalance(amount); err != nil {
		return err
	}

	w.HoldBalance -= amount
	return nil
}

func (w *Wallet) MoveHoldBalanceToFree(amount int64) error {
	if err := w.validateHoldBalance(amount); err != nil {
		return err
	}

	w.HoldBalance -= amount
	w.FreeBalance += amount
	return nil
}

func (w *Wallet) validateHoldBalance(amount int64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}
	if w.HoldBalance < amount {
		return errors.New("insufficient hold balance")
	}

	return nil
}

func (w *Wallet) validateFreeBalance(amount int64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}
	if w.FreeBalance < amount {
		return errors.New("insufficient free balance")
	}

	return nil
}

func validateAmount(amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	return nil
}
