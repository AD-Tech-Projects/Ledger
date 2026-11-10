package services

import (
	"errors"
	"fmt"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

var ErrTransactionFailed = errors.New("transaction failed")

type transactionError struct {
	operation domain.TransactionType
	message   string
	cause     error
}

func (e *transactionError) Error() string {
	message := fmt.Sprintf("TransactionError-{%s}: %s", e.operation, e.message)
	if e.cause != nil {
		message += " - " + e.cause.Error()
	}
	return message
}

// NOTE - lets errors.Is match both ErrTransactionFailed (mapped to HTTP 422 by the handler) and the original cause,
// such as the insufficient balance error returned by the wallet operations
func (e *transactionError) Unwrap() []error {
	if e.cause == nil {
		return []error{ErrTransactionFailed}
	}
	return []error{ErrTransactionFailed, e.cause}
}

type walletRequirement int

const (
	walletForbidden walletRequirement = iota // The wallet must not be informed (external side of the operation)
	walletOptional                           // The wallet may be informed
	walletRequired                           // The wallet must be informed
)

type transactionWalletRule struct {
	source      walletRequirement
	destination walletRequirement
}

// NOTE - which wallets take part in each transaction type (REVERSAL derives its wallets from the original transaction)
var transactionWalletRules = map[domain.TransactionType]transactionWalletRule{
	domain.TRANSACTION_DEPOSIT:    {source: walletForbidden, destination: walletRequired},
	domain.TRANSACTION_WITHDRAWAL: {source: walletRequired, destination: walletForbidden},
	domain.TRANSACTION_TRANSFER:   {source: walletRequired, destination: walletRequired},

	domain.TRANSACTION_YIELD:    {source: walletOptional, destination: walletRequired}, // source: wallet that generated the yield (not debited)
	domain.TRANSACTION_SALE:     {source: walletOptional, destination: walletRequired}, // source: buyer wallet (debited when informed)
	domain.TRANSACTION_FEE:      {source: walletRequired, destination: walletOptional}, // destination: fee revenue wallet (credited when informed)
	domain.TRANSACTION_PURCHASE: {source: walletRequired, destination: walletOptional}, // destination: seller wallet (credited when informed)

	domain.TRANSACTION_RESERVE: {source: walletRequired, destination: walletForbidden},
	domain.TRANSACTION_RELEASE: {source: walletRequired, destination: walletForbidden},
	domain.TRANSACTION_HOLD:    {source: walletRequired, destination: walletForbidden},
	domain.TRANSACTION_UNHOLD:  {source: walletRequired, destination: walletForbidden},
}

// NOTE - balance effect that undoes each transaction type when applied to the reversal wallets (see getReversalWalletIDs)
var reversalEffectTypes = map[domain.TransactionType]domain.TransactionType{
	domain.TRANSACTION_DEPOSIT:    domain.TRANSACTION_WITHDRAWAL,
	domain.TRANSACTION_WITHDRAWAL: domain.TRANSACTION_DEPOSIT,
	domain.TRANSACTION_TRANSFER:   domain.TRANSACTION_TRANSFER,
	domain.TRANSACTION_YIELD:      domain.TRANSACTION_WITHDRAWAL, // the yield came from outside the platform, so it only leaves the credited wallet
	domain.TRANSACTION_FEE:        domain.TRANSACTION_SALE,       // refunds the payer, taking the funds back from the revenue wallet when there is one
	domain.TRANSACTION_SALE:       domain.TRANSACTION_PURCHASE,
	domain.TRANSACTION_PURCHASE:   domain.TRANSACTION_SALE,
	domain.TRANSACTION_RESERVE:    domain.TRANSACTION_RELEASE,
	domain.TRANSACTION_RELEASE:    domain.TRANSACTION_RESERVE,
	domain.TRANSACTION_HOLD:       domain.TRANSACTION_UNHOLD,
	domain.TRANSACTION_UNHOLD:     domain.TRANSACTION_HOLD,
}

// isBalanceStateTransaction returns true if the transaction type only moves funds between the balance states (free,
// reserved and hold) of a single wallet, without changing its total balance. Movement types (including FEE) change the
// total balance, moving funds between wallets or into/out of the platform, so their reversal must reverse the direction.
func isBalanceStateTransaction(transactionType domain.TransactionType) bool {
	switch transactionType {
	case domain.TRANSACTION_RESERVE, domain.TRANSACTION_RELEASE, domain.TRANSACTION_HOLD, domain.TRANSACTION_UNHOLD:
		return true
	default:
		return false
	}
}

// getReversalWalletIDs returns the wallets of a reversal based on the original
// transaction: balance state changes keep the original wallet, while movements
// swap source and destination, since the funds go back the opposite way.
func getReversalWalletIDs(origin *domain.Transaction) (sourceWalletID *string, destinationWalletID *string) {
	originSourceID := walletIDOf(origin.SourceWalletID, origin.SourceWallet)
	originDestinationID := walletIDOf(origin.DestinationWalletID, origin.DestinationWallet)

	if isBalanceStateTransaction(origin.Type) {
		return originSourceID, nil
	}

	sourceWalletID, destinationWalletID = originDestinationID, originSourceID
	// NOTE - a wallet that the reversal effect does not touch is left out, such as the wallet that generated a yield,
	// which was never debited and therefore is not credited back
	if rule, ok := transactionWalletRules[reversalEffectTypes[origin.Type]]; ok {
		if rule.source == walletForbidden {
			sourceWalletID = nil
		}
		if rule.destination == walletForbidden {
			destinationWalletID = nil
		}
	}
	return sourceWalletID, destinationWalletID
}

func applyTransaction(transaction *domain.Transaction) error {
	if transaction == nil {
		return &transactionValidationError{message: "transaction is required"}
	}

	effectType := transaction.Type
	if transaction.Type == domain.TRANSACTION_REVERSAL {
		if transaction.ReversalOrigin == nil {
			return &transactionValidationError{message: "reversal origin transaction is required"}
		}
		reversalType, ok := reversalEffectTypes[transaction.ReversalOrigin.Type]
		if !ok {
			return &transactionValidationError{message: fmt.Sprintf("transaction type %q cannot be reversed", transaction.ReversalOrigin.Type)}
		}
		effectType = reversalType
	}

	if err := applyBalanceEffect(effectType, transaction.SourceWallet, transaction.DestinationWallet, transaction.Amount); err != nil {
		return &transactionError{operation: transaction.Type, message: "failed to apply balance changes", cause: err}
	}
	return nil
}

func applyBalanceEffect(effectType domain.TransactionType, source *domain.Wallet, destination *domain.Wallet, amount int64) error {
	switch effectType {
	case domain.TRANSACTION_DEPOSIT, domain.TRANSACTION_YIELD:
		return creditWallet(destination, amount)
	case domain.TRANSACTION_WITHDRAWAL:
		return debitWallet(source, amount)
	case domain.TRANSACTION_TRANSFER:
		return transferFreeBalance(source, destination, amount)
	case domain.TRANSACTION_SALE:
		if source != nil {
			return transferFreeBalance(source, destination, amount)
		}
		return creditWallet(destination, amount)
	case domain.TRANSACTION_FEE, domain.TRANSACTION_PURCHASE:
		if destination != nil {
			return transferFreeBalance(source, destination, amount)
		}
		return debitWallet(source, amount)
	case domain.TRANSACTION_RESERVE:
		if err := requireWallet(source, "source"); err != nil {
			return err
		}
		return moveFreeBalanceToReserved(source, amount)
	case domain.TRANSACTION_RELEASE:
		if err := requireWallet(source, "source"); err != nil {
			return err
		}
		return moveReservedBalanceToFree(source, amount)
	case domain.TRANSACTION_HOLD:
		if err := requireWallet(source, "source"); err != nil {
			return err
		}
		return moveFreeBalanceToHold(source, amount)
	case domain.TRANSACTION_UNHOLD:
		if err := requireWallet(source, "source"); err != nil {
			return err
		}
		return moveHoldBalanceToFree(source, amount)
	default:
		return fmt.Errorf("unsupported transaction type: %s", effectType)
	}
}

func creditWallet(wallet *domain.Wallet, amount int64) error {
	if err := requireWallet(wallet, "destination"); err != nil {
		return err
	}
	return increaseFreeBalance(wallet, amount)
}

func debitWallet(wallet *domain.Wallet, amount int64) error {
	if err := requireWallet(wallet, "source"); err != nil {
		return err
	}
	return decreaseFreeBalance(wallet, amount)
}

// transferFreeBalance moves funds in two steps through the hold balance
// (source free -> source hold -> destination hold -> destination free), so a failure never leaves a partial transfer.
func transferFreeBalance(source *domain.Wallet, destination *domain.Wallet, amount int64) error {
	if err := requireWallet(source, "source"); err != nil {
		return err
	}
	if err := requireWallet(destination, "destination"); err != nil {
		return err
	}

	if err := moveFreeBalanceToHold(source, amount); err != nil {
		return err
	}
	if err := increaseHoldBalance(destination, amount); err != nil {
		_ = moveHoldBalanceToFree(source, amount)
		return err
	}
	if err := decreaseHoldBalance(source, amount); err != nil {
		_ = decreaseHoldBalance(destination, amount)
		_ = moveHoldBalanceToFree(source, amount)
		return err
	}
	return moveHoldBalanceToFree(destination, amount)
}

func requireWallet(wallet *domain.Wallet, side string) error {
	if wallet == nil {
		return fmt.Errorf("%s wallet is required", side)
	}
	return nil
}

func walletIDOf(id *string, wallet *domain.Wallet) *string {
	if id != nil {
		return id
	}
	if wallet != nil && wallet.ID != "" {
		walletID := wallet.ID
		return &walletID
	}
	return nil
}
