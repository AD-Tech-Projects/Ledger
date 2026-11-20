package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/models"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/repositories"
	"github.com/AD-Tech-Projects/Ledger/packages/services/helpers"
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
	ErrInvalidOwner  = errors.New("invalid owner")
)

func (s *OwnerService) FindByUniqueIdentifier(
	ctx context.Context,
	countryCode domain.CountryCode,
	ownerType domain.OwnerType,
	document string,
) (*domain.Owner, error) {
	if document == "" {
		return nil, &ownerValidationError{message: "owner document is required"}
	}

	uniqueIdentifier := document
	if !isOwnerUniqueIdentifierBuilt(document) {
		uniqueIdentifier = buildOwnerUniqueIdentifier(countryCode, ownerType, document)
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
		return nil, &ownerValidationError{message: "owner is required"}
	}

	if !isOwnerUniqueIdentifierBuilt(owner.UniqueIdentifier) {
		owner.UniqueIdentifier = buildOwnerUniqueIdentifier(owner.CountryCode, owner.Type, owner.UniqueIdentifier)
	}
	if err := validateOwner(owner); err != nil {
		return nil, err
	}

	if owner.ID == "" {
		owner.ID = uuid.NewString()
	}

	if err := s.ownerRepository.Create(ctx, owner); err != nil {
		return nil, fmt.Errorf("failed to create owner: %w", err)
	}

	return owner, nil
}

// An identifier is considered already built when it has the three parts COUNTRY:TYPE:DOCUMENT
func isOwnerUniqueIdentifierBuilt(identifier string) bool {
	return len(strings.Split(identifier, ":")) == 3
}

// Builds the identifier in the format COUNTRY:TYPE:DOCUMENT in uppercase, keeping only letters and digits of the document
func buildOwnerUniqueIdentifier(countryCode domain.CountryCode, ownerType domain.OwnerType, document string) string {
	return strings.ToUpper(fmt.Sprintf("%s:%s:%s", countryCode, ownerType, helpers.NormalizeCpf(document)))
}
