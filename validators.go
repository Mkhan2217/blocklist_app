// validators.go — phone formatting and validation helpers.
//
// Keep normalisation and validation logic here so handlers remain small.
//
// formatPhoneNumber: normalizes input (URL-decoded) by removing non-digit
// characters except '+', and ensures a leading '+' is present if digits found.
//
// validatePhoneNumber: enforces the DB CHECK regex: /^\+[1-9][0-9]{9,14}$/.
package main

import (
	"log"
	"net/url"
	"regexp"
	"strings"
)

// phoneRegex matches DB CHECK constraint: + followed by 10-15 digits,
// first digit after + must be non-zero.
var phoneRegex = regexp.MustCompile(`^\+[1-9][0-9]{9,14}$`)

// formatPhoneNumber returns normalized phone or empty string when input empty.
// It preserves an existing country code; only strips noise characters.
func formatPhoneNumber(phone string) string {
	decoded, err := url.QueryUnescape(phone)
	if err == nil {
		phone = decoded
	}

	// keep digits and plus
	reg := regexp.MustCompile(`[^\d+]`)
	cleaned := reg.ReplaceAllString(phone, "")

	if cleaned == "" {
		return ""
	}
	if !strings.HasPrefix(cleaned, "+") {
		cleaned = "+" + cleaned
	}

	// keep a single info log for visibility in dev; remove/redirect in prod.
	log.Printf("Formatted phone: %s", cleaned)
	return cleaned
}

// validatePhoneNumber returns true when phone matches expected international format.
func validatePhoneNumber(phone string) bool {
	return phoneRegex.MatchString(phone)
}
