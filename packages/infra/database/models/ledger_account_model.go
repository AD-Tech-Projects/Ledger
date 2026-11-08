package models

import (
	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

type LedgerAccountModel struct {
	BaseModel
	Type                   string      `gorm:"not null"`
	Status                 string      `gorm:"not null"`
	NormalBalanceDirection string      `gorm:"not null"`
	Name                   string      `gorm:"not null"`
	Code                   string      `gorm:"not null;uniqueIndex"`
	AssetID                string      `gorm:"type:uuid;not null;index"`
	Asset                  *AssetModel `gorm:"foreignKey:AssetID;references:ID"`
}

func (LedgerAccountModel) TableName() string {
	return "ledger_accounts"
}

func LedgerAccountToDomain(m *LedgerAccountModel) *domain.LedgerAccount {
	if m == nil {
		return nil
	}
	return &domain.LedgerAccount{
		ID:                     m.ID,
		Type:                   domain.LedgerAccountType(m.Type),
		Status:                 domain.LedgerAccountStatus(m.Status),
		NormalBalanceDirection: domain.LedgerBalanceDirection(m.NormalBalanceDirection),
		Name:                   m.Name,
		Code:                   m.Code,
		AssetID:                m.AssetID,
		Asset:                  AssetToDomain(m.Asset),
	}
}

func LedgerAccountFromDomain(d *domain.LedgerAccount) *LedgerAccountModel {
	if d == nil {
		return nil
	}
	return &LedgerAccountModel{
		BaseModel:              BaseModel{ID: d.ID},
		Type:                   string(d.Type),
		Status:                 string(d.Status),
		NormalBalanceDirection: string(d.NormalBalanceDirection),
		Name:                   d.Name,
		Code:                   d.Code,
		AssetID:                d.AssetID,
	}
}
