package holder

func NewIndividualAccountHolder(id string, document string) AccountHolder {
	return AccountHolder{
		Id:               id,
		UniqueIdentifier: document,
		HolderType:       INDIVIDUAL,
	}
}

func NewCompanyAccountHolder(id string, document string) AccountHolder {
	return AccountHolder{
		Id:               id,
		UniqueIdentifier: document,
		HolderType:       COMPANY,
	}
}
