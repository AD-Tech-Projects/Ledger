package transaction

type Transaction struct {
	id                   string
	amount               int64
	status               TransactionStatus
	transactionType      TransactionType
	sourceAccountId      string
	destinationAccountId string
}

type TransactionType string

const (
	DEPOSIT    TransactionType = "DEPOSIT"
	WITHDRAWAL TransactionType = "WITHDRAWAL"
)

type TransactionStatus string

const (
	PENDING   TransactionStatus = "PENDING"
	COMPLETED TransactionStatus = "COMPLETED"
	FAILED    TransactionStatus = "FAILED"
)

type TransactionError struct {
	message string
}
