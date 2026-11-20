package helpers

import (
	"strings"
	"unicode"
)

var (
	cpfFirstCheckDigitWeights   = []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
	cpfSecondCheckDigitWeights  = []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
	cnpjFirstCheckDigitWeights  = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	cnpjSecondCheckDigitWeights = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
)

// NormalizeCpf removes all non-alphanumeric characters from a CPF and returns the normalized string
func NormalizeCpf(cpf string) string {
	normalizedCpf := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, cpf)

	return normalizedCpf
}

// IsValidCPF checks a CPF of 9 alphanumeric characters followed by two check digits
// NOTE - like the CNPJ, the CPF may have uppercase letters before the check digits; each character is worth its ASCII
// code minus 48 (digits keep their value, A = 17), so the numeric CPF is validated by the same rule
func IsValidCPF(document string) bool {
	isWellFormed := len(document) == 11 && isUpperAlphanumeric(document[:9]) && isDigits(document[9:]) &&
		!hasSingleRepeatedCharacter(document)
	if !isWellFormed {
		return false
	}

	values := documentValues(document)
	isFirstCheckDigitValid := mod11CheckDigit(values[:9], cpfFirstCheckDigitWeights) == values[9]
	isSecondCheckDigitValid := mod11CheckDigit(values[:10], cpfSecondCheckDigitWeights) == values[10]
	return isFirstCheckDigitValid && isSecondCheckDigitValid
}

// IsValidCNPJ checks a CNPJ of 12 alphanumeric characters followed by two check digits
// NOTE - since July 2026 the CNPJ may have uppercase letters in its first 12 characters; each character is worth its
// ASCII code minus 48 (digits keep their value, A = 17), so the numeric CNPJ is validated by the same rule
func IsValidCNPJ(document string) bool {
	isWellFormed := len(document) == 14 && isUpperAlphanumeric(document[:12]) && isDigits(document[12:]) &&
		!hasSingleRepeatedCharacter(document)
	if !isWellFormed {
		return false
	}

	values := documentValues(document)
	isFirstCheckDigitValid := mod11CheckDigit(values[:12], cnpjFirstCheckDigitWeights) == values[12]
	isSecondCheckDigitValid := mod11CheckDigit(values[:13], cnpjSecondCheckDigitWeights) == values[13]
	return isFirstCheckDigitValid && isSecondCheckDigitValid
}

// IsValidSSN checks the 9 digits of an SSN and rejects the numbers never issued: area 000, 666 or 900-999, group 00 and serial 0000
// NOTE - ITINs (area 9XX) are not accepted as the document of individuals
func IsValidSSN(document string) bool {
	isWellFormed := len(document) == 9 && isDigits(document)
	if !isWellFormed {
		return false
	}

	area, group, serial := document[:3], document[3:5], document[5:]
	isIssuedArea := area != "000" && area != "666" && area[0] != '9'
	return isIssuedArea && group != "00" && serial != "0000"
}

// IsValidEIN checks the 9 digits of an EIN
func IsValidEIN(document string) bool {
	return len(document) == 9 && isDigits(document) && document != "000000000"
}

// mod11CheckDigit is the check digit of the values with the weights: 0 when the remainder of the weighted sum by 11 is
// lower than 2, otherwise 11 minus the remainder
func mod11CheckDigit(values []int, weights []int) int {
	sum := 0
	for index, value := range values {
		sum += value * weights[index]
	}

	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}

// documentValues converts each character to its ASCII code minus 48 (the digit value for digits)
func documentValues(document string) []int {
	values := make([]int, len(document))
	for index, character := range []byte(document) {
		values[index] = int(character) - '0'
	}
	return values
}

func isDigits(value string) bool {
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return value != ""
}

func isUpperAlphanumeric(value string) bool {
	for _, character := range value {
		isDigit := character >= '0' && character <= '9'
		isUpperLetter := character >= 'A' && character <= 'Z'
		if !isDigit && !isUpperLetter {
			return false
		}
	}
	return value != ""
}

func hasSingleRepeatedCharacter(value string) bool {
	return value != "" && strings.Count(value, value[:1]) == len(value)
}
