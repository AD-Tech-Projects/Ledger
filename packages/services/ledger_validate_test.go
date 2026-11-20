package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

const (
	testLedgerAssetID       = "44444444-4444-4444-4444-444444444444"
	testLedgerOtherAssetID  = "55555555-5555-5555-5555-555555555555"
	testLedgerTransactionID = "66666666-6666-6666-6666-666666666666"
)

// # Main Flux
func TestLedgerAccountNormalDirectionsCases(t *testing.T) {
	tests := map[domain.LedgerAccountType]domain.LedgerBalanceDirection{
		domain.LEDGER_ACCOUNT_ASSET:     domain.LEDGER_ACCOUNT_DEBIT,
		domain.LEDGER_ACCOUNT_EXPENSE:   domain.LEDGER_ACCOUNT_DEBIT,
		domain.LEDGER_ACCOUNT_LIABILITY: domain.LEDGER_ACCOUNT_CREDIT,
		domain.LEDGER_ACCOUNT_EQUITY:    domain.LEDGER_ACCOUNT_CREDIT,
		domain.LEDGER_ACCOUNT_REVENUE:   domain.LEDGER_ACCOUNT_CREDIT,
	}

	if len(ledgerAccountNormalDirections) != len(tests) {
		t.Fatalf("ledgerAccountNormalDirections has %d types, want %d", len(ledgerAccountNormalDirections), len(tests))
	}
	for accountType, want := range tests {
		t.Run(string(accountType), func(t *testing.T) {
			account := newTestLedgerAccount(accountType, want, domain.LEDGER_ACCOUNT_ACTIVE)

			if err := validateLedgerAccount(account); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if account.NormalBalanceDirection != want {
				t.Fatalf("normal balance direction = %s, want %s", account.NormalBalanceDirection, want)
			}
		})
	}
}

func TestValidateLedgerAccountAcceptsInactiveAccount(t *testing.T) {
	account := newTestLedgerAccount(domain.LEDGER_ACCOUNT_LIABILITY, domain.LEDGER_ACCOUNT_CREDIT, domain.LEDGER_ACCOUNT_INACTIVE)

	if err := validateLedgerAccount(account); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if account.Status != domain.LEDGER_ACCOUNT_INACTIVE {
		t.Fatalf("ledger account status = %s, want %s", account.Status, domain.LEDGER_ACCOUNT_INACTIVE)
	}
}

func TestValidateLedgerEntriesCases(t *testing.T) {
	tests := []struct {
		name    string
		entries func() []*domain.LedgerEntry
	}{
		{"debit and credit in sequence", func() []*domain.LedgerEntry {
			return []*domain.LedgerEntry{newTestLedgerEntry(1, domain.LEDGER_ACCOUNT_DEBIT), newTestLedgerEntry(2, domain.LEDGER_ACCOUNT_CREDIT)}
		}},
		{"entries out of order", func() []*domain.LedgerEntry {
			return []*domain.LedgerEntry{
				newTestLedgerEntry(3, domain.LEDGER_ACCOUNT_CREDIT),
				newTestLedgerEntry(1, domain.LEDGER_ACCOUNT_DEBIT),
				newTestLedgerEntry(2, domain.LEDGER_ACCOUNT_CREDIT),
			}
		}},
		{"single entry", func() []*domain.LedgerEntry {
			return []*domain.LedgerEntry{newTestLedgerEntry(1, domain.LEDGER_ACCOUNT_DEBIT)}
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entries := test.entries()
			wantSequences := ledgerEntrySequences(entries)

			if err := validateLedgerEntries(entries); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertLedgerEntrySequences(t, entries, wantSequences)
		})
	}
}

// # Failures
func TestValidateLedgerAccountRejectsInvalidAccounts(t *testing.T) {
	tests := []struct {
		name   string
		change func(*domain.LedgerAccount)
	}{
		{"asset with credit direction", func(a *domain.LedgerAccount) { a.NormalBalanceDirection = domain.LEDGER_ACCOUNT_CREDIT }},
		{"expense with credit direction", func(a *domain.LedgerAccount) {
			a.Type, a.NormalBalanceDirection = domain.LEDGER_ACCOUNT_EXPENSE, domain.LEDGER_ACCOUNT_CREDIT
		}},
		{"liability with debit direction", func(a *domain.LedgerAccount) {
			a.Type, a.NormalBalanceDirection = domain.LEDGER_ACCOUNT_LIABILITY, domain.LEDGER_ACCOUNT_DEBIT
		}},
		{"equity with debit direction", func(a *domain.LedgerAccount) {
			a.Type, a.NormalBalanceDirection = domain.LEDGER_ACCOUNT_EQUITY, domain.LEDGER_ACCOUNT_DEBIT
		}},
		{"revenue with debit direction", func(a *domain.LedgerAccount) {
			a.Type, a.NormalBalanceDirection = domain.LEDGER_ACCOUNT_REVENUE, domain.LEDGER_ACCOUNT_DEBIT
		}},
		{"empty direction", func(a *domain.LedgerAccount) { a.NormalBalanceDirection = "" }},
		{"empty type", func(a *domain.LedgerAccount) { a.Type = "" }},
		{"invalid type", func(a *domain.LedgerAccount) { a.Type = "INCOME" }},
		{"invalid status", func(a *domain.LedgerAccount) { a.Status = "CLOSED" }},
		{"empty code", func(a *domain.LedgerAccount) { a.Code = " " }},
		{"empty name", func(a *domain.LedgerAccount) { a.Name = "" }},
		{"invalid asset id", func(a *domain.LedgerAccount) { a.AssetID = "BRL" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account := newTestLedgerAccount(domain.LEDGER_ACCOUNT_ASSET, domain.LEDGER_ACCOUNT_DEBIT, domain.LEDGER_ACCOUNT_ACTIVE)
			test.change(account)
			want := *account

			err := validateLedgerAccount(account)
			var validationErr *ledgerAccountValidationError
			if !errors.As(err, &validationErr) || !errors.Is(err, ErrInvalidLedgerAccount) {
				t.Fatalf("expected ledger account validation error, got %v", err)
			}
			if *account != want {
				t.Fatalf("ledger account changed: got %+v, want %+v", account, &want)
			}
		})
	}
}

func TestValidateLedgerAccountRejectsNilAccount(t *testing.T) {
	if err := validateLedgerAccount(nil); !errors.Is(err, ErrInvalidLedgerAccount) {
		t.Fatalf("expected invalid ledger account error, got %v", err)
	}
}

func TestValidateLedgerEntriesRejectsInvalidEntries(t *testing.T) {
	balancedEntries := func(change func(entries []*domain.LedgerEntry)) func() []*domain.LedgerEntry {
		return func() []*domain.LedgerEntry {
			entries := []*domain.LedgerEntry{newTestLedgerEntry(1, domain.LEDGER_ACCOUNT_DEBIT), newTestLedgerEntry(2, domain.LEDGER_ACCOUNT_CREDIT)}
			change(entries)
			return entries
		}
	}

	tests := []struct {
		name    string
		entries func() []*domain.LedgerEntry
	}{
		{"nil entries", func() []*domain.LedgerEntry { return nil }},
		{"empty entries", func() []*domain.LedgerEntry { return []*domain.LedgerEntry{} }},
		{"nil entry", balancedEntries(func(e []*domain.LedgerEntry) { e[1] = nil })},
		{"nil first entry", balancedEntries(func(e []*domain.LedgerEntry) { e[0] = nil })},
		{"entries of different ledger transactions", balancedEntries(func(e []*domain.LedgerEntry) { e[1].LedgerTransactionID = uuid.NewString() })},
		{"zero amount", balancedEntries(func(e []*domain.LedgerEntry) { e[0].Amount = 0 })},
		{"negative amount", balancedEntries(func(e []*domain.LedgerEntry) { e[1].Amount = -10 })},
		{"invalid direction", balancedEntries(func(e []*domain.LedgerEntry) { e[0].Direction = "BOTH" })},
		{"ledger account not loaded", balancedEntries(func(e []*domain.LedgerEntry) { e[0].LedgerAccount = nil })},
		{"loaded ledger account is another account", balancedEntries(func(e []*domain.LedgerEntry) { e[0].LedgerAccountID = uuid.NewString() })},
		{"inactive ledger account", balancedEntries(func(e []*domain.LedgerEntry) { e[1].LedgerAccount.Status = domain.LEDGER_ACCOUNT_INACTIVE })},
		{"entry asset differs from the account asset", balancedEntries(func(e []*domain.LedgerEntry) { e[0].AssetID = testLedgerOtherAssetID })},
		{"account asset differs from the entry asset", balancedEntries(func(e []*domain.LedgerEntry) { e[1].LedgerAccount.AssetID = testLedgerOtherAssetID })},
		{"sequence starting at zero", balancedEntries(func(e []*domain.LedgerEntry) { e[0].Sequence, e[1].Sequence = 0, 1 })},
		{"sequence starting at two", balancedEntries(func(e []*domain.LedgerEntry) { e[0].Sequence, e[1].Sequence = 2, 3 })},
		{"negative sequence", balancedEntries(func(e []*domain.LedgerEntry) { e[0].Sequence = -1 })},
		{"sequence with a gap", balancedEntries(func(e []*domain.LedgerEntry) { e[1].Sequence = 3 })},
		{"repeated sequence", balancedEntries(func(e []*domain.LedgerEntry) { e[1].Sequence = 1 })},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entries := test.entries()
			wantSequences := ledgerEntrySequences(entries)

			err := validateLedgerEntries(entries)
			var validationErr *ledgerEntryValidationError
			if !errors.As(err, &validationErr) || !errors.Is(err, ErrInvalidLedgerEntry) {
				t.Fatalf("expected ledger entry validation error, got %v", err)
			}
			assertLedgerEntrySequences(t, entries, wantSequences)
		})
	}
}

// helpers
func newTestLedgerAccount(
	accountType domain.LedgerAccountType,
	direction domain.LedgerBalanceDirection,
	status domain.LedgerAccountStatus,
) *domain.LedgerAccount {
	return &domain.LedgerAccount{
		ID:                     uuid.NewString(),
		Type:                   accountType,
		Status:                 status,
		NormalBalanceDirection: direction,
		Name:                   "Test account",
		Code:                   "1001-BRL",
		AssetID:                testLedgerAssetID,
	}
}

func newTestLedgerEntry(sequence int64, direction domain.LedgerBalanceDirection) *domain.LedgerEntry {
	account := newTestLedgerAccount(domain.LEDGER_ACCOUNT_LIABILITY, domain.LEDGER_ACCOUNT_CREDIT, domain.LEDGER_ACCOUNT_ACTIVE)
	return &domain.LedgerEntry{
		ID:                  uuid.NewString(),
		LedgerTransactionID: testLedgerTransactionID,
		LedgerAccountID:     account.ID,
		LedgerAccount:       account,
		AssetID:             testLedgerAssetID,
		Amount:              100,
		Direction:           direction,
		Sequence:            sequence,
	}
}

// ledgerEntrySequences returns the sequence of each entry (-1 for nil entries), to check that validation keeps them
func ledgerEntrySequences(entries []*domain.LedgerEntry) []int64 {
	sequences := make([]int64, len(entries))
	for index, entry := range entries {
		sequences[index] = -1
		if entry != nil {
			sequences[index] = entry.Sequence
		}
	}
	return sequences
}

func assertLedgerEntrySequences(t *testing.T, entries []*domain.LedgerEntry, want []int64) {
	t.Helper()
	got := ledgerEntrySequences(entries)
	if len(got) != len(want) {
		t.Fatalf("ledger entries = %d, want %d", len(got), len(want))
	}
	for index := range got {
		if got[index] != want[index] {
			t.Fatalf("ledger entry sequences = %v, want %v", got, want)
		}
	}
}
