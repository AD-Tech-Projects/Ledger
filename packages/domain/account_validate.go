package domain

type accountError struct {
	message string
}

func validateAmount(amount int64) error {
	if amount <= 0 {
		return &accountError{message: "Amount must be greater than zero"}
	}

	return nil
}

func (a *Account) validateHoldBalance(amount int64) error {
	validateAmount(amount)
	if a.HoldBalance < amount {
		return &accountError{message: "Insufficient hold balance"}
	}

	return nil
}

func (a *Account) validateFreeBalance(amount int64) error {
	validateAmount(amount)
	if a.FreeBalance < amount {
		return &accountError{message: "Insufficient free balance"}
	}

	return nil
}
