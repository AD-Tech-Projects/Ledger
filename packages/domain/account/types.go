package account

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

type AccountError struct {
	message string
}
