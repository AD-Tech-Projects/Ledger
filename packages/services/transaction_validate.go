package services

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

type transactionValidationError struct {
	message string
}

func (e *transactionValidationError) Error() string {
	return "TransactionValidationError:" + e.message
}

func (e *transactionValidationError) Unwrap() error {
	return ErrInvalidTransaction
}

type CreateTransactionInput struct {
	Amount              int64                  // Amount in the smallest unit of the asset; for REVERSAL, zero means the full original amount
	Type                domain.TransactionType // Type of the transaction; REVERSAL is only allowed for completed transactions that can be reversed
	SourceWalletID      *string                // SourceWalletID is required for transactions that debit a wallet (WITHDRAWAL, TRANSFER, RESERVE, RELEASE, HOLD, UNHOLD, PURCHASE, FEE); optional for YIELD and SALE; forbidden for REVERSAL
	DestinationWalletID *string                // DestinationWalletID is required for transactions that credit a wallet (DEPOSIT, TRANSFER, YIELD, SALE, PURCHASE, FEE); optional for RESERVE, RELEASE, HOLD, UNHOLD; forbidden for REVERSAL
	ReversalOriginID    *string                // ReversalOriginID is required only for REVERSAL transactions and must be the ID of a completed transaction that can be reversed; forbidden for all other transaction types
	ReferenceID         *string                // ReferenceID is a user-defined identifier for the transaction
	Description         string                 // Description is a user-defined description for the transaction
}

// validateCreateTransactionInput checks the request before touching the database:
// type, amount, identifiers and which wallets each transaction type accepts
func validateCreateTransactionInput(input *CreateTransactionInput) error {
	input.SourceWalletID = normalizeOptional(input.SourceWalletID)
	input.DestinationWalletID = normalizeOptional(input.DestinationWalletID)
	input.ReversalOriginID = normalizeOptional(input.ReversalOriginID)
	input.ReferenceID = normalizeOptional(input.ReferenceID)
	input.Description = strings.TrimSpace(input.Description)

	if err := validateTransactionType(input.Type); err != nil {
		return err
	}

	if input.Type == domain.TRANSACTION_REVERSAL {
		if input.ReversalOriginID == nil {
			return &transactionValidationError{message: "reversal origin id (origin transaction id) is required"}
		}
		if err := validateUUID(*input.ReversalOriginID, "reversal origin id"); err != nil {
			return err
		}
		if input.SourceWalletID != nil || input.DestinationWalletID != nil {
			return &transactionValidationError{message: "reversal wallets are taken from the original transaction and must not be informed"}
		}
		if input.Amount < 0 {
			return &transactionValidationError{message: "transaction amount must not be negative"}
		}
		return nil
	}

	if input.ReversalOriginID != nil {
		return &transactionValidationError{message: "reversal origin transaction id is only allowed for reversals"}
	}
	if input.Amount <= 0 {
		return &transactionValidationError{message: "transaction amount must be greater than zero"}
	}
	if err := validateWalletRule(input.Type, input.SourceWalletID != nil, input.DestinationWalletID != nil); err != nil {
		return err
	}
	if input.SourceWalletID != nil {
		if err := validateUUID(*input.SourceWalletID, "source wallet id"); err != nil {
			return err
		}
	}
	if input.DestinationWalletID != nil {
		if err := validateUUID(*input.DestinationWalletID, "destination wallet id"); err != nil {
			return err
		}
	}
	if input.SourceWalletID != nil && input.DestinationWalletID != nil && *input.SourceWalletID == *input.DestinationWalletID {
		return &transactionValidationError{message: "source and destination wallets cannot be the same"}
	}

	return nil
}

// validateTransaction checks a transaction with its wallets (and reversal origin) already loaded.
func validateTransaction(transaction *domain.Transaction) error {
	if transaction == nil {
		return &transactionValidationError{message: "transaction is required"}
	}
	if err := validateTransactionType(transaction.Type); err != nil {
		return err
	}
	if transaction.Amount <= 0 {
		return &transactionValidationError{message: "transaction amount must be greater than zero"}
	}

	if transaction.Type == domain.TRANSACTION_REVERSAL {
		if err := validateReversal(transaction); err != nil {
			return err
		}
	} else if err := validateWalletRule(transaction.Type, transaction.SourceWallet != nil, transaction.DestinationWallet != nil); err != nil {
		return err
	}

	return validateTransactionWallets(transaction)
}

func validateReversal(transaction *domain.Transaction) error {
	origin := transaction.ReversalOrigin
	if origin == nil {
		return &transactionValidationError{message: "reversal origin transaction is required"}
	}
	if origin.Type == domain.TRANSACTION_REVERSAL {
		return &transactionValidationError{message: "a reversal cannot be reversed"}
	}
	if _, ok := reversalEffectTypes[origin.Type]; !ok {
		return &transactionValidationError{message: fmt.Sprintf("transaction type %q cannot be reversed", origin.Type)}
	}
	if origin.Status != domain.TRANSACTION_COMPLETED {
		return &transactionValidationError{message: "only completed transactions can be reversed"}
	}
	if transaction.Amount != origin.Amount {
		return &transactionValidationError{message: "reversal amount must be equal to the original transaction amount"}
	}

	expectedSourceID, expectedDestinationID := getReversalWalletIDs(origin)
	if !sameWalletID(expectedSourceID, transaction.SourceWallet) || !sameWalletID(expectedDestinationID, transaction.DestinationWallet) {
		return &transactionValidationError{message: "reversal wallets must match the original transaction wallets"}
	}

	return nil
}

func validateTransactionWallets(transaction *domain.Transaction) error {
	source, destination := transaction.SourceWallet, transaction.DestinationWallet
	if source == nil && destination == nil {
		return &transactionValidationError{message: "at least one wallet is required"}
	}
	if source != nil && destination != nil && source.ID == destination.ID {
		return &transactionValidationError{message: "source and destination wallets cannot be the same"}
	}

	assetID := transaction.AssetID
	for _, wallet := range []*domain.Wallet{source, destination} {
		if wallet == nil {
			continue
		}
		if wallet.Status != domain.WALLET_ACTIVE {
			return &transactionValidationError{message: fmt.Sprintf("wallet %s must be active", wallet.ID)}
		}
		if assetID == "" {
			assetID = wallet.AssetID
		}
		if wallet.AssetID != assetID {
			return &transactionValidationError{message: "transaction wallets must use the same asset"}
		}
	}

	return nil
}

func validateWalletRule(transactionType domain.TransactionType, hasSource bool, hasDestination bool) error {
	rule, ok := transactionWalletRules[transactionType]
	if !ok {
		return &transactionValidationError{message: fmt.Sprintf("unsupported transaction type %q", transactionType)}
	}
	if err := checkWalletRequirement(transactionType, "source", rule.source, hasSource); err != nil {
		return err
	}
	return checkWalletRequirement(transactionType, "destination", rule.destination, hasDestination)
}

func checkWalletRequirement(transactionType domain.TransactionType, side string, requirement walletRequirement, informed bool) error {
	switch {
	case requirement == walletRequired && !informed:
		return &transactionValidationError{message: fmt.Sprintf("%s transactions require a %s wallet", transactionType, side)}
	case requirement == walletForbidden && informed:
		return &transactionValidationError{message: fmt.Sprintf("%s transactions do not accept a %s wallet", transactionType, side)}
	default:
		return nil
	}
}

func validateTransactionType(transactionType domain.TransactionType) error {
	if transactionType == domain.TRANSACTION_REVERSAL {
		return nil
	}
	if _, ok := transactionWalletRules[transactionType]; ok {
		return nil
	}
	return &transactionValidationError{message: fmt.Sprintf("unsupported transaction type %q", transactionType)}
}

func validateUUID(value string, field string) error {
	if _, err := uuid.Parse(value); err != nil {
		return &transactionValidationError{message: field + " must be a valid UUID"}
	}
	return nil
}

func sameWalletID(expected *string, wallet *domain.Wallet) bool {
	if expected == nil || wallet == nil {
		return expected == nil && wallet == nil
	}
	return *expected == wallet.ID
}

func normalizeOptional(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
