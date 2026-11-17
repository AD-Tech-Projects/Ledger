package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

// # Main Flux
func TestWalletTypeAccountRulesCases(t *testing.T) {
	tests := map[domain.WalletType]accountRequirement{
		domain.WALLET_CHECKING_CUSTODY:  accountOptional,
		domain.WALLET_CHECKING_EXTERNAL: accountRequired,
		domain.WALLET_SAVINGS:           accountForbidden,
		domain.WALLET_INVESTMENT:        accountOptional,
		domain.WALLET_CRYPTO_CUSTODY:    accountForbidden,
		domain.WALLET_CRYPTO_EXTERNAL:   accountRequired,
	}

	if len(walletTypeAccountRules) != len(tests) {
		t.Fatalf("walletTypeAccountRules has %d wallet types, want %d", len(walletTypeAccountRules), len(tests))
	}
	for walletType, want := range tests {
		got, ok := walletTypeAccountRules[walletType]
		if !ok {
			t.Fatalf("walletTypeAccountRules has no rule for %s", walletType)
		}
		if got != want {
			t.Fatalf("walletTypeAccountRules[%s] = %v, want %v", walletType, got, want)
		}
	}
}

func TestValidateWalletExternalAccountCases(t *testing.T) {
	accountID := uuid.NewString()

	tests := []struct {
		name       string
		walletType domain.WalletType
		accountID  *string
	}{
		{"checking custody without external account", domain.WALLET_CHECKING_CUSTODY, nil},
		{"checking custody with external account", domain.WALLET_CHECKING_CUSTODY, &accountID},
		{"checking external with external account", domain.WALLET_CHECKING_EXTERNAL, &accountID},
		{"savings without external account", domain.WALLET_SAVINGS, nil},
		{"investment without external account", domain.WALLET_INVESTMENT, nil},
		{"investment with external account", domain.WALLET_INVESTMENT, &accountID},
		{"crypto custody without external account", domain.WALLET_CRYPTO_CUSTODY, nil},
		{"crypto external with external account", domain.WALLET_CRYPTO_EXTERNAL, &accountID},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := newTestValidWallet(test.walletType, test.accountID)

			if err := validateWallet(wallet); err != nil {
				t.Fatalf("unexpected error for %s: %v", test.walletType, err)
			}
			assertWalletAccountUnchanged(t, wallet, test.accountID)
		})
	}
}

// # Failures
func TestValidateWalletRejectsInvalidExternalAccount(t *testing.T) {
	accountID := uuid.NewString()
	invalidAccountID := "not-a-uuid"

	tests := []struct {
		name       string
		walletType domain.WalletType
		accountID  *string
	}{
		{"checking external without external account", domain.WALLET_CHECKING_EXTERNAL, nil},
		{"crypto external without external account", domain.WALLET_CRYPTO_EXTERNAL, nil},
		{"savings with external account", domain.WALLET_SAVINGS, &accountID},
		{"crypto custody with external account", domain.WALLET_CRYPTO_CUSTODY, &accountID},
		{"external account id is not a UUID", domain.WALLET_CHECKING_EXTERNAL, &invalidAccountID},
		{"optional external account id is not a UUID", domain.WALLET_CHECKING_CUSTODY, &invalidAccountID},
		{"removed wallet type", domain.WalletType("CHECKING"), nil},
		{"empty wallet type", domain.WalletType(""), nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := newTestValidWallet(test.walletType, test.accountID)

			err := validateWallet(wallet)
			if !errors.Is(err, ErrInvalidWallet) {
				t.Fatalf("expected ErrInvalidWallet for %s with external account %v, got %v", test.walletType, test.accountID, err)
			}
			assertWalletAccountUnchanged(t, wallet, test.accountID)
		})
	}
}

// helpers
func newTestValidWallet(walletType domain.WalletType, accountID *string) *domain.Wallet {
	return &domain.Wallet{
		Type:      walletType,
		Status:    domain.WALLET_ACTIVE,
		OwnerID:   uuid.NewString(),
		AssetID:   uuid.NewString(),
		AccountID: accountID,
	}
}

func assertWalletAccountUnchanged(t *testing.T, wallet *domain.Wallet, want *string) {
	t.Helper()
	if wallet.AccountID != want {
		t.Fatalf("wallet external account changed: got %v, want %v", wallet.AccountID, want)
	}
}
