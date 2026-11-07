package services

import (
	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

type walletValidationError struct {
	message string
}

func (e *walletValidationError) Error() string {
	return "WalletValidationError:" + e.message
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

	switch wallet.Type {
	case domain.WALLET_CHECKING, domain.WALLET_SAVINGS, domain.WALLET_INVESTMENT, domain.WALLET_CRYPTO:
	default:
		return &walletValidationError{message: "wallet type is invalid"}
	}

	if wallet.Type.RequiresAccount() {
		if wallet.AccountID == nil {
			return &walletValidationError{message: "wallet type requires an external account"}
		}
		if _, err := uuid.Parse(*wallet.AccountID); err != nil {
			return &walletValidationError{message: "external account id must be a valid UUID"}
		}
	}

	return nil
}
