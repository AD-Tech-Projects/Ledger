package transaction

import (
	"fmt"

	account "github.com/AD-Tech-Projects/Ledger/packages/domain/account"
)

func validateTransaction(
	amount int64,
	transactionType TransactionType,
	sourceAccount *account.Account,
	destinationAccount *account.Account,
) error {
	if amount <= 0 {
		return &TransactionError{message: "Transaction amount must be greater than zero", operation: transactionType}
	}

	if sourceAccount.Id == destinationAccount.Id {
		return &TransactionError{message: "Source and destination accounts cannot be the same", operation: transactionType}
	}
	if sourceAccount.Currency != destinationAccount.Currency {
		return &TransactionError{message: "Source and destination accounts must have the same currency", operation: transactionType}
	}
	if sourceAccount.Status != account.ACTIVE || destinationAccount.Status != account.ACTIVE {
		return &TransactionError{message: "Both source and destination accounts must be active", operation: transactionType}
	}

	if transactionType == DEPOSIT && sourceAccount.FreeBalance < amount {
		return &TransactionError{message: "Insufficient funds in source account for deposit", operation: transactionType}
	}
	if transactionType == WITHDRAWAL && destinationAccount.FreeBalance < amount {
		return &TransactionError{message: "Insufficient funds in destination account for withdrawal", operation: transactionType}
	}

	return nil
}

func (e *TransactionError) Error() string {
	return fmt.Errorf("{%s}: %s", e.operation, e.message).Error()
}
