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

type OwnerService struct {
	ownerRepository *repositories.GormRepository[models.OwnerModel, domain.Owner]
}

func NewOwnerService(
	ownerRepository *repositories.GormRepository[models.OwnerModel, domain.Owner],
) *OwnerService {
	return &OwnerService{
		ownerRepository: ownerRepository,
	}
}

var (
	ErrOwnerNotFound = errors.New("owner not found")
)

func (s *OwnerService) FindByUniqueIdentifier(ctx context.Context, uniqueIdentifier string) (*domain.Owner, error) {
	if uniqueIdentifier == "" {
		return nil, fmt.Errorf("unique identifier is required")
	}

	owner := &domain.Owner{UniqueIdentifier: uniqueIdentifier}
	foundOwner, err := s.ownerRepository.FindOne(ctx, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to find owner by unique identifier: %w", err)
	}

	if foundOwner == nil {
		return nil, fmt.Errorf("%w: owner with unique identifier '%s' not found", ErrOwnerNotFound, uniqueIdentifier)
	}

	return foundOwner, nil
}

func (s *OwnerService) Create(ctx context.Context, owner *domain.Owner) (*domain.Owner, error) {
	if owner == nil {
		return nil, fmt.Errorf("Owner is required")
	}

	if owner.ID == "" {
		owner.ID = uuid.NewString()
	}

	if err := s.ownerRepository.Create(ctx, owner); err != nil {
		return nil, fmt.Errorf("failed to create owner: %w", err)
	}

	return owner, nil
}
