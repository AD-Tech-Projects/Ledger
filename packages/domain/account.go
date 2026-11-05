package domain

// TODO - mudar nome para Wallet, Account sera a entidade que representa a conta do cliente (e dados bancarios)
type Account struct {
	ID string

	Type    AccountType
	Status  AccountStatus
	Version int64

	// NOTE - store the balance in the smallest unit of the currency (e.g., cents for USD)
	FreeBalance     int64
	ReservedBalance int64
	// NOTE - balance that is on hold, such as for pending transactions or disputes
	HoldBalance int64

	HolderID string
	Holder   *AccountHolder

	AssetID string
	Asset   *Asset
}

type AccountType string

const (
	CHECKING AccountType = "CHECKING"
	SAVINGS  AccountType = "SAVINGS"
)

type AccountStatus string

const (
	ACTIVE   AccountStatus = "ACTIVE"
	INACTIVE AccountStatus = "INACTIVE"
)

func NewAccount(
	id string,
	status AccountStatus,
	accountType AccountType,
	balance int64,
	holder *AccountHolder,
) Account {
	return Account{
		ID:              id,
		Status:          status,
		Type:            accountType,
		FreeBalance:     balance,
		ReservedBalance: 0,
		HoldBalance:     0,
		HolderID:        holder.ID,
		Holder:          nil,
	}
}

func (e *accountError) Error() string {
	return e.message
}

func (a *Account) moveFreeBalanceToHold(amount int64) {
	err := a.validateFreeBalance(amount)
	if err != nil {
		panic(err)
	}

	a.FreeBalance -= amount
	a.HoldBalance += amount
}

func (a *Account) increaseHoldBalance(amount int64) {
	a.HoldBalance += amount
}

func (a *Account) decreaseHoldBalance(amount int64) {
	err := a.validateHoldBalance(amount)
	if err != nil {
		panic(err)
	}

	a.HoldBalance -= amount
}

func (a *Account) moveHoldBalanceToFree(amount int64) {
	err := a.validateHoldBalance(amount)
	if err != nil {
		panic(err)
	}

	a.HoldBalance -= amount
	a.FreeBalance += amount
}
