package account

func validateAmount(amount int64) error {
	if amount <= 0 {
		return &AccountError{message: "Amount must be greater than zero"}
	}

	return nil
}

func (a *Account) validateHoldBalance(amount int64) error {
	validateAmount(amount)
	if a.holdBalance < amount {
		return &AccountError{message: "Insufficient hold balance"}
	}

	return nil
}

func (a *Account) validateFreeBalance(amount int64) error {
	validateAmount(amount)
	if a.freeBalance < amount {
		return &AccountError{message: "Insufficient free balance"}
	}

	return nil
}
