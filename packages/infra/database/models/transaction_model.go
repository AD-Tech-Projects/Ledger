package models

import (
	"time"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

type TransactionModel struct {
	BaseModel
	Amount              int64             `gorm:"not null;check:amount > 0"`
	Status              string            `gorm:"not null"`
	Type                string            `gorm:"column:transaction_type;not null"`
	AssetID             string            `gorm:"type:uuid;not null;index"`
	Asset               *AssetModel       `gorm:"foreignKey:AssetID;references:ID"`
	SourceWalletID      *string           `gorm:"type:uuid;index"`
	SourceWallet        *WalletModel      `gorm:"foreignKey:SourceWalletID;references:ID"`
	DestinationWalletID *string           `gorm:"type:uuid;index"`
	DestinationWallet   *WalletModel      `gorm:"foreignKey:DestinationWalletID;references:ID"`
	ReferenceID         *string           `gorm:"index"`
	ReversalOriginID    *string           `gorm:"type:uuid;index"`
	ReversalOrigin      *TransactionModel `gorm:"foreignKey:ReversalOriginID;references:ID"`
	Description         string            `gorm:"not null"`
	CompletedAt         *time.Time
}

func (TransactionModel) TableName() string {
	return "transactions"
}

func TransactionToDomain(m *TransactionModel) *domain.Transaction {
	if m == nil {
		return nil
	}
	return &domain.Transaction{
		ID:                  m.ID,
		Amount:              m.Amount,
		Status:              domain.TransactionStatus(m.Status),
		Type:                domain.TransactionType(m.Type),
		AssetID:             m.AssetID,
		Asset:               AssetToDomain(m.Asset),
		SourceWalletID:      m.SourceWalletID,
		SourceWallet:        transactionWalletToDomain(m.SourceWallet, m.SourceWalletID),
		DestinationWalletID: m.DestinationWalletID,
		DestinationWallet:   transactionWalletToDomain(m.DestinationWallet, m.DestinationWalletID),
		ReferenceID:         m.ReferenceID,
		ReversalOriginID:    m.ReversalOriginID,
		ReversalOrigin:      TransactionToDomain(m.ReversalOrigin),
		Description:         m.Description,
		CreatedAt:           m.CreatedAt,
		CompletedAt:         m.CompletedAt,
	}
}

func transactionWalletToDomain(m *WalletModel, id *string) *domain.Wallet {
	if m == nil {
		if id == nil {
			return nil
		}
		return &domain.Wallet{ID: *id}
	}
	return WalletToDomain(m)
}

func TransactionFromDomain(d *domain.Transaction) *TransactionModel {
	var sourceAccountID, destinationAccountID *string
	if d.SourceWallet != nil {
		sourceAccountID = &d.SourceWallet.ID
	} else {
		sourceAccountID = d.SourceWalletID
	}
	if d.DestinationWallet != nil {
		destinationAccountID = &d.DestinationWallet.ID
	} else {
		destinationAccountID = d.DestinationWalletID
	}
	return &TransactionModel{
		ID:                  d.ID,
		CreatedAt:           d.CreatedAt,
		Amount:              d.Amount,
		Status:              string(d.Status),
		Type:                string(d.Type),
		AssetID:             d.AssetID,
		SourceWalletID:      sourceAccountID,
		DestinationWalletID: destinationAccountID,
		ReferenceID:         d.ReferenceID,
		ReversalOriginID:    d.ReversalOriginID,
		Description:         d.Description,
		CompletedAt:         d.CompletedAt,
	}
}
