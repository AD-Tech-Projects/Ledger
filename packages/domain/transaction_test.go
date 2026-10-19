package domain

import (
	"testing"
)

func TestSuccessfulTransaction(t *testing.T) {
	companyHolder := newCompanyAccountHolder("1", "12345678000195")
	personHolder := newIndividualAccountHolder("2", "123456789")

	companyAccount := createAccount("1", 10000, companyHolder)
	personAccount := createAccount("2", 0, personHolder)

	ledgerTransaction := makeTransaction(
		"1",
		5000,
		DEPOSIT,
		&companyAccount,
		&personAccount,
	)
	ledgerTransaction.commit()

	assertTransactionCompletion(t, &ledgerTransaction)
	assertCompanyAccountFreeBalance(t, &companyAccount)
	assertPersonAccountFreeBalance(t, &personAccount)
	assertHeldBalance(t, &companyAccount, &personAccount)
}

func assertTransactionCompletion(t *testing.T, transaction *Transaction) {
	const expected = COMPLETED
	result := transaction.Status

	if result != expected {
		t.Fatalf("expected transaction status %q, got %q", COMPLETED, transaction.Status)
	}
}

func assertCompanyAccountFreeBalance(t *testing.T, companyAccount *Account) {
	const expected = 5000
	result := companyAccount.FreeBalance

	if result != expected {
		t.Fatalf("expected company free balance 5000, got %d", companyAccount.FreeBalance)
	}
}

func assertPersonAccountFreeBalance(t *testing.T, personAccount *Account) {
	const expected = 5000
	result := personAccount.FreeBalance

	if result != expected {
		t.Fatalf("expected person free balance 5000, got %d", personAccount.FreeBalance)
	}
}

func assertHeldBalance(t *testing.T, companyAccount *Account, personAccount *Account) {
	if companyAccount.HoldBalance != 0 || personAccount.HoldBalance != 0 {
		t.Fatalf(
			"expected no held balance after commit, got company=%d person=%d",
			companyAccount.HoldBalance,
			personAccount.HoldBalance,
		)
	}
}

func TestTransactionRollback(t *testing.T) {
	companyHolder := newCompanyAccountHolder("1", "12345678000195")
	personHolder := newIndividualAccountHolder("2", "123456789")

	companyAccount := createAccount("1", 10000, companyHolder)
	personAccount := createAccount("2", 0, personHolder)

	assertTransactionCreationPanics(t, func() {
		makeTransaction(
			"1",
			15000,
			DEPOSIT,
			&companyAccount,
			&personAccount,
		)
	})

	if companyAccount.FreeBalance != 10000 || companyAccount.HoldBalance != 0 {
		t.Fatalf(
			"expected company account to remain unchanged, got free balance=%d, hold balance=%d",
			companyAccount.FreeBalance,
			companyAccount.HoldBalance,
		)
	}
	if personAccount.FreeBalance != 0 || personAccount.HoldBalance != 0 {
		t.Fatalf(
			"expected person account to remain unchanged, got free balance=%d, hold balance=%d",
			personAccount.FreeBalance,
			personAccount.HoldBalance,
		)
	}
}

func assertTransactionCreationPanics(t *testing.T, createTransaction func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Fatal("expected transaction creation to panic")
		}
	}()

	createTransaction()
}

func createAccount(
	id string,
	balance int64,
	holder AccountHolder) Account {
	return NewAccount(
		id,
		ACTIVE,
		CHECKING,
		USD,
		balance,
		&holder,
	)
}
