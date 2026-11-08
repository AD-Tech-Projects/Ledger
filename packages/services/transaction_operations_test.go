package services

import (
	"errors"
	"testing"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

func TestApplyTransactionCases(t *testing.T) {
	tests := []struct {
		name            string
		transactionType domain.TransactionType
		sourceFree      int64
		sourceHold      int64
		destFree        int64
		wantSourceFree  int64
		wantSourceHold  int64
		wantDestFree    int64
		wantDestHold    int64
	}{
		{"deposit", domain.TRANSACTION_DEPOSIT, 100, 0, 0, 50, 0, 50, 0},
		{"withdrawal", domain.TRANSACTION_WITHDRAWAL, 0, 0, 100, 50, 0, 50, 0},
		{"transfer", domain.TRANSACTION_TRANSFER, 100, 0, 0, 50, 0, 50, 0},
		{"reversal", domain.TRANSACTION_REVERSAL, 0, 0, 100, 50, 0, 50, 0},
		{"yield", domain.TRANSACTION_YIELD, 0, 0, 0, 0, 0, 50, 0},
		{"fee", domain.TRANSACTION_FEE, 100, 0, 0, 50, 0, 50, 0},
		{"purchase", domain.TRANSACTION_PURCHASE, 100, 0, 0, 50, 0, 50, 0},
		{"sale", domain.TRANSACTION_SALE, 0, 0, 0, 0, 0, 50, 0},
		{"reserve", domain.TRANSACTION_RESERVE, 100, 0, 0, 50, 50, 0, 0},
		{"release", domain.TRANSACTION_RELEASE, 0, 50, 0, 50, 0, 0, 0},
		{"hold", domain.TRANSACTION_HOLD, 100, 0, 0, 50, 50, 0, 0},
		{"unhold", domain.TRANSACTION_UNHOLD, 0, 50, 0, 50, 0, 0, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := testWallet("source", test.sourceFree, test.sourceHold)
			destination := testWallet("destination", test.destFree, 0)
			transaction := testTransaction(test.transactionType, 50, &source, &destination)

			if err := validateTransaction(transaction); err != nil {
				t.Fatal(err)
			}
			if err := applyTransaction(transaction); err != nil {
				t.Fatal(err)
			}

			assertBalances(t, "source", &source, test.wantSourceFree, test.wantSourceHold)
			assertBalances(t, "destination", &destination, test.wantDestFree, test.wantDestHold)
		})
	}
}

// NOTE - covers the original rollback scenario: an amount greater than the available balance
// must fail without changing any wallet (no partial transfer).
func TestApplyTransactionInsufficientBalanceKeepsWalletsUnchanged(t *testing.T) {
	tests := []struct {
		name            string
		transactionType domain.TransactionType
		sourceFree      int64
		sourceHold      int64
		destFree        int64
	}{
		{"deposit", domain.TRANSACTION_DEPOSIT, 100, 0, 0},
		{"transfer", domain.TRANSACTION_TRANSFER, 100, 0, 0},
		{"fee", domain.TRANSACTION_FEE, 100, 0, 0},
		{"purchase", domain.TRANSACTION_PURCHASE, 100, 0, 0},
		{"withdrawal", domain.TRANSACTION_WITHDRAWAL, 0, 0, 100},
		{"reversal", domain.TRANSACTION_REVERSAL, 0, 0, 100},
		{"reserve", domain.TRANSACTION_RESERVE, 100, 0, 0},
		{"hold", domain.TRANSACTION_HOLD, 100, 0, 0},
		{"release", domain.TRANSACTION_RELEASE, 0, 100, 0},
		{"unhold", domain.TRANSACTION_UNHOLD, 0, 100, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := testWallet("source", test.sourceFree, test.sourceHold)
			destination := testWallet("destination", test.destFree, 0)
			transaction := testTransaction(test.transactionType, 150, &source, &destination)

			if err := validateTransaction(transaction); err != nil {
				t.Fatal(err)
			}
			if err := applyTransaction(transaction); err == nil {
				t.Fatal("expected insufficient balance error")
			}

			assertBalances(t, "source", &source, test.sourceFree, test.sourceHold)
			assertBalances(t, "destination", &destination, test.destFree, 0)
		})
	}
}

// NOTE - covers the original scenarios of deposits/withdrawals without destination/source wallet:
// invalid transactions are rejected before any balance is changed.
func TestValidateTransactionRejectsInvalidTransactions(t *testing.T) {
	inactive := func(w *domain.Wallet) { w.Status = domain.WALLET_INACTIVE }
	otherAsset := func(w *domain.Wallet) { w.AssetID = "other-asset" }

	tests := []struct {
		name            string
		transactionType domain.TransactionType
		amount          int64
		withSource      bool
		withDestination bool
		sameWallet      bool
		changeSource    func(*domain.Wallet)
		changeDest      func(*domain.Wallet)
	}{
		{"deposit without destination wallet", domain.TRANSACTION_DEPOSIT, 50, true, false, false, nil, nil},
		{"withdrawal without source wallet", domain.TRANSACTION_WITHDRAWAL, 50, false, true, false, nil, nil},
		{"zero amount", domain.TRANSACTION_TRANSFER, 0, true, true, false, nil, nil},
		{"negative amount", domain.TRANSACTION_TRANSFER, -50, true, true, false, nil, nil},
		{"same source and destination wallet", domain.TRANSACTION_TRANSFER, 50, true, true, true, nil, nil},
		{"wallets with different assets", domain.TRANSACTION_TRANSFER, 50, true, true, false, nil, otherAsset},
		{"inactive source wallet", domain.TRANSACTION_TRANSFER, 50, true, true, false, inactive, nil},
		{"inactive destination wallet", domain.TRANSACTION_TRANSFER, 50, true, true, false, nil, inactive},
		{"unsupported transaction type", domain.TransactionType("UNKNOWN"), 50, true, true, false, nil, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := testWallet("source", 100, 0)
			destination := testWallet("destination", 100, 0)
			if test.changeSource != nil {
				test.changeSource(&source)
			}
			if test.changeDest != nil {
				test.changeDest(&destination)
			}

			var sourceWallet, destinationWallet *domain.Wallet
			if test.withSource {
				sourceWallet = &source
			}
			if test.withDestination {
				destinationWallet = &destination
			}
			if test.sameWallet {
				destinationWallet = sourceWallet
			}

			err := validateTransaction(testTransaction(test.transactionType, test.amount, sourceWallet, destinationWallet))
			var validationErr *transactionValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("expected transaction validation error, got %v", err)
			}

			assertBalances(t, "source", &source, 100, 0)
			assertBalances(t, "destination", &destination, 100, 0)
		})
	}
}

func TestValidateTransactionRejectsNilTransaction(t *testing.T) {
	var validationErr *transactionValidationError
	if err := validateTransaction(nil); !errors.As(err, &validationErr) {
		t.Fatalf("expected transaction validation error, got %v", err)
	}
}

func testWallet(id string, freeBalance, holdBalance int64) domain.Wallet {
	return domain.Wallet{
		ID:          id,
		Status:      domain.WALLET_ACTIVE,
		AssetID:     "asset",
		FreeBalance: freeBalance,
		HoldBalance: holdBalance,
	}
}

func testTransaction(transactionType domain.TransactionType, amount int64, source, destination *domain.Wallet) *domain.Transaction {
	return &domain.Transaction{
		Type:              transactionType,
		Amount:            amount,
		SourceWallet:      source,
		DestinationWallet: destination,
	}
}

func assertBalances(t *testing.T, label string, wallet *domain.Wallet, wantFree, wantHold int64) {
	t.Helper()

	if wallet.FreeBalance != wantFree || wallet.HoldBalance != wantHold {
		t.Fatalf("unexpected %s balances: free=%d hold=%d, want free=%d hold=%d",
			label, wallet.FreeBalance, wallet.HoldBalance, wantFree, wantHold)
	}
}
