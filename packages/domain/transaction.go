package domain

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID        string
	CreatedAt time.Time

	Amount      int64
	Description string

	Type        TransactionType
	Status      TransactionStatus
	CompletedAt *time.Time

	AssetID string
	Asset   *Asset

	SourceAccountID *string
	SourceAccount   *Account

	DestinationAccountID *string
	DestinationAccount   *Account

	// NOTE - identifier for the operation that originated this transaction
	ReferenceID *string

	// NOTE - ID of the transaction that originated this transaction, in case of a reversal
	ReversalOriginID *string
	ReversalOrigin   *Transaction
}

type TransactionType string

const (
	TRANSACTION_DEPOSIT    TransactionType = "DEPOSIT"
	TRANSACTION_WITHDRAWAL TransactionType = "WITHDRAWAL"
	TRANSACTION_TRANSFER   TransactionType = "TRANSFER"
	TRANSACTION_REVERSAL   TransactionType = "REVERSAL"

	TRANSACTION_YIELD    TransactionType = "YIELD"
	TRANSACTION_FEE      TransactionType = "FEE"
	TRANSACTION_PURCHASE TransactionType = "PURCHASE"
	TRANSACTION_SALE     TransactionType = "SALE"

	TRANSACTION_RESERVE TransactionType = "RESERVE"
	TRANSACTION_RELEASE TransactionType = "RELEASE"
	TRANSACTION_HOLD    TransactionType = "HOLD"
	TRANSACTION_UNHOLD  TransactionType = "UNHOLD"
)

type TransactionStatus string

const (
	TRANSACTION_PENDING   TransactionStatus = "PENDING"
	TRANSACTION_COMPLETED TransactionStatus = "COMPLETED"
	TRANSACTION_FAILED    TransactionStatus = "FAILED"
)

func (e *transactionError) Error() string {
	return fmt.Errorf("{%s}: %s", e.operation, e.message).Error()
}

// TODO - mover lógica para service
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
	case TRANSACTION_DEPOSIT:
		sourceAccount.moveFreeBalanceToHold(amount)
	case TRANSACTION_WITHDRAWAL:
		destinationAccount.moveFreeBalanceToHold(amount)
		// TODO - adiconar cases para os outros tipos de transações
	}

	return Transaction{
		ID:                 id,
		Amount:             amount,
		Status:             TRANSACTION_PENDING,
		Type:               transactionType,
		SourceAccount:      sourceAccount,
		DestinationAccount: destinationAccount,
	}
}

func (t *Transaction) Commit() {
	transactionStatus := TRANSACTION_FAILED

	switch t.Type {
	case TRANSACTION_DEPOSIT:
		if err := t.SourceAccount.validateHoldBalance(t.Amount); err != nil {
			panic(err)
		}

		t.DestinationAccount.increaseHoldBalance(t.Amount)
		t.SourceAccount.decreaseHoldBalance(t.Amount)
		t.DestinationAccount.moveHoldBalanceToFree(t.Amount)
	case TRANSACTION_WITHDRAWAL:
		if err := t.DestinationAccount.validateHoldBalance(t.Amount); err != nil {
			panic(err)
		}

		t.SourceAccount.increaseHoldBalance(t.Amount)
		t.DestinationAccount.decreaseHoldBalance(t.Amount)
		t.SourceAccount.moveHoldBalanceToFree(t.Amount)

		// TODO - adiconar cases para os outros tipos de transações
	}
	transactionStatus = TRANSACTION_COMPLETED

	t.Status = transactionStatus
}
