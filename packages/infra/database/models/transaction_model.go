package models

import "github.com/AD-Tech-Projects/Ledger/packages/domain"

type TransactionModel struct {
	BaseModel
	Amount               int64         `gorm:"not null"`
	Status               string        `gorm:"not null"`
	TransactionType      string        `gorm:"not null"`
	SourceAccountID      string        `gorm:"type:uuid;not null;index"`
	SourceAccount        *AccountModel `gorm:"foreignKey:SourceAccountID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DestinationAccountID string        `gorm:"type:uuid;not null;index"`
	DestinationAccount   *AccountModel `gorm:"foreignKey:DestinationAccountID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
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
		SourceAccount:      transactionAccountToDomain(m.SourceAccount, m.SourceAccountID),
		DestinationAccount: transactionAccountToDomain(m.DestinationAccount, m.DestinationAccountID),
	}
}

func transactionAccountToDomain(m *AccountModel, id string) *domain.Account {
	if m == nil {
		return &domain.Account{Id: id}
	}

	return AccountToDomain(m)
}

func TransactionFromDomain(d *domain.Transaction) *TransactionModel {
	var sourceAccountID, destinationAccountID string
	if d.SourceAccount != nil {
		sourceAccountID = d.SourceAccount.Id
	}
	if d.DestinationAccount != nil {
		destinationAccountID = d.DestinationAccount.Id
	}

	return &TransactionModel{
		BaseModel:            BaseModel{ID: d.Id},
		Amount:               d.Amount,
		Status:               string(d.Status),
		TransactionType:      string(d.TransactionType),
		SourceAccountID:      sourceAccountID,
		DestinationAccountID: destinationAccountID,
	}
}
