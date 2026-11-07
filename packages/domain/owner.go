package domain

type Owner struct {
	// NOTE - database unique identifier for the owner
	ID string

	Type OwnerType

	// NOTE - unique document identifier for the owner, such as CPF/CNPJ in Brazil, or SSN/EIN in the United States
	UniqueIdentifier string
}

type OwnerType string

const (
	OWNER_INDIVIDUAL OwnerType = "INDIVIDUAL"
	OWNER_COMPANY    OwnerType = "COMPANY"
)
