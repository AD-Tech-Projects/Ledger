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
		id:              id,
		status:          status,
		accountType:     accountType,
		currency:        currency,
		freeBalance:     balance,
		reservedBalance: 0,
		holdBalance:     0,
		holderId:        holder.GetId(),
	}
}

func (a *Account) GetId() string {
	return a.id
}

func (a *Account) GetStatus() AccountStatus {
	return a.status
}

func (a *Account) GetCurrency() AccountCurrency {
	return a.currency
}

func (a *Account) GetFreeBalance() int64 {
	return a.freeBalance
}

func (a *Account) IncreaseFreeBalance(amount int64) error {
	validateAmount(amount)

	a.freeBalance += amount

	return nil
}

func (a *Account) MoveFreeBalanceToHold(amount int64) error {
	a.validateFreeBalance(amount)

	a.freeBalance -= amount
	a.holdBalance += amount

	return nil
}

func (a *Account) DecreaseHoldBalance(amount int64) error {
	a.validateHoldBalance(amount)

	a.holdBalance -= amount

	return nil
}

func (a *Account) MoveHoldBalanceToFree(amount int64) error {
	a.validateHoldBalance(amount)

	a.holdBalance -= amount
	a.freeBalance += amount

	return nil
}

func (e *AccountError) Error() string {
	return e.message
}
