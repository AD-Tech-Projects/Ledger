package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

var ErrInvalidAccount = errors.New("invalid external account")

type accountValidationError struct {
	message string
}

func (e *accountValidationError) Error() string {
	return "AccountValidationError:" + e.message
}

func (e *accountValidationError) Unwrap() error {
	return ErrInvalidAccount
}

// validateAccount checks an external account: its owner, its country and the fields of its kind. A BANK account needs
// the currency and the account number or IBAN, a CRYPTO account needs the network and the address, and neither kind
// accepts the fields of the other.
func validateAccount(account *domain.Account) error {
	if account == nil {
		return &accountValidationError{message: "external account is required"}
	}
	if _, err := uuid.Parse(account.OwnerID); err != nil {
		return &accountValidationError{message: "owner id must be a valid UUID"}
	}
	isKnownCountry := account.CountryCode == domain.COUNTRY_BR || account.CountryCode == domain.COUNTRY_US ||
		account.CountryCode == domain.COUNTRY_EU || account.CountryCode == domain.COUNTRY_UNKNOWN
	if !isKnownCountry {
		return &accountValidationError{message: fmt.Sprintf("external account country code %q is not supported", account.CountryCode)}
	}

	hasBankFields := hasAnyValue(account.Currency, account.InstitutionName, account.BranchCode, account.AccountNumber,
		account.RoutingCode, account.IBAN, account.BIC)
	hasCryptoFields := hasAnyValue(account.CryptoNetwork, account.CryptoAddress, account.CryptoMemo)

	switch account.Kind {
	case domain.ACCOUNT_BANK:
		if hasCryptoFields {
			return &accountValidationError{message: "bank accounts do not accept crypto fields"}
		}
		if !hasAnyValue(account.Currency) {
			return &accountValidationError{message: "bank accounts require a currency"}
		}
		if !hasAnyValue(account.AccountNumber, account.IBAN) {
			return &accountValidationError{message: "bank accounts require an account number or an IBAN"}
		}
	case domain.ACCOUNT_CRYPTO:
		if hasBankFields {
			return &accountValidationError{message: "crypto accounts do not accept bank fields"}
		}
		isAddressInformed := hasAnyValue(account.CryptoNetwork) && hasAnyValue(account.CryptoAddress)
		if !isAddressInformed {
			return &accountValidationError{message: "crypto accounts require a network and an address"}
		}
	default:
		return &accountValidationError{message: fmt.Sprintf("external account kind %q is invalid", account.Kind)}
	}

	return nil
}

func hasAnyValue(values ...string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return true
		}
	}
	return false
}
