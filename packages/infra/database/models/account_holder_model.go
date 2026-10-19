package models

import "github.com/AD-Tech-Projects/Ledger/packages/domain"

type AccountHolderModel struct {
	BaseModel
	UniqueIdentifier string `gorm:"not null;uniqueIndex"`
	HolderType       string `gorm:"not null"`
}

func (AccountHolderModel) TableName() string {
	return "account_holders"
}

func AccountHolderToDomain(m *AccountHolderModel) *domain.AccountHolder {
	return &domain.AccountHolder{
		Id:               m.ID,
		UniqueIdentifier: m.UniqueIdentifier,
		HolderType:       domain.HolderType(m.HolderType),
	}
}

func AccountHolderFromDomain(d *domain.AccountHolder) *AccountHolderModel {
	return &AccountHolderModel{
		BaseModel:        BaseModel{ID: d.Id},
		UniqueIdentifier: d.UniqueIdentifier,
		HolderType:       string(d.HolderType),
	}
}
