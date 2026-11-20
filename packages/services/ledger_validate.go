package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

var (
	ErrInvalidLedgerAccount = errors.New("invalid ledger account")
	ErrInvalidLedgerEntry   = errors.New("invalid ledger entry")
)

type ledgerAccountValidationError struct {
	message string
}

func (e *ledgerAccountValidationError) Error() string {
	return "LedgerAccountValidationError:" + e.message
}

func (e *ledgerAccountValidationError) Unwrap() error {
	return ErrInvalidLedgerAccount
}

type ledgerEntryValidationError struct {
	message string
}

func (e *ledgerEntryValidationError) Error() string {
	return "LedgerEntryValidationError:" + e.message
}

func (e *ledgerEntryValidationError) Unwrap() error {
	return ErrInvalidLedgerEntry
}

// ledgerAccountNormalDirections defines the supported ledger account types and the side that increases each one
// (left side of the accounting equation by DEBIT, right side by CREDIT)
var ledgerAccountNormalDirections = map[domain.LedgerAccountType]domain.LedgerBalanceDirection{
	domain.LEDGER_ACCOUNT_ASSET:     domain.LEDGER_ACCOUNT_DEBIT,
	domain.LEDGER_ACCOUNT_EXPENSE:   domain.LEDGER_ACCOUNT_DEBIT,
	domain.LEDGER_ACCOUNT_LIABILITY: domain.LEDGER_ACCOUNT_CREDIT,
	domain.LEDGER_ACCOUNT_EQUITY:    domain.LEDGER_ACCOUNT_CREDIT,
	domain.LEDGER_ACCOUNT_REVENUE:   domain.LEDGER_ACCOUNT_CREDIT,
}

// validateLedgerAccount checks an account of the chart of accounts, including that its normal balance direction is the
// one of its type
func validateLedgerAccount(account *domain.LedgerAccount) error {
	if account == nil {
		return &ledgerAccountValidationError{message: "ledger account is required"}
	}

	normalDirection, isKnownType := ledgerAccountNormalDirections[account.Type]
	if !isKnownType {
		return &ledgerAccountValidationError{message: fmt.Sprintf("ledger account type %q is invalid", account.Type)}
	}
	if account.NormalBalanceDirection != normalDirection {
		return &ledgerAccountValidationError{message: fmt.Sprintf(
			"normal balance direction of %s accounts must be %s, got %q", account.Type, normalDirection, account.NormalBalanceDirection)}
	}

	isKnownStatus := account.Status == domain.LEDGER_ACCOUNT_ACTIVE || account.Status == domain.LEDGER_ACCOUNT_INACTIVE
	if !isKnownStatus {
		return &ledgerAccountValidationError{message: fmt.Sprintf("ledger account status %q is invalid", account.Status)}
	}
	if strings.TrimSpace(account.Code) == "" || strings.TrimSpace(account.Name) == "" {
		return &ledgerAccountValidationError{message: "ledger account code and name are required"}
	}
	if _, err := uuid.Parse(account.AssetID); err != nil {
		return &ledgerAccountValidationError{message: "asset id must be a valid UUID"}
	}

	return nil
}

// validateLedgerEntries checks the entries of one ledger transaction, each one with its ledger account loaded: positive
// amount, valid direction, an ACTIVE account of the same asset as the entry, and sequences from 1 to the number of entries
// without gaps or repetitions (in any order)
// NOTE - the balance of debits and credits per asset is not checked here
func validateLedgerEntries(entries []*domain.LedgerEntry) error {
	if len(entries) == 0 {
		return &ledgerEntryValidationError{message: "at least one ledger entry is required"}
	}

	ledgerTransactionID := ""
	if entries[0] != nil {
		ledgerTransactionID = entries[0].LedgerTransactionID
	}
	entryCount := int64(len(entries))
	usedSequences := make(map[int64]bool, len(entries))

	for _, entry := range entries {
		if entry == nil {
			return &ledgerEntryValidationError{message: "ledger entry is required"}
		}
		if entry.LedgerTransactionID != ledgerTransactionID {
			return &ledgerEntryValidationError{message: "ledger entries must belong to the same ledger transaction"}
		}
		if entry.Amount <= 0 {
			return &ledgerEntryValidationError{message: fmt.Sprintf("amount of entry %d must be greater than zero", entry.Sequence)}
		}
		isKnownDirection := entry.Direction == domain.LEDGER_ACCOUNT_DEBIT || entry.Direction == domain.LEDGER_ACCOUNT_CREDIT
		if !isKnownDirection {
			return &ledgerEntryValidationError{message: fmt.Sprintf("direction %q of entry %d is invalid", entry.Direction, entry.Sequence)}
		}

		account := entry.LedgerAccount
		isAccountLoaded := account != nil && account.ID == entry.LedgerAccountID
		if !isAccountLoaded {
			return &ledgerEntryValidationError{message: fmt.Sprintf("ledger account of entry %d must be loaded", entry.Sequence)}
		}
		if account.Status != domain.LEDGER_ACCOUNT_ACTIVE {
			return &ledgerEntryValidationError{message: fmt.Sprintf("ledger account %s is not active", account.Code)}
		}
		if entry.AssetID != account.AssetID {
			return &ledgerEntryValidationError{message: fmt.Sprintf("asset of entry %d must be the asset of ledger account %s", entry.Sequence, account.Code)}
		}

		isSequenceInRange := entry.Sequence >= 1 && entry.Sequence <= entryCount
		if !isSequenceInRange || usedSequences[entry.Sequence] {
			return &ledgerEntryValidationError{message: fmt.Sprintf(
				"ledger entry sequences must go from 1 to %d without gaps or repetitions, got %d", entryCount, entry.Sequence)}
		}
		usedSequences[entry.Sequence] = true
	}

	return nil
}
