package domain

import (
	"fmt"
)

type Transaction struct {
	Id                 string
	Amount             int64
	Status             TransactionStatus
	TransactionType    TransactionType
	SourceAccount      *Account
	DestinationAccount *Account
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

func (e *transactionError) Error() string {
	return fmt.Errorf("{%s}: %s", e.operation, e.message).Error()
}

func MakeTransaction(
	id string,
	amount int64,
	transactionType TransactionType,
	sourceAccount *Account,
	destinationAccount *Account,
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
		sourceAccount.moveFreeBalanceToHold(amount)
	case WITHDRAWAL:
		destinationAccount.moveFreeBalanceToHold(amount)
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
		t.DestinationAccount.increaseHoldBalance(t.Amount)
		t.SourceAccount.decreaseHoldBalance(t.Amount)
		t.DestinationAccount.moveHoldBalanceToFree(t.Amount)
	case WITHDRAWAL:
		t.SourceAccount.increaseHoldBalance(t.Amount)
		t.DestinationAccount.decreaseHoldBalance(t.Amount)
		t.SourceAccount.moveHoldBalanceToFree(t.Amount)
	}
	transactionStatus = COMPLETED

	t.Status = transactionStatus
}
