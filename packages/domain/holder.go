package domain

type AccountHolder struct {
	// database unique identifier for the holder
	Id string
	// unique document identifier for the holder, such as CPF/CNPJ in Brazil, or SSN/EIN in the United States
	UniqueIdentifier string
	HolderType       HolderType
}

type HolderType string

const (
	INDIVIDUAL HolderType = "INDIVIDUAL"
	COMPANY    HolderType = "COMPANY"
)

func newIndividualAccountHolder(id string, document string) AccountHolder {
	return AccountHolder{
		Id:               id,
		UniqueIdentifier: document,
		HolderType:       INDIVIDUAL,
	}
}

func newCompanyAccountHolder(id string, document string) AccountHolder {
	return AccountHolder{
		Id:               id,
		UniqueIdentifier: document,
		HolderType:       COMPANY,
	}
}
