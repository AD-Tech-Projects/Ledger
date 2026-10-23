package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/models"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/repositories"
)

var ErrInvalidTransaction = errors.New("invalid transaction")

type TransactionService struct {
	transactionRepository *repositories.GormRepository[models.TransactionModel, domain.Transaction]
}

func NewTransactionService(
	transactionRepository *repositories.GormRepository[models.TransactionModel, domain.Transaction],
) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
	}
}

func (s *TransactionService) Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
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

	if err := s.transactionRepository.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}
