package models

import (
	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

type AccountModel struct {
	BaseModel
	WalletID        string       `gorm:"type:uuid;not null;uniqueIndex"`
	Wallet          *WalletModel `gorm:"foreignKey:WalletID;references:ID"`
	Kind            string       `gorm:"not null"`
	CountryCode     string       `gorm:"size:2"`
	Currency        string       `gorm:"not null"`
	InstitutionName string
	BranchCode      string
	AccountNumber   string
	RoutingCode     string
	IBAN            string
	BIC             string
	CryptoNetwork   string
	CryptoAddress   string
	CryptoMemo      string
	Identifiers     []AccountIdentifierModel `gorm:"foreignKey:AccountID"`
}

func (AccountModel) TableName() string {
	return "external_accounts"
}

func AccountToDomain(m *AccountModel) *domain.Account {
	if m == nil {
		return nil
	}
	account := &domain.Account{ID: m.ID, WalletID: m.WalletID, Kind: domain.AccountKind(m.Kind),
		CountryCode: domain.AccountCountryCode(m.CountryCode), Currency: m.Currency, InstitutionName: m.InstitutionName,
		BranchCode: m.BranchCode, AccountNumber: m.AccountNumber, RoutingCode: m.RoutingCode,
		IBAN: m.IBAN, BIC: m.BIC, CryptoNetwork: m.CryptoNetwork,
		CryptoAddress: m.CryptoAddress, CryptoMemo: m.CryptoMemo}
	for _, identifier := range m.Identifiers {
		account.Identifiers = append(account.Identifiers, identifierToDomain(&identifier))
	}
	return account
}

func AccountFromDomain(d *domain.Account) *AccountModel {
	if d == nil {
		return nil
	}
	model := &AccountModel{
		BaseModel: BaseModel{ID: d.ID}, WalletID: d.WalletID, Kind: string(d.Kind),
		CountryCode: string(d.CountryCode), Currency: d.Currency, InstitutionName: d.InstitutionName,
		BranchCode: d.BranchCode, AccountNumber: d.AccountNumber, RoutingCode: d.RoutingCode,
		IBAN: d.IBAN, BIC: d.BIC, CryptoNetwork: d.CryptoNetwork, CryptoAddress: d.CryptoAddress,
		CryptoMemo: d.CryptoMemo,
	}
	for _, identifier := range d.Identifiers {
		model.Identifiers = append(model.Identifiers, identifierFromDomain(identifier))
	}
	return model
}
