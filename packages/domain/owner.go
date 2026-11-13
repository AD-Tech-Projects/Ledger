package domain

type Owner struct {
	// NOTE - database unique identifier for the owner
	ID string

	Type OwnerType

	CountryCode      CountryCode // CountryCode is the jurisdiction of the owner's unique document.
	UniqueIdentifier string      // Unique document identifier for the owner, such as CPF/CNPJ in Brazil, or SSN/EIN in the United States
}

type OwnerType string

const (
	OWNER_INDIVIDUAL OwnerType = "INDIVIDUAL" // Natural person, identified by CPF/SSN
	OWNER_COMPANY    OwnerType = "COMPANY"    // Legal entity, identified by CNPJ/EIN
)
