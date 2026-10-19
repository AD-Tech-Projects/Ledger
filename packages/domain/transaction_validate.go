package domain

type transactionError struct {
	operation TransactionType
	message   string
}

func validateTransaction(
	amount int64,
	transactionType TransactionType,
	sourceAccount *Account,
	destinationAccount *Account,
) error {
	if amount <= 0 {
		return &transactionError{message: "Transaction amount must be greater than zero", operation: transactionType}
	}

	if sourceAccount.Id == destinationAccount.Id {
		return &transactionError{message: "Source and destination accounts cannot be the same", operation: transactionType}
	}
	if sourceAccount.Currency != destinationAccount.Currency {
		return &transactionError{message: "Source and destination accounts must have the same currency", operation: transactionType}
	}
	if sourceAccount.Status != ACTIVE || destinationAccount.Status != ACTIVE {
		return &transactionError{message: "Both source and destination accounts must be active", operation: transactionType}
	}

	if transactionType == DEPOSIT && sourceAccount.FreeBalance < amount {
		return &transactionError{message: "Insufficient funds in source account for deposit", operation: transactionType}
	}
	if transactionType == WITHDRAWAL && destinationAccount.FreeBalance < amount {
		return &transactionError{message: "Insufficient funds in destination account for withdrawal", operation: transactionType}
	}

	return nil
}
