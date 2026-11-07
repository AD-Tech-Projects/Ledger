package services

import (
	"fmt"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

type transactionValidationError struct {
	message string
}

func (e *transactionValidationError) Error() string {
	return "TransactionValidationError:" + e.message
}

func validateTransaction(transaction *domain.Transaction) error {

	if transaction == nil {
		return &transactionValidationError{message: "transaction is required"}
	}

	if invalidTypeErr := validateTransactionType(transaction); invalidTypeErr != nil {
		return invalidTypeErr
	}

	if transaction.SourceWallet == nil || transaction.DestinationWallet == nil {
		return &transactionValidationError{message: "source and destination wallets are required"}
	}
	if transaction.Amount <= 0 {
		return &transactionValidationError{message: "transaction amount must be greater than zero"}
	}
	if transaction.SourceWallet.ID == transaction.DestinationWallet.ID {
		return &transactionValidationError{message: "source and destination wallets cannot be the same"}
	}
	if transaction.SourceWallet.AssetID != transaction.DestinationWallet.AssetID {
		return &transactionValidationError{message: "source and destination wallets must use the same asset"}
	}
	if transaction.SourceWallet.Status != domain.WALLET_ACTIVE || transaction.DestinationWallet.Status != domain.WALLET_ACTIVE {
		return &transactionValidationError{message: "both source and destination wallets must be active"}
	}

	return nil
}

func validateTransactionType(transaction *domain.Transaction) error {
	switch transaction.Type {
	case domain.TRANSACTION_DEPOSIT,
		domain.TRANSACTION_WITHDRAWAL,
		domain.TRANSACTION_TRANSFER,
		domain.TRANSACTION_REVERSAL,
		domain.TRANSACTION_YIELD,
		domain.TRANSACTION_FEE,
		domain.TRANSACTION_PURCHASE,
		domain.TRANSACTION_SALE,
		domain.TRANSACTION_RESERVE,
		domain.TRANSACTION_RELEASE,
		domain.TRANSACTION_HOLD,
		domain.TRANSACTION_UNHOLD:
		return nil
	default:
		return &transactionValidationError{message: fmt.Sprintf("unsupported transaction type %q", transaction.Type)}
	}
}
