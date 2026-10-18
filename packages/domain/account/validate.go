package account

func validateAmount(amount int64) error {
	if amount <= 0 {
		return &AccountError{message: "Amount must be greater than zero"}
	}

	return nil
}

func (a *Account) validateHoldBalance(amount int64) error {
	validateAmount(amount)
	if a.HoldBalance < amount {
		return &AccountError{message: "Insufficient hold balance"}
	}

	return nil
}

func (a *Account) validateFreeBalance(amount int64) error {
	validateAmount(amount)
	if a.FreeBalance < amount {
		return &AccountError{message: "Insufficient free balance"}
	}

	return nil
}

func (e *AccountError) Error() string {
	return e.message
}
