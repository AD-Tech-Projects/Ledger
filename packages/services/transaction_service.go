package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/models"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/repositories"
)

type TransactionService struct {
	transactionRepository *repositories.GormRepository[models.TransactionModel, domain.Transaction]
	walletRepository      *repositories.GormRepository[models.WalletModel, domain.Wallet]
}

func NewTransactionService(
	transactionRepository *repositories.GormRepository[models.TransactionModel, domain.Transaction],
	walletRepository *repositories.GormRepository[models.WalletModel, domain.Wallet],
) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		walletRepository:      walletRepository,
	}
}

type transactionError struct {
	operation domain.TransactionType
	message   string
}

func (e *transactionError) Error() string {
	return fmt.Errorf("TransactionError-{%s}: %s", e.operation, e.message).Error()
}

var (
	ErrInvalidTransaction  = errors.New("invalid transaction")
	ErrTransactionNotFound = errors.New("transaction not found")
)

func (s *TransactionService) Create(ctx context.Context, amount int64, transactionType domain.TransactionType, sourceWalletID string, destinationWalletID string) (*domain.Transaction, error) {
	transaction := &domain.Transaction{
		Amount: amount,
		Status: domain.TRANSACTION_PENDING,
		Type:   transactionType,
		SourceWallet: &domain.Wallet{
			ID: sourceWalletID,
		},
		DestinationWallet: &domain.Wallet{
			ID: destinationWalletID,
		},
	}

	if transaction.SourceWallet == nil || transaction.DestinationWallet == nil {
		return nil, fmt.Errorf("%w: source and destination wallets are required", ErrInvalidTransaction)
	}
	if _, err := uuid.Parse(transaction.SourceWallet.ID); err != nil {
		return nil, fmt.Errorf("%w: source wallet id must be a valid UUID", ErrInvalidTransaction)
	}
	if _, err := uuid.Parse(transaction.DestinationWallet.ID); err != nil {
		return nil, fmt.Errorf("%w: destination wallet id must be a valid UUID", ErrInvalidTransaction)
	}

	if transaction.ID == "" {
		transaction.ID = uuid.NewString()
	}

	var transactionFailure error
	err := s.transactionRepository.Transaction(ctx, func(tx *gorm.DB) error {
		// TODO - lock queries for source and destination wallets to prevent race conditions
		walletRepository := s.walletRepository.WithDB(tx)
		transactionRepository := s.transactionRepository.WithDB(tx)

		sourceWallet, err := walletRepository.GetByID(ctx, transaction.SourceWallet.ID)
		if err != nil {
			return fmt.Errorf("failed to get source wallet: %w", err)
		}
		destinationWallet, err := walletRepository.GetByID(ctx, transaction.DestinationWallet.ID)
		if err != nil {
			return fmt.Errorf("failed to get destination wallet: %w", err)
		}
		if sourceWallet == nil || destinationWallet == nil {
			return ErrWalletNotFound
		}

		transaction.SourceWallet = sourceWallet
		transaction.DestinationWallet = destinationWallet

		if err := validateTransaction(transaction); err != nil {
			transactionFailure = err
		}
		if err := applyTransaction(transaction); err != nil {
			transactionFailure = err
		}
		if transactionFailure != nil {
			transaction.Status = domain.TRANSACTION_FAILED
			return transactionFailure
		}

		if err := walletRepository.Update(ctx, sourceWallet); err != nil {
			return fmt.Errorf("failed to update source wallet: %w", err)
		}
		if err := walletRepository.Update(ctx, destinationWallet); err != nil {
			return fmt.Errorf("failed to update destination wallet: %w", err)
		}

		transaction.Status = domain.TRANSACTION_COMPLETED
		now := time.Now()
		transaction.CompletedAt = &now
		if err := transactionRepository.Create(ctx, transaction); err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		return nil
	})
	if err != nil {
		if transactionFailure != nil {
			if persistErr := s.persistFailedTransaction(ctx, transaction); persistErr != nil {
				return nil, fmt.Errorf("%w; failed to persist failed transaction: %v", err, persistErr)
			}
		}
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) persistFailedTransaction(
	ctx context.Context,
	transaction *domain.Transaction,
) error {
	existingTransaction, err := s.transactionRepository.GetByID(ctx, transaction.ID)
	if err != nil {
		return fmt.Errorf("failed to find transaction for failure update: %w", err)
	}

	if existingTransaction == nil {
		if err := s.transactionRepository.Create(ctx, transaction); err != nil {
			return fmt.Errorf("failed to create failed transaction: %w", err)
		}
		return nil
	}

	if err := s.transactionRepository.Update(ctx, transaction); err != nil {
		return fmt.Errorf("failed to update failed transaction: %w", err)
	}

	return nil
}

func applyTransaction(transaction *domain.Transaction) error {
	source := transaction.SourceWallet
	destination := transaction.DestinationWallet
	amount := transaction.Amount

	switch transaction.Type {
	case domain.TRANSACTION_DEPOSIT, domain.TRANSACTION_TRANSFER:
		return transferFreeBalance(source, destination, amount, transaction.Type)
	case domain.TRANSACTION_WITHDRAWAL, domain.TRANSACTION_REVERSAL:
		return transferFreeBalance(destination, source, amount, transaction.Type)
	case domain.TRANSACTION_YIELD, domain.TRANSACTION_SALE:
		if err := destination.IncreaseHoldBalance(amount); err != nil {
			return &transactionError{operation: transaction.Type, message: "failed to increase hold balance for destination wallet"}
		}
		return destination.MoveHoldBalanceToFree(amount)
	case domain.TRANSACTION_FEE, domain.TRANSACTION_PURCHASE:
		return transferFreeBalance(source, destination, amount, transaction.Type)
	case domain.TRANSACTION_RESERVE, domain.TRANSACTION_HOLD:
		return source.MoveFreeBalanceToHold(amount)
	case domain.TRANSACTION_RELEASE, domain.TRANSACTION_UNHOLD:
		return source.MoveHoldBalanceToFree(amount)
	default:
		return &transactionValidationError{message: fmt.Sprintf("unsupported transaction type: %s", transaction.Type)}
	}
}

func transferFreeBalance(source *domain.Wallet, destination *domain.Wallet, amount int64, op domain.TransactionType) error {
	if err := source.MoveFreeBalanceToHold(amount); err != nil {
		return &transactionError{operation: op, message: "failed to move free balance to hold for source wallet"}
	}
	if err := destination.IncreaseHoldBalance(amount); err != nil {
		_ = source.MoveHoldBalanceToFree(amount)
		return &transactionError{operation: op, message: "failed to increase hold balance for destination wallet"}
	}
	if err := source.DecreaseHoldBalance(amount); err != nil {
		return &transactionError{operation: op, message: "failed to decrease hold balance for source wallet"}
	}
	if err := destination.MoveHoldBalanceToFree(amount); err != nil {
		return &transactionError{operation: op, message: "failed to move hold balance to free for destination wallet"}
	}
	return nil
}

func (s *TransactionService) GetById(ctx context.Context, id string) (*domain.Transaction, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: transaction id must be a valid UUID", ErrInvalidTransaction)
	}

	transaction, err := s.transactionRepository.GetByID(ctx, id,
		repositories.Join[models.TransactionModel]("SourceWallet"),
		repositories.Join[models.TransactionModel]("SourceWallet.Owner"),
		repositories.Join[models.TransactionModel]("DestinationWallet"),
		repositories.Join[models.TransactionModel]("DestinationWallet.Owner"))
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by id: %w", err)
	}
	if transaction == nil {
		return nil, ErrTransactionNotFound
	}

	return transaction, nil
}

func (s *TransactionService) GetAll(ctx context.Context, filter *domain.Transaction) ([]*domain.Transaction, error) {
	transactions, err := s.transactionRepository.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}
