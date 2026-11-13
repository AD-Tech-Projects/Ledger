package domain

// NOTE - the account entity represents a financial account,
// which can be either a bank account or a crypto account.
// It belongs to an Owner, can be linked to multiple Wallets and can have multiple identifiers,
// such as personal code, corporate code, phone number, email, or random identifier.
type Account struct {
	ID string

	OwnerID string
	Owner   *Owner

	Kind        AccountKind // Kind defines whether the bank fields or the crypto fields are filled.
	CountryCode CountryCode // CountryCode is the jurisdiction of the account, which defines the available payment rails.

	Currency        string // Currency is the ISO currency code of the bank account (e.g. BRL, USD).
	InstitutionName string // InstitutionName is the name of the bank or financial institution.
	BranchCode      string // BranchCode is the branch identifier used by the institution, when applicable.
	AccountNumber   string // AccountNumber is the account number at the institution.
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
	ACCOUNT_BANK   AccountKind = "BANK"   // Account at a bank or payment institution
	ACCOUNT_CRYPTO AccountKind = "CRYPTO" // Address on a blockchain network
)
