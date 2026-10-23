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

type AccountService struct {
	accountRepository *repositories.GormRepository[models.AccountModel, domain.Account]
}

var ErrInvalidAccount = errors.New("Invalid Account")

func NewAccountService(
	accountRepository *repositories.GormRepository[models.AccountModel, domain.Account],
) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
	}
}

func (s *AccountService) Create(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	if account.HolderId == "" {
		return nil, fmt.Errorf("%w: holder id is required", ErrInvalidAccount)
	}

	if _, err := uuid.Parse(account.HolderId); err != nil {
		return nil, fmt.Errorf("%w: holder id must be a valid UUID", ErrInvalidAccount)
	}

	if account.Id == "" {
		account.Id = uuid.New().String()
	}

	if err := s.accountRepository.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (s *AccountService) GetById(ctx context.Context, id string) (*domain.Account, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: account id must be a valid UUID", ErrInvalidAccount)
	}

	account, err := s.accountRepository.GetByID(ctx, id, repositories.Join[models.AccountModel]("Holder"))
	if err != nil {
		return nil, fmt.Errorf("failed to get account by id: %w", err)
	}

	return account, nil
}
