package holder

func NewIndividualAccountHolder(id string, document string) AccountHolder {
	return AccountHolder{
		id:               id,
		uniqueIdentifier: document,
		holderType:       INDIVIDUAL,
	}
}

func NewCompanyAccountHolder(id string, document string) AccountHolder {
	return AccountHolder{
		id:               id,
		uniqueIdentifier: document,
		holderType:       COMPANY,
	}
}

func (h *AccountHolder) GetId() string {
	return h.id
}
