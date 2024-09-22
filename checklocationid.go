package fiskalhrgo

import (
	"regexp"
)

// ValidateLocationID validates the locationID
// It can contain only digits (0-9) and letters (a-z, A-Z), with a maximum length of 20.
func ValidateLocationID(locationID string) bool {
	// Regex pattern to match valid locationID
	validLocationID := regexp.MustCompile(`^[a-zA-Z0-9]{1,20}$`)
	return validLocationID.MatchString(locationID)
}
