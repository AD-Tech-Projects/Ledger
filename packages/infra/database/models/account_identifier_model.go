package models

import "github.com/AD-Tech-Projects/Ledger/packages/domain"

type AccountIdentifierModel struct {
	BaseModel
	AccountID      string `gorm:"type:uuid;not null;index"`
	Network        string `gorm:"not null"`
	IdentifierType string `gorm:"column:identifier_type;not null"`
	Value          string `gorm:"not null"`
}

func (AccountIdentifierModel) TableName() string {
	return "account_identifiers"
}

func identifierToDomain(m *AccountIdentifierModel) domain.AccountIdentifier {
	return domain.AccountIdentifier{ID: m.ID, AccountID: m.AccountID, Network: domain.AccountIdentifierNetwork(m.Network),
		Type: domain.AccountIdentifierType(m.IdentifierType), Value: m.Value}
}

func identifierFromDomain(d domain.AccountIdentifier) AccountIdentifierModel {
	return AccountIdentifierModel{ID: d.ID, AccountID: d.AccountID, Network: string(d.Network),
		IdentifierType: string(d.Type), Value: d.Value}
}
