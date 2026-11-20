package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

// # Main Flux
func TestValidateAccountCases(t *testing.T) {
	tests := []struct {
		name   string
		change func(*domain.Account)
		kind   domain.AccountKind
	}{
		{"bank account with account number", nil, domain.ACCOUNT_BANK},
		{"bank account with IBAN only", func(a *domain.Account) { a.AccountNumber, a.IBAN = "", "DE89370400440532013000" }, domain.ACCOUNT_BANK},
		{"bank account of an unknown country", func(a *domain.Account) { a.CountryCode = domain.COUNTRY_UNKNOWN }, domain.ACCOUNT_BANK},
		{"crypto account", nil, domain.ACCOUNT_CRYPTO},
		{"crypto account with memo", func(a *domain.Account) { a.CryptoMemo = "123456" }, domain.ACCOUNT_CRYPTO},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account := newTestAccount(test.kind)
			if test.change != nil {
				test.change(account)
			}
			want := *account

			if err := validateAccount(account); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertAccountUnchanged(t, account, &want)
		})
	}
}

// # Failures
func TestValidateAccountRejectsInvalidAccounts(t *testing.T) {
	tests := []struct {
		name   string
		kind   domain.AccountKind
		change func(*domain.Account)
	}{
		{"empty owner id", domain.ACCOUNT_BANK, func(a *domain.Account) { a.OwnerID = "" }},
		{"invalid owner id", domain.ACCOUNT_BANK, func(a *domain.Account) { a.OwnerID = "not-a-uuid" }},
		{"unsupported country", domain.ACCOUNT_BANK, func(a *domain.Account) { a.CountryCode = "XX" }},
		{"empty kind", "", nil},
		{"invalid kind", "PAYPAL", nil},
		{"bank account without currency", domain.ACCOUNT_BANK, func(a *domain.Account) { a.Currency = " " }},
		{"bank account without account number and IBAN", domain.ACCOUNT_BANK, func(a *domain.Account) { a.AccountNumber, a.IBAN = "", "" }},
		{"bank account with crypto address", domain.ACCOUNT_BANK, func(a *domain.Account) { a.CryptoAddress = "bc1qxyz" }},
		{"bank account with crypto memo", domain.ACCOUNT_BANK, func(a *domain.Account) { a.CryptoMemo = "123" }},
		{"crypto account without network", domain.ACCOUNT_CRYPTO, func(a *domain.Account) { a.CryptoNetwork = "" }},
		{"crypto account without address", domain.ACCOUNT_CRYPTO, func(a *domain.Account) { a.CryptoAddress = "" }},
		{"crypto account with currency", domain.ACCOUNT_CRYPTO, func(a *domain.Account) { a.Currency = "BRL" }},
		{"crypto account with IBAN", domain.ACCOUNT_CRYPTO, func(a *domain.Account) { a.IBAN = "DE89370400440532013000" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account := newTestAccount(test.kind)
			if test.change != nil {
				test.change(account)
			}
			want := *account

			err := validateAccount(account)
			var validationErr *accountValidationError
			if !errors.As(err, &validationErr) || !errors.Is(err, ErrInvalidAccount) {
				t.Fatalf("expected external account validation error, got %v", err)
			}
			assertAccountUnchanged(t, account, &want)
		})
	}
}

func TestValidateAccountRejectsNilAccount(t *testing.T) {
	if err := validateAccount(nil); !errors.Is(err, ErrInvalidAccount) {
		t.Fatalf("expected invalid external account error, got %v", err)
	}
}

// helpers
func newTestAccount(kind domain.AccountKind) *domain.Account {
	account := &domain.Account{
		ID:          uuid.NewString(),
		OwnerID:     uuid.NewString(),
		Kind:        kind,
		CountryCode: domain.COUNTRY_BR,
	}
	switch kind {
	case domain.ACCOUNT_BANK:
		account.Currency = "BRL"
		account.InstitutionName = "Bank"
		account.BranchCode = "0001"
		account.AccountNumber = "123456-7"
	case domain.ACCOUNT_CRYPTO:
		account.CountryCode = domain.COUNTRY_UNKNOWN
		account.CryptoNetwork = "BITCOIN"
		account.CryptoAddress = "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq"
	}
	return account
}

func assertAccountUnchanged(t *testing.T, account *domain.Account, want *domain.Account) {
	t.Helper()
	if account.Kind != want.Kind || account.Currency != want.Currency || account.AccountNumber != want.AccountNumber ||
		account.IBAN != want.IBAN || account.CryptoNetwork != want.CryptoNetwork || account.CryptoAddress != want.CryptoAddress {
		t.Fatalf("external account changed: got %+v, want %+v", account, want)
	}
}
