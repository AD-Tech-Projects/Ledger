package models

import "github.com/AD-Tech-Projects/Ledger/packages/domain"

type AssetModel struct {
	BaseModel
	Name          string `gorm:"not null"`
	Type          string `gorm:"not null"`
	Status        string `gorm:"not null"`
	Code          string `gorm:"not null;uniqueIndex"`
	DecimalPlaces uint8  `gorm:"not null"`
}

func (AssetModel) TableName() string {
	return "assets"
}

func AssetToDomain(m *AssetModel) *domain.Asset {
	if m == nil {
		return nil
	}
	return &domain.Asset{
		ID:            m.ID,
		Name:          m.Name,
		Type:          domain.AssetType(m.Type),
		Status:        domain.AssetStatus(m.Status),
		Code:          domain.AssetCode(m.Code),
		DecimalPlaces: m.DecimalPlaces,
	}
}

func AssetFromDomain(d *domain.Asset) *AssetModel {
	if d == nil {
		return nil
	}
	return &AssetModel{
		BaseModel:     BaseModel{ID: d.ID},
		Name:          d.Name,
		Type:          string(d.Type),
		Status:        string(d.Status),
		Code:          string(d.Code),
		DecimalPlaces: d.DecimalPlaces,
	}
}
