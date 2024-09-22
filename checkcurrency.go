package fiskalhrgo

import (
	"regexp"
)

// Helper function to validate if the string is a valid currency format (with 2 decimal places)
func IsValidCurrencyFormat(amount string) bool {
	// Regex pattern to match valid decimal with exactly two decimal places
	validCurrency := regexp.MustCompile(`^\d+(\.\d{2})$`)
	return validCurrency.MatchString(amount)
}
