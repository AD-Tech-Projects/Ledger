package domain

type CountryCode string

const (
	COUNTRY_US CountryCode = "US" // United States
	COUNTRY_BR CountryCode = "BR" // Brazil
	COUNTRY_EU CountryCode = "EU" // European Union

	// NOTE - unknown country code, stored as an empty string, used when there is no specific country association, such as crypto accounts
	COUNTRY_UNKNOWN CountryCode = ""
)
