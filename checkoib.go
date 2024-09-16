package fiskalhrgo

// ValidateOIB checks if an OIB is valid using the Mod 11, 10 algorithm
func ValidateOIB(oib string) bool {
	if len(oib) != 11 {
		return false
	}

	// Convert the first 10 digits of OIB to integers
	var remainder int = 10
	for i := 0; i < 10; i++ {
		digit := int(oib[i] - '0') // Convert char to int by subtracting ASCII '0'
		if digit < 0 || digit > 9 {
			return false // If the character is not a digit, return false
		}
		remainder = (remainder + digit) % 10
		if remainder == 0 {
			remainder = 10
		}
		remainder = (remainder * 2) % 11
	}

	// Calculate the check digit
	checkDigit := (11 - remainder) % 10

	// Compare the calculated check digit with the last digit of the OIB
	lastDigit := int(oib[10] - '0')
	if lastDigit < 0 || lastDigit > 9 {
		return false
	}

	return checkDigit == lastDigit
}
