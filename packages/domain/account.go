package domain

// NOTE - the account entity represents a financial account,
// which can be either a bank account or a crypto account.
// It is associated with a Wallet and can have multiple identifiers,
// such as personal code, corporate code, phone number, email, or random identifier.
type Account struct {
	ID string

	WalletID string
	Wallet   *Wallet

	Kind        AccountKind
	CountryCode AccountCountryCode

	Currency        string
	InstitutionName string
	BranchCode      string // BranchCode is the branch identifier used by the institution, when applicable.
	AccountNumber   string
	RoutingCode     string // RoutingCode is the country-specific bank clearing or routing identifier.
	IBAN            string // International bank account number, used for international bank transfers
	BIC             string // International bank identifier code, used for international bank transfers (also known as SWIFT code)

	CryptoNetwork string // CryptoNetwork identifies the blockchain network for a crypto account.
	CryptoAddress string // CryptoAddress is the public blockchain address used for the account.
	CryptoMemo    string // CryptoMemo is an optional memo, tag, or destination tag required by a network.

	Identifiers []AccountIdentifier // Identifiers contains payment-network identifiers associated with the account.
}

type AccountKind string

const (
	ACCOUNT_BANK   AccountKind = "BANK"
	ACCOUNT_CRYPTO AccountKind = "CRYPTO"
)

type AccountCountryCode string

const (
	COUNTRY_US AccountCountryCode = "US"
	COUNTRY_BR AccountCountryCode = "BR"
	COUNTRY_EU AccountCountryCode = "EU"

	// NOTE - unknown country code, used for accounts that do not have a specific country association, such as crypto accounts
	COUNTRY_UNKNOWN AccountCountryCode = "UNKNOWN"
)
