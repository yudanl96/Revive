package validate

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile("^[a-zA-Z0-9_]+$").MatchString
	isValidPassword = regexp.MustCompile("^[a-zA-Z0-9_!@$%*&]+$").MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 5, 15); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only letters, digits, or underscore")
	}
	return nil
}

func ValidatePassword(value string) error {
	if err := ValidateString(value, 8, 30); err != nil {
		return err
	}
	if !isValidPassword(value) {
		return fmt.Errorf("password must contain letter, digits, or special character _!@$%%*&")
	}
	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("not a valid email address")
	}
	return nil
}
