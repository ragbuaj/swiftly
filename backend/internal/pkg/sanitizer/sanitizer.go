package sanitizer

import (
	"html"
	"regexp"
	"strings"
)

// Text cleans a string by escaping HTML tags and trimming whitespace.
func Text(input string) string {
	return html.EscapeString(strings.TrimSpace(input))
}

// Email normalizes an email address.
func Email(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// Username cleans and normalizes a username (lowercase, no spaces).
func Username(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

// Phone removes non-numeric characters from a phone number string.
func Phone(phone string) string {
	re := regexp.MustCompile(`[^0-9+]`)
	return re.ReplaceAllString(phone, "")
}
