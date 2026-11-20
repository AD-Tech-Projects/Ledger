package services

import (
	"fmt"
	"strings"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/services/helpers"
)

type ownerValidationError struct {
	message string
}

func (e *ownerValidationError) Error() string {
	return "OwnerValidationError:" + e.message
}

func (e *ownerValidationError) Unwrap() error {
	return ErrInvalidOwner
}

// validateOwner checks an owner before its creation: the type, a known country and the unique identifier, which must be
// COUNTRY:TYPE:DOCUMENT with a document whose format is valid for the country and the owner type
func validateOwner(owner *domain.Owner) error {
	if owner == nil {
		return &ownerValidationError{message: "owner is required"}
	}

	documentFormats, isSupportedCountry := ownerDocumentFormats[owner.CountryCode]
	if !isSupportedCountry {
		return &ownerValidationError{message: fmt.Sprintf("owner country code %q is not supported", owner.CountryCode)}
	}
	isValidDocumentFormat, isSupportedType := documentFormats[owner.Type]
	if !isSupportedType {
		return &ownerValidationError{message: fmt.Sprintf("owner type %q is invalid", owner.Type)}
	}

	parts := strings.Split(owner.UniqueIdentifier, ":")
	isIdentifierOfOwner := len(parts) == 3 && parts[0] == string(owner.CountryCode) && parts[1] == string(owner.Type)
	if !isIdentifierOfOwner {
		return &ownerValidationError{message: "owner unique identifier must be COUNTRY:TYPE:DOCUMENT with the owner country code and type"}
	}

	document := parts[2]
	if document == "" {
		return &ownerValidationError{message: "owner document is required"}
	}
	if !isValidDocumentFormat(document) {
		return &ownerValidationError{message: fmt.Sprintf("owner document is not a valid %s %s document", owner.CountryCode, owner.Type)}
	}

	return nil
}

// ownerDocumentFormats defines the supported countries and owner types, and how the normalized document of each one is validated
// NOTE - COUNTRY_UNKNOWN is not supported for new owners, since the document format depends on the country
var ownerDocumentFormats = map[domain.CountryCode]map[domain.OwnerType]func(document string) bool{
	domain.COUNTRY_BR: {
		domain.OWNER_INDIVIDUAL: helpers.IsValidCPF,
		domain.OWNER_COMPANY:    helpers.IsValidCNPJ,
	},
	domain.COUNTRY_US: {
		domain.OWNER_INDIVIDUAL: helpers.IsValidSSN,
		domain.OWNER_COMPANY:    helpers.IsValidEIN,
	},
	// NOTE - the European Union has no common personal or company document, so any normalized document is accepted
	domain.COUNTRY_EU: {
		domain.OWNER_INDIVIDUAL: helpers.IsNotEmpty,
		domain.OWNER_COMPANY:    helpers.IsNotEmpty,
	},
}
