package domain

type AccountIdentifier struct {
	ID string

	AccountID string
	Account   *Account

	Network AccountIdentifierNetwork // Network identifies the payment rail where this identifier can be used.
	Type    AccountIdentifierType    // Type identifies the semantic kind of value stored in Value.
	Value   string
}

type AccountIdentifierType string

const (
	ACCOUNT_IDENTIFIER_PERSONAL_CODE  AccountIdentifierType = "PERSONAL_CODE"  // CPF in Brazil or SSN in the United States
	ACCOUNT_IDENTIFIER_CORPORATE_CODE AccountIdentifierType = "CORPORATE_CODE" // CNPJ in Brazil or EIN in the United States
	ACCOUNT_IDENTIFIER_PHONE          AccountIdentifierType = "PHONE"
	ACCOUNT_IDENTIFIER_EMAIL          AccountIdentifierType = "EMAIL"
	ACCOUNT_IDENTIFIER_RANDOM         AccountIdentifierType = "RANDOM"
)

type AccountIdentifierNetwork string

const (
	ACCOUNT_IDENTIFIER_NETWORK_BR_PIX  AccountIdentifierNetwork = "BR_PIX"
	ACCOUNT_IDENTIFIER_NETWORK_EU_SEPA AccountIdentifierNetwork = "EU_SEPA"
	ACCOUNT_IDENTIFIER_NETWORK_US_WIRE AccountIdentifierNetwork = "US_WIRE"

	// NOTE - unknown network, used for identifiers that do not have a specific network association, such as random identifiers
	ACCOUNT_IDENTIFIER_NETWORK_UNKNOWN AccountIdentifierNetwork = "UNKNOWN"
)
