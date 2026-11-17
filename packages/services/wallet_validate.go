package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

type walletValidationError struct {
	message string
}

func (e *walletValidationError) Error() string {
	return "WalletValidationError:" + e.message
}

func (e *walletValidationError) Unwrap() error {
	return ErrInvalidWallet
}

func validateWallet(wallet *domain.Wallet) error {
	if wallet == nil {
		return &walletValidationError{message: "wallet is required"}
	}

	if wallet.OwnerID == "" {
		return &walletValidationError{message: "owner id is required"}
	}
	if _, err := uuid.Parse(wallet.OwnerID); err != nil {
		return &walletValidationError{message: "owner id must be a valid UUID"}
	}

	if wallet.AssetID == "" {
		return &walletValidationError{message: "asset id is required"}
	}
	if _, err := uuid.Parse(wallet.AssetID); err != nil {
		return &walletValidationError{message: "asset id must be a valid UUID"}
	}

	if wallet.Status != domain.WALLET_ACTIVE && wallet.Status != domain.WALLET_INACTIVE {
		return &walletValidationError{message: "wallet status is invalid"}
	}

	accountRule, ok := walletTypeAccountRules[wallet.Type]
	if !ok {
		return &walletValidationError{message: "wallet type is invalid"}
	}

	switch {
	case accountRule == accountRequired && wallet.AccountID == nil:
		return &walletValidationError{message: "wallet type requires an external account"}
	case accountRule == accountForbidden && wallet.AccountID != nil:
		return &walletValidationError{message: "wallet type does not accept an external account"}
	}
	if wallet.AccountID != nil {
		if _, err := uuid.Parse(*wallet.AccountID); err != nil {
			return &walletValidationError{message: "external account id must be a valid UUID"}
		}
	}

	return nil
}

type accountRequirement int

const (
	accountForbidden accountRequirement = iota // The wallet must not be linked to an external account
	accountOptional                            // The wallet may be linked to an external account
	accountRequired                            // The wallet must be linked to an external account
)

// walletTypeAccountRules defines the supported wallet types and whether each one must, may or must not be linked to an external account
var walletTypeAccountRules = map[domain.WalletType]accountRequirement{
	domain.WALLET_CHECKING_CUSTODY:  accountOptional, // external account that receives funds
	domain.WALLET_CHECKING_EXTERNAL: accountRequired, // bank account outside the platform that receives the payments
	domain.WALLET_SAVINGS:           accountForbidden,
	domain.WALLET_INVESTMENT:        accountOptional, // external account that receives investment returns
	domain.WALLET_CRYPTO_CUSTODY:    accountForbidden,
	domain.WALLET_CRYPTO_EXTERNAL:   accountRequired, // crypto address outside the platform that receives the transfers
}

func validateHoldBalance(wallet *domain.Wallet, amount int64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}
	if wallet.HoldBalance < amount {
		return errors.New("insufficient hold balance")
	}

	return nil
}

func validateReservedBalance(wallet *domain.Wallet, amount int64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}
	if wallet.ReservedBalance < amount {
		return errors.New("insufficient reserved balance")
	}

	return nil
}

func validateFreeBalance(wallet *domain.Wallet, amount int64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}
	if wallet.FreeBalance < amount {
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
