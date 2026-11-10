package services

import (
	"testing"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

func TestWalletTypeRequiresAccount(t *testing.T) {
	tests := map[domain.WalletType]bool{
		domain.WALLET_CHECKING:   true,
		domain.WALLET_SAVINGS:    true,
		domain.WALLET_INVESTMENT: false,
		domain.WALLET_CRYPTO:     false,
	}

	for walletType, want := range tests {
		if got := walletTypeRequiresAccount(walletType); got != want {
			t.Fatalf("walletTypeRequiresAccount(%s) = %v, want %v", walletType, got, want)
		}
	}
}
