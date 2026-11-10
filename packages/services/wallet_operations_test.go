package services

import (
	"testing"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

func TestWalletBalanceOperations(t *testing.T) {
	tests := []struct {
		name      string
		free      int64
		hold      int64
		operation func(*domain.Wallet) error
		wantErr   bool
		wantFree  int64
		wantHold  int64
	}{
		{"increase free", 10, 0, func(w *domain.Wallet) error { return increaseFreeBalance(w, 40) }, false, 50, 0},
		{"decrease free", 100, 0, func(w *domain.Wallet) error { return decreaseFreeBalance(w, 40) }, false, 60, 0},
		{"decrease free with insufficient free balance", 30, 0, func(w *domain.Wallet) error { return decreaseFreeBalance(w, 40) }, true, 30, 0},
		{"increase free with zero amount", 10, 0, func(w *domain.Wallet) error { return increaseFreeBalance(w, 0) }, true, 10, 0},
		{"move free to hold", 100, 0, func(w *domain.Wallet) error { return moveFreeBalanceToHold(w, 40) }, false, 60, 40},
		{"move free to hold with insufficient free balance", 30, 0, func(w *domain.Wallet) error { return moveFreeBalanceToHold(w, 40) }, true, 30, 0},
		{"move hold to free", 0, 40, func(w *domain.Wallet) error { return moveHoldBalanceToFree(w, 40) }, false, 40, 0},
		{"move hold to free with insufficient hold balance", 0, 30, func(w *domain.Wallet) error { return moveHoldBalanceToFree(w, 40) }, true, 0, 30},
		{"increase hold", 0, 10, func(w *domain.Wallet) error { return increaseHoldBalance(w, 40) }, false, 0, 50},
		{"decrease hold", 0, 50, func(w *domain.Wallet) error { return decreaseHoldBalance(w, 40) }, false, 0, 10},
		{"decrease hold with insufficient hold balance", 0, 30, func(w *domain.Wallet) error { return decreaseHoldBalance(w, 40) }, true, 0, 30},
		{"zero amount", 100, 100, func(w *domain.Wallet) error { return moveFreeBalanceToHold(w, 0) }, true, 100, 100},
		{"negative amount", 100, 100, func(w *domain.Wallet) error { return increaseHoldBalance(w, -1) }, true, 100, 100},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := domain.Wallet{ID: "wallet", Status: domain.WALLET_ACTIVE, FreeBalance: test.free, HoldBalance: test.hold}

			err := test.operation(&wallet)
			if (err != nil) != test.wantErr {
				t.Fatalf("unexpected error result: wantErr=%v err=%v", test.wantErr, err)
			}
			if wallet.FreeBalance != test.wantFree || wallet.HoldBalance != test.wantHold {
				t.Fatalf("unexpected balances: free=%d hold=%d", wallet.FreeBalance, wallet.HoldBalance)
			}
		})
	}
}

func TestWalletReservedBalanceOperations(t *testing.T) {
	tests := []struct {
		name         string
		free         int64
		reserved     int64
		operation    func(*domain.Wallet) error
		wantErr      bool
		wantFree     int64
		wantReserved int64
	}{
		{"move free to reserved", 100, 0, func(w *domain.Wallet) error { return moveFreeBalanceToReserved(w, 40) }, false, 60, 40},
		{"move free to reserved with insufficient free balance", 30, 0, func(w *domain.Wallet) error { return moveFreeBalanceToReserved(w, 40) }, true, 30, 0},
		{"move reserved to free", 0, 40, func(w *domain.Wallet) error { return moveReservedBalanceToFree(w, 40) }, false, 40, 0},
		{"move reserved to free with insufficient reserved balance", 0, 30, func(w *domain.Wallet) error { return moveReservedBalanceToFree(w, 40) }, true, 0, 30},
		{"negative amount", 100, 100, func(w *domain.Wallet) error { return moveReservedBalanceToFree(w, -1) }, true, 100, 100},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := domain.Wallet{ID: "wallet", Status: domain.WALLET_ACTIVE, FreeBalance: test.free, ReservedBalance: test.reserved}

			err := test.operation(&wallet)
			if (err != nil) != test.wantErr {
				t.Fatalf("unexpected error result: wantErr=%v err=%v", test.wantErr, err)
			}
			if wallet.FreeBalance != test.wantFree || wallet.ReservedBalance != test.wantReserved || wallet.HoldBalance != 0 {
				t.Fatalf("unexpected balances: free=%d reserved=%d hold=%d", wallet.FreeBalance, wallet.ReservedBalance, wallet.HoldBalance)
			}
		})
	}
}

func TestWalletRoundTripKeepsBalance(t *testing.T) {
	wallet := domain.Wallet{ID: "wallet", Status: domain.WALLET_ACTIVE, FreeBalance: 100}

	if err := moveFreeBalanceToHold(&wallet, 40); err != nil {
		t.Fatal(err)
	}
	if err := moveHoldBalanceToFree(&wallet, 40); err != nil {
		t.Fatal(err)
	}

	if wallet.FreeBalance != 100 || wallet.HoldBalance != 0 {
		t.Fatalf("unexpected balances: free=%d hold=%d", wallet.FreeBalance, wallet.HoldBalance)
	}
}
