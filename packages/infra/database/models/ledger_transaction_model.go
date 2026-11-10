package models

import (
	"time"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

type LedgerTransactionModel struct {
	BaseModel
	TransactionID    string            `gorm:"type:uuid;not null;uniqueIndex"`
	Transaction      *TransactionModel `gorm:"foreignKey:TransactionID;references:ID"`
	Status           string            `gorm:"not null"`
	Description      string            `gorm:"not null"`
	PostedAt         *time.Time
	ReversalOriginID *string                 `gorm:"type:uuid;index"`
	ReversalOrigin   *LedgerTransactionModel `gorm:"foreignKey:ReversalOriginID;references:ID"`
}

func (LedgerTransactionModel) TableName() string {
	return "ledger_transactions"
}

func LedgerTransactionToDomain(m *LedgerTransactionModel) *domain.LedgerTransaction {
	if m == nil {
		return nil
	}
	var reversalOriginID *string
	if m.ReversalOriginID != nil {
		value := *m.ReversalOriginID
		reversalOriginID = &value
	}
	return &domain.LedgerTransaction{
		ID:               m.ID,
		TransactionID:    m.TransactionID,
		Transaction:      TransactionToDomain(m.Transaction),
		Status:           domain.LedgerTransactionStatus(m.Status),
		Description:      m.Description,
		PostedAt:         m.PostedAt,
		ReversalOriginID: reversalOriginID,
		ReversalOrigin:   LedgerTransactionToDomain(m.ReversalOrigin),
	}
}

func LedgerTransactionFromDomain(d *domain.LedgerTransaction) *LedgerTransactionModel {
	if d == nil {
		return nil
	}
	return &LedgerTransactionModel{
		ID:               d.ID,
		TransactionID:    d.TransactionID,
		Status:           string(d.Status),
		Description:      d.Description,
		PostedAt:         d.PostedAt,
		ReversalOriginID: d.ReversalOriginID,
	}
}
