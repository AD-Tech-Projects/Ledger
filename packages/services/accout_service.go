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

var (
	ErrInvalidAccount  = errors.New("invalid account")
	ErrAccountNotFound = errors.New("account not found")
)

type AccountPatch struct {
	Status      *domain.AccountStatus
	AccountType *domain.AccountType
	Currency    *domain.AccountCurrency
	HolderId    *string
}

func NewAccountService(
	accountRepository *repositories.GormRepository[models.AccountModel, domain.Account],
) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
	}
}

func (s *AccountService) Create(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	if account == nil {
		return nil, fmt.Errorf("%w: account is required", ErrInvalidAccount)
	}

	if account.HolderId == "" {
		return nil, fmt.Errorf("%w: holder id is required", ErrInvalidAccount)
	}

	if _, err := uuid.Parse(account.HolderId); err != nil {
		return nil, fmt.Errorf("%w: holder id must be a valid UUID", ErrInvalidAccount)
	}

	if account.Id == "" {
		account.Id = uuid.New().String()
	}

	account.FreeBalance = 0
	account.ReservedBalance = 0
	account.HoldBalance = 0

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
	if account == nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func (s *AccountService) GetAll(ctx context.Context) ([]*domain.Account, error) {
	accounts, err := s.accountRepository.GetAll(ctx, repositories.Join[models.AccountModel]("Holder"))
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	return accounts, nil
}

func (s *AccountService) Update(ctx context.Context, id string, patch AccountPatch) (*domain.Account, error) {
	account, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if patch.Status != nil {
		account.Status = *patch.Status
	}
	if patch.AccountType != nil {
		account.AccountType = *patch.AccountType
	}
	if patch.Currency != nil {
		account.Currency = *patch.Currency
	}
	if patch.HolderId != nil {
		if _, err := uuid.Parse(*patch.HolderId); err != nil {
			return nil, fmt.Errorf("%w: holder id must be a valid UUID", ErrInvalidAccount)
		}
		account.HolderId = *patch.HolderId
	}

	if err := s.accountRepository.Update(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return s.GetById(ctx, id)
}

func (s *AccountService) Deactivate(ctx context.Context, id string) error {
	account, err := s.GetById(ctx, id)
	if err != nil {
		return err
	}

	account.Status = domain.INACTIVE
	if err := s.accountRepository.Update(ctx, account); err != nil {
		return fmt.Errorf("failed to deactivate account: %w", err)
	}

	return nil
}
