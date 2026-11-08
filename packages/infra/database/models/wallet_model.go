package models

import "github.com/AD-Tech-Projects/Ledger/packages/domain"

type WalletModel struct {
	BaseModel
	Type            string        `gorm:"not null"`
	Status          string        `gorm:"not null"`
	Version         int64         `gorm:"not null"`
	FreeBalance     int64         `gorm:"not null"`
	ReservedBalance int64         `gorm:"not null"`
	HoldBalance     int64         `gorm:"not null"`
	OwnerID         string        `gorm:"type:uuid;not null;index"`
	Owner           *OwnerModel   `gorm:"foreignKey:OwnerID;references:ID"`
	AssetID         string        `gorm:"type:uuid;not null;index"`
	Asset           *AssetModel   `gorm:"foreignKey:AssetID;references:ID"`
	AccountID       *string       `gorm:"type:uuid;index"`
	Account         *AccountModel `gorm:"foreignKey:WalletID;references:ID"`
}

func (WalletModel) TableName() string {
	return "wallets"
}

func WalletToDomain(m *WalletModel) *domain.Wallet {
	if m == nil {
		return nil
	}
	accountID := m.AccountID
	account := AccountToDomain(m.Account)
	if account != nil {
		accountID = &account.ID
	}
	return &domain.Wallet{ID: m.ID, Type: domain.WalletType(m.Type), Status: domain.WalletStatus(m.Status),
		Version: m.Version, FreeBalance: m.FreeBalance, ReservedBalance: m.ReservedBalance,
		HoldBalance: m.HoldBalance, OwnerID: m.OwnerID, Owner: OwnerToDomain(m.Owner),
		AssetID: m.AssetID, Asset: AssetToDomain(m.Asset), AccountID: accountID, Account: account}
}

func WalletFromDomain(d *domain.Wallet) *WalletModel {
	if d == nil {
		return nil
	}
	model := &WalletModel{
		BaseModel: BaseModel{ID: d.ID}, Type: string(d.Type), Status: string(d.Status),
		Version: d.Version, FreeBalance: d.FreeBalance, ReservedBalance: d.ReservedBalance,
		HoldBalance: d.HoldBalance, OwnerID: d.OwnerID, AssetID: d.AssetID, AccountID: d.AccountID,
	}
	return model
}
