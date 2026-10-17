package transaction

import (
	account "github.com/AD-Tech-Projects/Ledger/packages/domain/account"
)

func MakeTransaction(
	id string,
	amount int64,
	status TransactionStatus,
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

	sourceAccount.MoveFreeBalanceToHold(amount)
	destinationAccount.IncreaseFreeBalance(amount)
	sourceAccount.DecreaseHoldBalance(amount)

	return Transaction{
		id:                   id,
		amount:               amount,
		status:               status,
		transactionType:      transactionType,
		sourceAccountId:      sourceAccount.GetId(),
		destinationAccountId: destinationAccount.GetId(),
	}
}
