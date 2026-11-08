package domain

import "testing"

func TestWalletBalanceOperations(t *testing.T) {
	tests := []struct {
		name      string
		free      int64
		hold      int64
		operation func(*Wallet) error
		wantErr   bool
		wantFree  int64
		wantHold  int64
	}{
		{"move free to hold", 100, 0, func(w *Wallet) error { return w.MoveFreeBalanceToHold(40) }, false, 60, 40},
		{"move free to hold with insufficient free balance", 30, 0, func(w *Wallet) error { return w.MoveFreeBalanceToHold(40) }, true, 30, 0},
		{"move hold to free", 0, 40, func(w *Wallet) error { return w.MoveHoldBalanceToFree(40) }, false, 40, 0},
		{"move hold to free with insufficient hold balance", 0, 30, func(w *Wallet) error { return w.MoveHoldBalanceToFree(40) }, true, 0, 30},
		{"increase hold", 0, 10, func(w *Wallet) error { return w.IncreaseHoldBalance(40) }, false, 0, 50},
		{"decrease hold", 0, 50, func(w *Wallet) error { return w.DecreaseHoldBalance(40) }, false, 0, 10},
		{"decrease hold with insufficient hold balance", 0, 30, func(w *Wallet) error { return w.DecreaseHoldBalance(40) }, true, 0, 30},
		{"zero amount", 100, 100, func(w *Wallet) error { return w.MoveFreeBalanceToHold(0) }, true, 100, 100},
		{"negative amount", 100, 100, func(w *Wallet) error { return w.IncreaseHoldBalance(-1) }, true, 100, 100},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := Wallet{ID: "wallet", Status: WALLET_ACTIVE, FreeBalance: test.free, HoldBalance: test.hold}

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

func TestWalletRoundTripKeepsBalance(t *testing.T) {
	wallet := Wallet{ID: "wallet", Status: WALLET_ACTIVE, FreeBalance: 100}

	if err := wallet.MoveFreeBalanceToHold(40); err != nil {
		t.Fatal(err)
	}
	if err := wallet.MoveHoldBalanceToFree(40); err != nil {
		t.Fatal(err)
	}

	if wallet.FreeBalance != 100 || wallet.HoldBalance != 0 {
		t.Fatalf("unexpected balances: free=%d hold=%d", wallet.FreeBalance, wallet.HoldBalance)
	}
}

func TestWalletTypeRequiresAccount(t *testing.T) {
	tests := map[WalletType]bool{
		WALLET_CHECKING:   true,
		WALLET_SAVINGS:    true,
		WALLET_INVESTMENT: false,
		WALLET_CRYPTO:     false,
	}

	for walletType, want := range tests {
		if got := walletType.RequiresAccount(); got != want {
			t.Fatalf("%s.RequiresAccount() = %v, want %v", walletType, got, want)
		}
	}
}
