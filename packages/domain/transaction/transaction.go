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
		Id:                 id,
		Amount:             amount,
		Status:             PENDING,
		TransactionType:    transactionType,
		SourceAccount:      sourceAccount,
		DestinationAccount: destinationAccount,
	}
}

func (t *Transaction) Commit() {
	transactionStatus := FAILED

	switch t.TransactionType {
	case DEPOSIT:
		t.DestinationAccount.IncreaseHoldBalance(t.Amount)
		t.SourceAccount.DecreaseHoldBalance(t.Amount)
		t.DestinationAccount.MoveHoldBalanceToFree(t.Amount)
	case WITHDRAWAL:
		t.SourceAccount.IncreaseHoldBalance(t.Amount)
		t.DestinationAccount.DecreaseHoldBalance(t.Amount)
		t.SourceAccount.MoveHoldBalanceToFree(t.Amount)
	}
	transactionStatus = COMPLETED

	t.Status = transactionStatus
}
