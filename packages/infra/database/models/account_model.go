package models

import (
	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

type AccountModel struct {
	Id              string `gorm:"primaryKey"`
	Status          string `gorm:"not null"`
	AccountType     string `gorm:"not null"`
	Currency        string `gorm:"not null"`
	FreeBalance     int64  `gorm:"not null"`
	ReservedBalance int64  `gorm:"not null"`
	HoldBalance     int64  `gorm:"not null"`
	HolderId        string `gorm:"not null"`
}

func (AccountModel) TableName() string {
	return "accounts"
}

func AccountToDomain(m *AccountModel) *domain.Account {
	return &domain.Account{
		Id:              m.Id,
		Status:          domain.AccountStatus(m.Status),
		AccountType:     domain.AccountType(m.AccountType),
		Currency:        domain.AccountCurrency(m.Currency),
		FreeBalance:     m.FreeBalance,
		ReservedBalance: m.ReservedBalance,
		HoldBalance:     m.HoldBalance,
		HolderId:        m.HolderId,
	}
}

func AccountFromDomain(d *domain.Account) *AccountModel {
	return &AccountModel{
		Id:              d.Id,
		Status:          string(d.Status),
		AccountType:     string(d.AccountType),
		Currency:        string(d.Currency),
		FreeBalance:     d.FreeBalance,
		ReservedBalance: d.ReservedBalance,
		HoldBalance:     d.HoldBalance,
		HolderId:        d.HolderId,
	}
}
