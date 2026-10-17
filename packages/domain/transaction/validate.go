package transaction

import (
	account "github.com/AD-Tech-Projects/Ledger/packages/domain/account"
)

func validateTransaction(
	amount int64,
	transactionType TransactionType,
	sourceAccount *account.Account,
	destinationAccount *account.Account,
) error {
	if amount <= 0 {
		return &TransactionError{message: "Transaction amount must be greater than zero"}
	}

	if sourceAccount.GetId() == destinationAccount.GetId() {
		return &TransactionError{message: "Source and destination accounts cannot be the same"}
	}
	if sourceAccount.GetCurrency() != destinationAccount.GetCurrency() {
		return &TransactionError{message: "Source and destination accounts must have the same currency"}
	}
	if sourceAccount.GetStatus() != account.ACTIVE || destinationAccount.GetStatus() != account.ACTIVE {
		return &TransactionError{message: "Both source and destination accounts must be active"}
	}

	if transactionType == DEPOSIT && sourceAccount.GetFreeBalance() < amount {
		return &TransactionError{message: "Insufficient funds in source account for deposit"}
	}
	if transactionType == WITHDRAWAL && destinationAccount.GetFreeBalance() < amount {
		return &TransactionError{message: "Insufficient funds in destination account for withdrawal"}
	}

	return nil
}

func (e *TransactionError) Error() string {
	return e.message
}
