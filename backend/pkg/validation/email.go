package validation

import (
	"errors"
	"regexp"
)

const minPasswordLength = 6

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
var letterRegex = regexp.MustCompile(`[A-Za-z]`)
var digitRegex = regexp.MustCompile(`\d`)
var specialCharRegex = regexp.MustCompile(`[^A-Za-z\d]`)

func IsEmailValid(email string) error {
	if len(email) == 0 {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func IsPasswordValid(pw string) error {
	err := errors.New("password must be at least 6 characters and contain a letter, number, and special character")
	if len(pw) < minPasswordLength {
		return err
	}

	if !letterRegex.MatchString(pw) || !digitRegex.MatchString(pw) || !specialCharRegex.MatchString(pw) {
		return err
	}
	return nil
}
