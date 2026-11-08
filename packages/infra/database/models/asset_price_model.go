package models

import (
	"time"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

type AssetPriceModel struct {
	BaseModel
	Price        int64       `gorm:"not null"`
	RecordedAt   time.Time   `gorm:"not null;index"`
	AssetID      string      `gorm:"type:uuid;not null;index"`
	Asset        *AssetModel `gorm:"foreignKey:AssetID;references:ID"`
	QuoteAssetID string      `gorm:"type:uuid;not null;index"`
	QuoteAsset   *AssetModel `gorm:"foreignKey:QuoteAssetID;references:ID"`
}

func (AssetPriceModel) TableName() string {
	return "asset_prices"
}

func AssetPriceToDomain(m *AssetPriceModel) *domain.AssetPrice {
	if m == nil {
		return nil
	}
	return &domain.AssetPrice{
		ID:           m.ID,
		Price:        m.Price,
		RecordedAt:   m.RecordedAt,
		AssetID:      m.AssetID,
		Asset:        AssetToDomain(m.Asset),
		QuoteAssetID: m.QuoteAssetID,
		QuoteAsset:   AssetToDomain(m.QuoteAsset),
	}
}

func AssetPriceFromDomain(d *domain.AssetPrice) *AssetPriceModel {
	if d == nil {
		return nil
	}
	return &AssetPriceModel{
		BaseModel:    BaseModel{ID: d.ID},
		Price:        d.Price,
		RecordedAt:   d.RecordedAt,
		AssetID:      d.AssetID,
		QuoteAssetID: d.QuoteAssetID,
	}
}
