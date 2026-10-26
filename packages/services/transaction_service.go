package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/models"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/repositories"
)

var ErrInvalidTransaction = errors.New("invalid transaction")
var ErrTransactionNotFound = errors.New("transaction not found")

type TransactionService struct {
	transactionRepository *repositories.GormRepository[models.TransactionModel, domain.Transaction]
	accountRepository     *repositories.GormRepository[models.AccountModel, domain.Account]
}

func NewTransactionService(
	transactionRepository *repositories.GormRepository[models.TransactionModel, domain.Transaction],
	accountRepository *repositories.GormRepository[models.AccountModel, domain.Account],
) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		accountRepository:     accountRepository,
	}
}

func (s *TransactionService) Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	if transaction == nil {
		return nil, fmt.Errorf("%w: transaction is required", ErrInvalidTransaction)
	}

	if transaction.SourceAccount == nil || transaction.DestinationAccount == nil {
		return nil, fmt.Errorf("%w: source and destination accounts are required", ErrInvalidTransaction)
	}

	if _, err := uuid.Parse(transaction.SourceAccount.Id); err != nil {
		return nil, fmt.Errorf("%w: source account id must be a valid UUID", ErrInvalidTransaction)
	}
	if _, err := uuid.Parse(transaction.DestinationAccount.Id); err != nil {
		return nil, fmt.Errorf("%w: destination account id must be a valid UUID", ErrInvalidTransaction)
	}

	if transaction.Id == "" {
		transaction.Id = uuid.NewString()
	}

	err := s.transactionRepository.Transaction(ctx, func(tx *gorm.DB) error {
		accountRepository := s.accountRepository.WithDB(tx)
		transactionRepository := s.transactionRepository.WithDB(tx)

		sourceAccount, err := accountRepository.GetByID(ctx, transaction.SourceAccount.Id)
		if err != nil {
			return fmt.Errorf("failed to get source account: %w", err)
		}
		if sourceAccount == nil {
			return ErrAccountNotFound
		}
		destinationAccount, err := accountRepository.GetByID(ctx, transaction.DestinationAccount.Id)
		if err != nil {
			return fmt.Errorf("failed to get destination account: %w", err)
		}
		if destinationAccount == nil {
			return ErrAccountNotFound
		}

		if err := makeTransactionSafely(transaction, sourceAccount, destinationAccount); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidTransaction, err)
		}

		if err := accountRepository.Update(ctx, sourceAccount); err != nil {
			return fmt.Errorf("failed to update source account: %w", err)
		}
		if err := accountRepository.Update(ctx, destinationAccount); err != nil {
			return fmt.Errorf("failed to update destination account: %w", err)
		}
		if err := transactionRepository.Create(ctx, transaction); err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		if err := commitTransactionSafely(transaction); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidTransaction, err)
		}

		if err := accountRepository.Update(ctx, sourceAccount); err != nil {
			return fmt.Errorf("failed to update source account after commit: %w", err)
		}
		if err := accountRepository.Update(ctx, destinationAccount); err != nil {
			return fmt.Errorf("failed to update destination account after commit: %w", err)
		}
		if err := transactionRepository.Update(ctx, transaction); err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func makeTransactionSafely(
	transaction *domain.Transaction,
	sourceAccount *domain.Account,
	destinationAccount *domain.Account,
) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("%v", recovered)
		}
	}()

	*transaction = domain.MakeTransaction(
		transaction.Id,
		transaction.Amount,
		transaction.TransactionType,
		sourceAccount,
		destinationAccount,
	)
	return nil
}

func commitTransactionSafely(transaction *domain.Transaction) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("%v", recovered)
		}
	}()

	transaction.Commit()
	return nil
}

func (s *TransactionService) GetById(ctx context.Context, id string) (*domain.Transaction, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: transaction id must be a valid UUID", ErrInvalidTransaction)
	}

	transaction, err := s.transactionRepository.GetByID(ctx, id,
		repositories.Join[models.TransactionModel]("SourceAccount"),
		repositories.Join[models.TransactionModel]("SourceAccount.Holder"),
		repositories.Join[models.TransactionModel]("DestinationAccount"),
		repositories.Join[models.TransactionModel]("DestinationAccount.Holder"))
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by id: %w", err)
	}
	if transaction == nil {
		return nil, ErrTransactionNotFound
	}

	return transaction, nil
}

func (s *TransactionService) GetAll(ctx context.Context, filter *domain.Transaction) ([]*domain.Transaction, error) {
	transactions, err := s.transactionRepository.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}
