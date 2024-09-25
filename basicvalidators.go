package fiskalhrgo

// SPDX-License-Identifier: MIT
// Copyright (c) 2024 L. D. T. d.o.o.
// Copyright (c) contributors for their respective contributions. See https://github.com/l-d-t/fiskalhrgo/graphs/contributors

import (
	"os"
	"path/filepath"
	"regexp"
)

// Helper function to validate if the string is a valid currency format (with 2 decimal places)
func IsValidCurrencyFormat(amount string) bool {
	// Regex pattern to match valid decimal with exactly two decimal places
	validCurrency := regexp.MustCompile(`^\d+(\.\d{2})$`)
	return validCurrency.MatchString(amount)
}

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

// ValidateLocationID validates the locationID
// It can contain only digits (0-9) and letters (a-z, A-Z), with a maximum length of 20.
func ValidateLocationID(locationID string) bool {
	// Regex pattern to match valid locationID
	validLocationID := regexp.MustCompile(`^[a-zA-Z0-9]{1,20}$`)
	return validLocationID.MatchString(locationID)
}

// IsFileReadable checks if the given file exists and is readable.
// It returns true if the file exists and is readable, otherwise false.
func IsFileReadable(filePath string) bool {
	// Get the absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return false
	}

	// Check if the file exists
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return false
	}

	// Check if the path is a regular file
	if !info.Mode().IsRegular() {
		return false
	}

	// Check if the file is readable
	file, err := os.Open(absPath)
	if err != nil {
		return false
	}
	defer file.Close()

	return true
}
