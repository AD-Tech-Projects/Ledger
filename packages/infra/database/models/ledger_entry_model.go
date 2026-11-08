package models

import "github.com/AD-Tech-Projects/Ledger/packages/domain"

type LedgerEntryModel struct {
	BaseModel
	LedgerTransactionID string                  `gorm:"type:uuid;not null;index"`
	LedgerTransaction   *LedgerTransactionModel `gorm:"foreignKey:LedgerTransactionID;references:ID"`
	LedgerAccountID     string                  `gorm:"type:uuid;not null;index"`
	LedgerAccount       *LedgerAccountModel     `gorm:"foreignKey:LedgerAccountID;references:ID"`
	AssetID             string                  `gorm:"type:uuid;not null;index"`
	Asset               *AssetModel             `gorm:"foreignKey:AssetID;references:ID"`
	Amount              int64                   `gorm:"not null"`
	Direction           string                  `gorm:"not null"`
	Sequence            int64                   `gorm:"not null"`
}

func (LedgerEntryModel) TableName() string {
	return "ledger_entries"
}

func LedgerEntryToDomain(m *LedgerEntryModel) *domain.LedgerEntry {
	if m == nil {
		return nil
	}
	return &domain.LedgerEntry{
		ID:                  m.ID,
		LedgerTransactionID: m.LedgerTransactionID,
		LedgerTransaction:   LedgerTransactionToDomain(m.LedgerTransaction),
		LedgerAccountID:     m.LedgerAccountID,
		LedgerAccount:       LedgerAccountToDomain(m.LedgerAccount),
		AssetID:             m.AssetID,
		Asset:               AssetToDomain(m.Asset),
		Amount:              m.Amount,
		Direction:           domain.LedgerBalanceDirection(m.Direction),
		Sequence:            m.Sequence,
	}
}

func LedgerEntryFromDomain(d *domain.LedgerEntry) *LedgerEntryModel {
	if d == nil {
		return nil
	}
	return &LedgerEntryModel{
		BaseModel:           BaseModel{ID: d.ID},
		LedgerTransactionID: d.LedgerTransactionID,
		LedgerAccountID:     d.LedgerAccountID,
		AssetID:             d.AssetID,
		Amount:              d.Amount,
		Direction:           string(d.Direction),
		Sequence:            d.Sequence,
	}
}
