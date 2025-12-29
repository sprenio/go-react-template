package validation

import "regexp"

const minPasswordLength = 6

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
var letterRegex = regexp.MustCompile(`[A-Za-z]`)
var digitRegex = regexp.MustCompile(`\d`)
var specialCharRegex = regexp.MustCompile(`[^A-Za-z\d]`)


func IsEmailValid(email string) bool {
	return emailRegex.MatchString(email)
}

func IsPasswordValid(pw string) bool {
	if len(pw) < minPasswordLength {
		return false
	}

	return letterRegex.MatchString(pw) && digitRegex.MatchString(pw) && specialCharRegex.MatchString(pw)
}