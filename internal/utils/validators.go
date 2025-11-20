package utils

import (
	"log"
	"net/url"
	"regexp"
	"strings"
)

var phoneRegex = regexp.MustCompile(`^\+[1-9][0-9]{9,14}$`)

func FormatPhoneNumber(phone string) string {
	decoded, _ := url.QueryUnescape(phone)
	phone = decoded

	reg := regexp.MustCompile(`[^\d+]`)
	cleaned := reg.ReplaceAllString(phone, "")

	if cleaned == "" {
		return ""
	}
	if !strings.HasPrefix(cleaned, "+") {
		cleaned = "+" + cleaned
	}

	log.Printf("Formatted phone: %s", cleaned)
	return cleaned
}

func ValidatePhoneNumber(phone string) bool {
	return phoneRegex.MatchString(phone)
}
