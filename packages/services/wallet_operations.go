package services

import (
	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

// NOTE - these functions only change the wallet in memory; TransactionService.Create calls them on wallets already
// locked with SELECT ... FOR UPDATE and persists them in the same database transaction

func increaseFreeBalance(wallet *domain.Wallet, amount int64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}

	wallet.FreeBalance += amount
	return nil
}

func decreaseFreeBalance(wallet *domain.Wallet, amount int64) error {
	if err := validateFreeBalance(wallet, amount); err != nil {
		return err
	}

	wallet.FreeBalance -= amount
	return nil
}

func moveFreeBalanceToReserved(wallet *domain.Wallet, amount int64) error {
	if err := validateFreeBalance(wallet, amount); err != nil {
		return err
	}

	wallet.FreeBalance -= amount
	wallet.ReservedBalance += amount
	return nil
}

func moveReservedBalanceToFree(wallet *domain.Wallet, amount int64) error {
	if err := validateReservedBalance(wallet, amount); err != nil {
		return err
	}

	wallet.ReservedBalance -= amount
	wallet.FreeBalance += amount
	return nil
}

func moveFreeBalanceToHold(wallet *domain.Wallet, amount int64) error {
	if err := validateFreeBalance(wallet, amount); err != nil {
		return err
	}

	wallet.FreeBalance -= amount
	wallet.HoldBalance += amount
	return nil
}

func increaseHoldBalance(wallet *domain.Wallet, amount int64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}

	wallet.HoldBalance += amount
	return nil
}

func decreaseHoldBalance(wallet *domain.Wallet, amount int64) error {
	if err := validateHoldBalance(wallet, amount); err != nil {
		return err
	}

	wallet.HoldBalance -= amount
	return nil
}

func moveHoldBalanceToFree(wallet *domain.Wallet, amount int64) error {
	if err := validateHoldBalance(wallet, amount); err != nil {
		return err
	}

	wallet.HoldBalance -= amount
	wallet.FreeBalance += amount
	return nil
}
