package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/models"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/repositories"
)

type AccountHolderService struct {
	accountHolderRepository *repositories.GormRepository[models.AccountHolderModel, domain.AccountHolder]
}

func NewAccountHolderService(
	accountHolderRepository *repositories.GormRepository[models.AccountHolderModel, domain.AccountHolder],
) *AccountHolderService {
	return &AccountHolderService{
		accountHolderRepository: accountHolderRepository,
	}
}

func (s *AccountHolderService) Create(ctx context.Context, holder *domain.AccountHolder) (*domain.AccountHolder, error) {
	if holder == nil {
		return nil, fmt.Errorf("account holder is required")
	}

	if holder.Id == "" {
		holder.Id = uuid.NewString()
	}

	if err := s.accountHolderRepository.Create(ctx, holder); err != nil {
		return nil, fmt.Errorf("failed to create account holder: %w", err)
	}

	return holder, nil
}
