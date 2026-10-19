package models

import (
	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

type AccountModel struct {
	BaseModel
	Status          string              `gorm:"not null"`
	AccountType     string              `gorm:"not null"`
	Currency        string              `gorm:"not null"`
	FreeBalance     int64               `gorm:"not null"`
	ReservedBalance int64               `gorm:"not null"`
	HoldBalance     int64               `gorm:"not null"`
	HolderID        string              `gorm:"type:uuid;not null;index"`
	Holder          *AccountHolderModel `gorm:"foreignKey:HolderID;references:ID"`
}

func (AccountModel) TableName() string {
	return "accounts"
}

func AccountToDomain(m *AccountModel) *domain.Account {
	return &domain.Account{
		Id:              m.ID,
		Status:          domain.AccountStatus(m.Status),
		AccountType:     domain.AccountType(m.AccountType),
		Currency:        domain.AccountCurrency(m.Currency),
		FreeBalance:     m.FreeBalance,
		ReservedBalance: m.ReservedBalance,
		HoldBalance:     m.HoldBalance,
		HolderId:        m.HolderID,
		Holder:          accountHolderToDomain(m.Holder),
	}
}

func accountHolderToDomain(m *AccountHolderModel) *domain.AccountHolder {
	if m == nil {
		return nil
	}

	return AccountHolderToDomain(m)
}

func AccountFromDomain(d *domain.Account) *AccountModel {
	return &AccountModel{
		BaseModel:       BaseModel{ID: d.Id},
		Status:          string(d.Status),
		AccountType:     string(d.AccountType),
		Currency:        string(d.Currency),
		FreeBalance:     d.FreeBalance,
		ReservedBalance: d.ReservedBalance,
		HoldBalance:     d.HoldBalance,
		HolderID:        d.HolderId,
	}
}
