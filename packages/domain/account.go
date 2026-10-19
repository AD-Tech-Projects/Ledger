package domain

type Account struct {
	Id          string
	Status      AccountStatus
	AccountType AccountType
	Currency    AccountCurrency
	// store the balance in the smallest unit of the currency (e.g., cents for USD)
	FreeBalance     int64
	ReservedBalance int64
	HoldBalance     int64
	// database unique identifier for the holder of the account
	HolderId string
}

type AccountCurrency string

const (
	USD AccountCurrency = "USD"
	EUR AccountCurrency = "EUR"
	BRL AccountCurrency = "BRL"
)

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
	currency AccountCurrency,
	balance int64,
	holder *AccountHolder,
) Account {
	return Account{
		Id:              id,
		Status:          status,
		AccountType:     accountType,
		Currency:        currency,
		FreeBalance:     balance,
		ReservedBalance: 0,
		HoldBalance:     0,
		HolderId:        holder.Id,
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
