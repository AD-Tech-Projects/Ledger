package account

type Account struct {
	id          string
	status      AccountStatus
	accountType AccountType
	currency    AccountCurrency
	// store the balance in the smallest unit of the currency (e.g., cents for USD)
	freeBalance     int64
	reservedBalance int64
	holdBalance     int64
	// database unique identifier for the holder of the account
	holderId string
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
