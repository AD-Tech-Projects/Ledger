package holder

type AccountHolder struct {
	// database unique identifier for the holder
	id string
	// unique document identifier for the holder, such as CPF/CNPJ in Brazil, or SSN/EIN in the United States
	uniqueIdentifier string
	holderType       HolderType
}

type HolderType string

const (
	INDIVIDUAL HolderType = "INDIVIDUAL"
	COMPANY    HolderType = "COMPANY"
)
