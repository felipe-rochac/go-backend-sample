package common

import (
	"regexp"

	"github.com/google/uuid"
)

func IsValidEmail(email string) bool {
	// Regular expression for validating an email
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex
	re := regexp.MustCompile(regex)

	// Match the email string
	return re.MatchString(email)
}

func IsValidUuid(idStr string) bool {
	_, err := uuid.Parse(idStr)

	return err == nil
}

func StringMinMaxLength(text string, min, max int) bool {
	if len(text) < min {
		return false
	}

	if len(text) > max {
		return false
	}

	return true
}
