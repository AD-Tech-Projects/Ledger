package services

import (
	"context"
	"fmt"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/models"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/repositories"
)

type AssetService struct {
	assetRepository *repositories.GormRepository[models.AssetModel, domain.Asset]
}

func NewAssetService(
	assetRepository *repositories.GormRepository[models.AssetModel, domain.Asset],
) *AssetService {
	return &AssetService{
		assetRepository: assetRepository,
	}
}

func (s *AssetService) FindByCode(ctx context.Context, code string) (*domain.Asset, error) {
	if code == "" {
		return nil, fmt.Errorf("asset code is required")
	}

	assetCode := domain.AssetCode(code)
	if assetCode == "" {
		return nil, fmt.Errorf("invalid asset code")
	}

	asset := &domain.Asset{Code: assetCode}
	foundAsset, err := s.assetRepository.FindOne(ctx, asset)
	if err != nil {
		return nil, fmt.Errorf("failed to find asset by code: %w", err)
	}

	if foundAsset == nil {
		return nil, fmt.Errorf("asset with code %s not found", code)
	}

	return foundAsset, nil
}
