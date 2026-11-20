package services

import (
	"errors"
	"testing"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

// # Main Flux
func TestValidateOwnerCases(t *testing.T) {
	tests := []struct {
		name        string
		countryCode domain.CountryCode
		ownerType   domain.OwnerType
		document    string
	}{
		{"brazilian individual with a CPF", domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "529.982.247-25"},
		{"brazilian individual with a CPF ending in zero check digit", domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "123.456.789-09"},
		{"brazilian individual with an alphanumeric CPF", domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "12A.BC3.456-02"},
		{"brazilian individual with a lowercase alphanumeric CPF", domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "12a.bc3.456-02"},
		{"brazilian company with a numeric CNPJ", domain.COUNTRY_BR, domain.OWNER_COMPANY, "11.222.333/0001-81"},
		{"brazilian company with an alphanumeric CNPJ", domain.COUNTRY_BR, domain.OWNER_COMPANY, "12.ABC.345/01DE-35"},
		{"brazilian company with a lowercase alphanumeric CNPJ", domain.COUNTRY_BR, domain.OWNER_COMPANY, "12.abc.345/01de-35"},
		{"american individual with a SSN", domain.COUNTRY_US, domain.OWNER_INDIVIDUAL, "123-45-6789"},
		{"american company with an EIN", domain.COUNTRY_US, domain.OWNER_COMPANY, "12-3456789"},
		{"european individual with any document", domain.COUNTRY_EU, domain.OWNER_INDIVIDUAL, "X1234567L"},
		{"european company with any document", domain.COUNTRY_EU, domain.OWNER_COMPANY, "DE123456789"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			owner := newTestOwner(test.countryCode, test.ownerType, test.document)
			identifier := owner.UniqueIdentifier

			if err := validateOwner(owner); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertOwnerIdentifierUnchanged(t, owner, identifier)
		})
	}
}

// # Failures
func TestValidateOwnerRejectsInvalidOwners(t *testing.T) {
	tests := []struct {
		name  string
		owner *domain.Owner
	}{
		{"nil owner", nil},
		{"unknown country", newTestOwner(domain.COUNTRY_UNKNOWN, domain.OWNER_INDIVIDUAL, "52998224725")},
		{"unsupported country", newTestOwner("XX", domain.OWNER_INDIVIDUAL, "52998224725")},
		{"empty owner type", newTestOwner(domain.COUNTRY_BR, "", "52998224725")},
		{"invalid owner type", newTestOwner(domain.COUNTRY_BR, "FOO", "52998224725")},
		{"lowercase owner type", newTestOwner(domain.COUNTRY_BR, "individual", "52998224725")},
		{"empty document", newTestOwner(domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "")},
		{"document with only punctuation", newTestOwner(domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "...-")},
		{"identifier without the document part", &domain.Owner{CountryCode: domain.COUNTRY_BR, Type: domain.OWNER_INDIVIDUAL, UniqueIdentifier: "BR:INDIVIDUAL"}},
		{"identifier of another country", &domain.Owner{CountryCode: domain.COUNTRY_BR, Type: domain.OWNER_INDIVIDUAL, UniqueIdentifier: "US:INDIVIDUAL:52998224725"}},
		{"identifier of another owner type", &domain.Owner{CountryCode: domain.COUNTRY_BR, Type: domain.OWNER_INDIVIDUAL, UniqueIdentifier: "BR:COMPANY:52998224725"}},
		{"CPF with a wrong first check digit", newTestOwner(domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "52998224735")},
		{"CPF with a wrong second check digit", newTestOwner(domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "52998224726")},
		{"CPF with repeated digits", newTestOwner(domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "111.111.111-11")},
		{"CPF with zeros", newTestOwner(domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "000.000.000-00")},
		{"CPF too short", newTestOwner(domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "5299822472")},
		{"CPF too long", newTestOwner(domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "529982247250")},
		{"CPF with letters in the check digits", newTestOwner(domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "5299822472A")},
		{"alphanumeric CPF with a wrong check digit", newTestOwner(domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "12ABC345603")},
		{"CPF with repeated letters", newTestOwner(domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "AAAAAAAAAAA")},
		{"CPF with a non ASCII letter", newTestOwner(domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "12ÇBC345602")},
		{"CNPJ given for an individual", newTestOwner(domain.COUNTRY_BR, domain.OWNER_INDIVIDUAL, "11222333000181")},
		{"CPF given for a company", newTestOwner(domain.COUNTRY_BR, domain.OWNER_COMPANY, "52998224725")},
		{"CNPJ with a wrong first check digit", newTestOwner(domain.COUNTRY_BR, domain.OWNER_COMPANY, "11222333000191")},
		{"CNPJ with a wrong second check digit", newTestOwner(domain.COUNTRY_BR, domain.OWNER_COMPANY, "11222333000182")},
		{"CNPJ with repeated digits", newTestOwner(domain.COUNTRY_BR, domain.OWNER_COMPANY, "00000000000000")},
		{"alphanumeric CNPJ with a wrong check digit", newTestOwner(domain.COUNTRY_BR, domain.OWNER_COMPANY, "12ABC34501DE36")},
		{"CNPJ with letters in the check digits", newTestOwner(domain.COUNTRY_BR, domain.OWNER_COMPANY, "12ABC34501DE3A")},
		{"CNPJ with a non ASCII letter", newTestOwner(domain.COUNTRY_BR, domain.OWNER_COMPANY, "12ÇBC34501DE35")},
		{"SSN with area 000", newTestOwner(domain.COUNTRY_US, domain.OWNER_INDIVIDUAL, "000-45-6789")},
		{"SSN with area 666", newTestOwner(domain.COUNTRY_US, domain.OWNER_INDIVIDUAL, "666-45-6789")},
		{"SSN with area 9XX (ITIN)", newTestOwner(domain.COUNTRY_US, domain.OWNER_INDIVIDUAL, "912-45-6789")},
		{"SSN with group 00", newTestOwner(domain.COUNTRY_US, domain.OWNER_INDIVIDUAL, "123-00-6789")},
		{"SSN with serial 0000", newTestOwner(domain.COUNTRY_US, domain.OWNER_INDIVIDUAL, "123-45-0000")},
		{"SSN too short", newTestOwner(domain.COUNTRY_US, domain.OWNER_INDIVIDUAL, "12345678")},
		{"EIN with letters", newTestOwner(domain.COUNTRY_US, domain.OWNER_COMPANY, "12-345678A")},
		{"EIN with zeros", newTestOwner(domain.COUNTRY_US, domain.OWNER_COMPANY, "00-0000000")},
		{"EIN too long", newTestOwner(domain.COUNTRY_US, domain.OWNER_COMPANY, "12-34567890")},
		{"european owner without document", newTestOwner(domain.COUNTRY_EU, domain.OWNER_COMPANY, "")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			identifier := ""
			if test.owner != nil {
				identifier = test.owner.UniqueIdentifier
			}

			err := validateOwner(test.owner)
			var validationErr *ownerValidationError
			if !errors.As(err, &validationErr) || !errors.Is(err, ErrInvalidOwner) {
				t.Fatalf("expected owner validation error, got %v", err)
			}
			if test.owner != nil {
				assertOwnerIdentifierUnchanged(t, test.owner, identifier)
			}
		})
	}
}

// helpers
func newTestOwner(countryCode domain.CountryCode, ownerType domain.OwnerType, document string) *domain.Owner {
	return &domain.Owner{
		Type:             ownerType,
		CountryCode:      countryCode,
		UniqueIdentifier: buildOwnerUniqueIdentifier(countryCode, ownerType, document),
	}
}

func assertOwnerIdentifierUnchanged(t *testing.T, owner *domain.Owner, want string) {
	t.Helper()
	if owner.UniqueIdentifier != want {
		t.Fatalf("owner unique identifier = %q, want %q", owner.UniqueIdentifier, want)
	}
}
