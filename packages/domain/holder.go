package domain

// TODO - mudar nome para Holder
type AccountHolder struct {
	// NOTE - database unique identifier for the holder
	ID string

	Type HolderType
	// NOTE - unique document identifier for the holder, such as CPF/CNPJ in Brazil, or SSN/EIN in the United States
	UniqueIdentifier string
}

type HolderType string

const (
	HOLDER_INDIVIDUAL HolderType = "INDIVIDUAL"
	HOLDER_COMPANY    HolderType = "COMPANY"
)

func newIndividualAccountHolder(id string, document string) AccountHolder {
	return AccountHolder{
		ID:               id,
		UniqueIdentifier: document,
		Type:             HOLDER_INDIVIDUAL,
	}
}

func newCompanyAccountHolder(id string, document string) AccountHolder {
	return AccountHolder{
		ID:               id,
		UniqueIdentifier: document,
		Type:             HOLDER_COMPANY,
	}
}
