package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

var ErrInvalidAssetPrice = errors.New("invalid asset price")

type assetPriceValidationError struct {
	message string
}

func (e *assetPriceValidationError) Error() string {
	return "AssetPriceValidationError:" + e.message
}

func (e *assetPriceValidationError) Unwrap() error {
	return ErrInvalidAssetPrice
}

// validateAssetPrice checks a quote: a positive price of one unit of the asset in the smallest unit of the quote asset,
// both assets and the moment it was observed
// NOTE - the asset may be quoted in itself (the seed quotes USD in USD), so it may be equal to the quote asset
func validateAssetPrice(price *domain.AssetPrice) error {
	if price == nil {
		return &assetPriceValidationError{message: "asset price is required"}
	}
	if price.Price <= 0 {
		return &assetPriceValidationError{message: "price must be greater than zero"}
	}
	if _, err := uuid.Parse(price.AssetID); err != nil {
		return &assetPriceValidationError{message: "asset id must be a valid UUID"}
	}
	if _, err := uuid.Parse(price.QuoteAssetID); err != nil {
		return &assetPriceValidationError{message: "quote asset id must be a valid UUID"}
	}
	if price.RecordedAt.IsZero() {
		return &assetPriceValidationError{message: "recorded at is required"}
	}

	return nil
}
