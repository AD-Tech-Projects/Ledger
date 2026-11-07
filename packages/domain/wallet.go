package domain

import "errors"

type Wallet struct {
	ID string

	Type    WalletType
	Status  WalletStatus
	Version int64

	// NOTE - store the balance in the smallest unit of the currency (e.g., cents for USD)
	FreeBalance     int64
	ReservedBalance int64
	// NOTE - balance that is on hold, such as for pending transactions or disputes
	HoldBalance int64

	OwnerID string
	Owner   *Owner

	AssetID string
	Asset   *Asset

	AccountID *string
	Account   *Account
}

type WalletType string

const (
	WALLET_CHECKING   WalletType = "CHECKING"
	WALLET_SAVINGS    WalletType = "SAVINGS"
	WALLET_INVESTMENT WalletType = "INVESTMENT"
	WALLET_CRYPTO     WalletType = "CRYPTO"
)

type WalletStatus string

const (
	WALLET_ACTIVE   WalletStatus = "ACTIVE"
	WALLET_INACTIVE WalletStatus = "INACTIVE"
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
