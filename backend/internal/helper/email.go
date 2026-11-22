package helper

import "strings"

// disposableDomains contains a hash set of known temporary email provider domains.
// Using a map with empty structs provides O(1) lookup time and optimal memory usage.
var disposableDomains = map[string]struct{}{
	"mailinator.com":    {},
	"niepodam.pl":       {},
	"yopmail.com":       {},
	"tempmail.com":      {},
	"10minutemail.com":  {},
	"guerrillamail.com": {},
	"trashmail.com":     {},
	// Add more domains as needed
}

// IsDisposableEmail checks whether the given email address belongs to a known
// disposable or temporary email provider.
func IsDisposableEmail(email string) bool {
	email = CleanEmail(email)

	_, domain := extractEmailParts(email)
	if domain == "" {
		return false // Invalid email format or missing domain
	}

	_, isDisposable := disposableDomains[domain]
	return isDisposable
}

// NormalizeEmail converts an email address to its canonical form
// by removing provider-specific sub-addressing (aliases) and dots.
func NormalizeEmail(email string) string {
	email = CleanEmail(email)

	localPart, domain := extractEmailParts(email)
	if domain == "" {
		return email // Invalid email structure, return as-is
	}

	switch domain {
	case "gmail.com", "googlemail.com":
		// Google ignores dots and everything after the first '+' sign
		if plusIndex := strings.Index(localPart, "+"); plusIndex != -1 {
			localPart = localPart[:plusIndex]
		}
		localPart = strings.ReplaceAll(localPart, ".", "")
		domain = "gmail.com" // Standardize domain alias to gmail.com

	case "outlook.com", "hotmail.com", "live.com", "msn.com", "icloud.com", "me.com", "mac.com", "proton.me", "protonmail.com":
		// Providers supporting plus-addressing (+) while preserving dots
		if plusIndex := strings.Index(localPart, "+"); plusIndex != -1 {
			localPart = localPart[:plusIndex]
		}

	case "yahoo.com", "ymail.com":
		// Yahoo uses hyphens (-) for disposable email address aliases (e.g., user-alias@yahoo.com)
		if hyphenIndex := strings.Index(localPart, "-"); hyphenIndex != -1 {
			localPart = localPart[:hyphenIndex]
		}
	}

	return localPart + "@" + domain
}

// CleanEmail strips whitespace and converts the email to lower case.
func CleanEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// extractEmailParts extracts the local part and domain from an email address.
// It is used internally across email-related functions.
func extractEmailParts(email string) (string, string) {
	atPosition := strings.LastIndex(email, "@")
	if atPosition == -1 || atPosition == 0 || atPosition == len(email)-1 {
		return "", "" // Invalid email format
	}
	return email[:atPosition], email[atPosition+1:]
}
