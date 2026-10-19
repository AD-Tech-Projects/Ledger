package models

import "github.com/AD-Tech-Projects/Ledger/packages/domain"

type TransactionModel struct {
	BaseModel
	Amount               int64  `gorm:"not null"`
	Status               string `gorm:"not null"`
	TransactionType      string `gorm:"not null"`
	SourceAccountID      string `gorm:"type:uuid;not null;index"`
	DestinationAccountID string `gorm:"type:uuid;not null;index"`
}

func (TransactionModel) TableName() string {
	return "transactions"
}

func TransactionToDomain(m *TransactionModel) *domain.Transaction {
	return &domain.Transaction{
		Id:                 m.ID,
		Amount:             m.Amount,
		Status:             domain.TransactionStatus(m.Status),
		TransactionType:    domain.TransactionType(m.TransactionType),
		SourceAccount:      &domain.Account{Id: m.SourceAccountID},
		DestinationAccount: &domain.Account{Id: m.DestinationAccountID},
	}
}

func TransactionFromDomain(d *domain.Transaction) *TransactionModel {
	return &TransactionModel{
		BaseModel:            BaseModel{ID: d.Id},
		Amount:               d.Amount,
		Status:               string(d.Status),
		TransactionType:      string(d.TransactionType),
		SourceAccountID:      d.SourceAccount.Id,
		DestinationAccountID: d.DestinationAccount.Id,
	}
}
