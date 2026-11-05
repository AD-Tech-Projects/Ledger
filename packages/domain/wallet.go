package domain

type Wallet struct {
	ID string

	Type    WalletType
	Status  WalletStatus
	Version int64

	// NOTE - store the balance in the smallest unit of the currency (e.g., cents for USD)
	FreeBalance     int64
	ReservedBalance int64
	// NOTE - balance that is on hold, such as for pending transactions or disputes
	HoldBalance int64

	OwnerID string
	Owner   *Owner

	AssetID string
	Asset   *Asset
}

type WalletType string

const (
	WALLET_CHECKING WalletType = "CHECKING"
	WALLET_SAVINGS  WalletType = "SAVINGS"
)

type WalletStatus string

const (
	WALLET_ACTIVE   WalletStatus = "ACTIVE"
	WALLET_INACTIVE WalletStatus = "INACTIVE"
)

func (e *accountError) Error() string {
	return e.message
}

func NewAccount(
	id string,
	status WalletStatus,
	walletType WalletType,
	balance int64,
	owner *Owner,
) Wallet {
	return Wallet{
		ID:              id,
		Status:          status,
		Type:            walletType,
		FreeBalance:     balance,
		ReservedBalance: 0,
		HoldBalance:     0,
		OwnerID:         owner.ID,
		Owner:           nil,
	}
}

func (w *Wallet) moveFreeBalanceToHold(amount int64) {
	err := w.validateFreeBalance(amount)
	if err != nil {
		panic(err)
	}

	w.FreeBalance -= amount
	w.HoldBalance += amount
}

func (w *Wallet) increaseHoldBalance(amount int64) {
	w.HoldBalance += amount
}

func (w *Wallet) decreaseHoldBalance(amount int64) {
	err := w.validateHoldBalance(amount)
	if err != nil {
		panic(err)
	}

	w.HoldBalance -= amount
}

func (w *Wallet) moveHoldBalanceToFree(amount int64) {
	err := w.validateHoldBalance(amount)
	if err != nil {
		panic(err)
	}

	w.HoldBalance -= amount
	w.FreeBalance += amount
}
