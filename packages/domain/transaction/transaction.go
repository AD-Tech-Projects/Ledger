package transaction

import (
	account "github.com/AD-Tech-Projects/Ledger/packages/domain/account"
)

func MakeTransaction(
	id string,
	amount int64,
	transactionType TransactionType,
	sourceAccount *account.Account,
	destinationAccount *account.Account,
) Transaction {
	if err := validateTransaction(
		amount,
		transactionType,
		sourceAccount,
		destinationAccount); err != nil {
		panic(err)
	}

	switch transactionType {
	case DEPOSIT:
		sourceAccount.MoveFreeBalanceToHold(amount)
	case WITHDRAWAL:
		destinationAccount.MoveFreeBalanceToHold(amount)
	}

	return Transaction{
		id:                 id,
		amount:             amount,
		status:             PENDING,
		transactionType:    transactionType,
		sourceAccount:      sourceAccount,
		destinationAccount: destinationAccount,
	}
}

func (t *Transaction) Commit() {
	transactionStatus := FAILED

	switch t.transactionType {
	case DEPOSIT:
		t.destinationAccount.IncreaseHoldBalance(t.amount)
		t.sourceAccount.DecreaseHoldBalance(t.amount)
		t.destinationAccount.MoveHoldBalanceToFree(t.amount)
	case WITHDRAWAL:
		t.sourceAccount.IncreaseHoldBalance(t.amount)
		t.destinationAccount.DecreaseHoldBalance(t.amount)
		t.sourceAccount.MoveHoldBalanceToFree(t.amount)
	}
	transactionStatus = COMPLETED

	t.status = transactionStatus
}
