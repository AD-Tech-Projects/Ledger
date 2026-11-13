package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"

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

func (s *OwnerService) FindByUniqueIdentifier(
	ctx context.Context,
	countryCode domain.CountryCode,
	ownerType domain.OwnerType,
	document string,
) (*domain.Owner, error) {
	if document == "" {
		return nil, fmt.Errorf("unique identifier is required")
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
		return nil, fmt.Errorf("Owner is required")
	}

	if owner.ID == "" {
		owner.ID = uuid.NewString()
	}

	if !isOwnerUniqueIdentifierBuilt(owner.UniqueIdentifier) {
		owner.UniqueIdentifier = buildOwnerUniqueIdentifier(owner.CountryCode, owner.Type, owner.UniqueIdentifier)
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
	normalizedDocument := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, document)

	return strings.ToUpper(fmt.Sprintf("%s:%s:%s", countryCode, ownerType, normalizedDocument))
}
