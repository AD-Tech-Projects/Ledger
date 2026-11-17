package repositories

import (
	"testing"
)

type testOwnerModel struct {
	ID string
}

type testWalletModel struct {
	ID      string
	Owner   *testOwnerModel
	OwnerID string
}

type testTransactionModel struct {
	SourceWallet      *testWalletModel
	DestinationWallet *testWalletModel
	Wallets           []testWalletModel
	ValueWallet       testWalletModel
}

// # Main Flux
func TestMustAssociationCases(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"pointer field of the target type", MustAssociation[testWalletModel, testOwnerModel]("Owner").path, "Owner"},
		{"one of two fields of the same type", MustAssociation[testTransactionModel, testWalletModel]("DestinationWallet").path, "DestinationWallet"},
		{
			"nested association",
			Nested(
				MustAssociation[testTransactionModel, testWalletModel]("SourceWallet"),
				MustAssociation[testWalletModel, testOwnerModel]("Owner"),
			).path,
			"SourceWallet.Owner",
		},
		{
			"nested association with an undeclared parent",
			Nested(Association[testTransactionModel, testWalletModel]{}, MustAssociation[testWalletModel, testOwnerModel]("Owner")).path,
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("association path = %q, want %q", test.got, test.want)
			}
		})
	}
}

// # Failures
func TestMustAssociationRejectsUndeclaredField(t *testing.T) {
	tests := []struct {
		name    string
		declare func()
	}{
		{"missing field", func() { MustAssociation[testWalletModel, testOwnerModel]("Asset") }},
		{"empty field name", func() { MustAssociation[testWalletModel, testOwnerModel]("") }},
		{"field of another type", func() { MustAssociation[testWalletModel, testWalletModel]("Owner") }},
		{"non association field", func() { MustAssociation[testWalletModel, string]("OwnerID") }},
		{"slice field", func() { MustAssociation[testTransactionModel, testWalletModel]("Wallets") }},
		{"value field", func() { MustAssociation[testTransactionModel, testWalletModel]("ValueWallet") }},
		{"model that is not a struct", func() { MustAssociation[string, testOwnerModel]("Owner") }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertPanics(t, test.declare)
		})
	}
}

// helpers
func assertPanics(t *testing.T, declare func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatalf("expected MustAssociation to panic")
		}
	}()
	declare()
}
