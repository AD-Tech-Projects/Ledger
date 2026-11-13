package models

import "github.com/AD-Tech-Projects/Ledger/packages/domain"

type OwnerModel struct {
	BaseModel
	UniqueIdentifier string `gorm:"not null;uniqueIndex"`
	OwnerType        string `gorm:"not null"`
	CountryCode      string `gorm:"not null;size:2;default:''"`
}

func (OwnerModel) TableName() string {
	return "owners"
}

func OwnerToDomain(m *OwnerModel) *domain.Owner {
	if m == nil {
		return nil
	}
	return &domain.Owner{
		ID:               m.ID,
		UniqueIdentifier: m.UniqueIdentifier,
		Type:             domain.OwnerType(m.OwnerType),
		CountryCode:      domain.CountryCode(m.CountryCode),
	}
}

func OwnerFromDomain(d *domain.Owner) *OwnerModel {
	return &OwnerModel{
		ID:               d.ID,
		UniqueIdentifier: d.UniqueIdentifier,
		OwnerType:        string(d.Type),
		CountryCode:      string(d.CountryCode),
	}
}
