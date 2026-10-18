package account

import (
	holder "github.com/AD-Tech-Projects/Ledger/packages/domain/holder"
)

func NewAccount(
	id string,
	status AccountStatus,
	accountType AccountType,
	currency AccountCurrency,
	balance int64,
	holder *holder.AccountHolder,
) Account {
	return Account{
		Id:              id,
		Status:          status,
		AccountType:     accountType,
		Currency:        currency,
		FreeBalance:     balance,
		ReservedBalance: 0,
		HoldBalance:     0,
		HolderId:        holder.Id,
	}
}

func (a *Account) IncreaseFreeBalance(amount int64) {
	err := validateAmount(amount)
	if err != nil {
		panic(err)
	}

	a.FreeBalance += amount
}

func (a *Account) MoveFreeBalanceToHold(amount int64) {
	err := a.validateFreeBalance(amount)
	if err != nil {
		panic(err)
	}

	a.FreeBalance -= amount
	a.HoldBalance += amount
}

func (a *Account) IncreaseHoldBalance(amount int64) {
	a.HoldBalance += amount
}

func (a *Account) DecreaseHoldBalance(amount int64) {
	err := a.validateHoldBalance(amount)
	if err != nil {
		panic(err)
	}

	a.HoldBalance -= amount
}

func (a *Account) MoveHoldBalanceToFree(amount int64) {
	err := a.validateHoldBalance(amount)
	if err != nil {
		panic(err)
	}

	a.HoldBalance -= amount
	a.FreeBalance += amount
}
