package utils

import "testing"

func TestFormatPhoneNumber(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"+1234567890", "+1234567890"},
		{"1234567890", "+1234567890"},
		{"(123) 456-7890", "+1234567890"},
		{"abc1234567890xyz", "+1234567890"},
		{"", ""},
	}
	

	for _, tt := range tests {
		got := FormatPhoneNumber(tt.input)
		if got != tt.output {
			t.Errorf("FormatPhoneNumber(%q) = %q; want %q", tt.input, got, tt.output)
		}
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	validNumbers := []string{
		"+1234567890",
		"+112345678901",
		"+911234567890",
	}
	invalidNumbers := []string{
		"1234567890",
		"+0123456789",
		"+123",
		"",
		"abcdef",
	}

	for _, num := range validNumbers {
		if !ValidatePhoneNumber(num) {
			t.Errorf("ValidatePhoneNumber(%q) = false; want true", num)
		}
	}

	for _, num := range invalidNumbers {
		if ValidatePhoneNumber(num) {
			t.Errorf("ValidatePhoneNumber(%q) = true; want false", num)
		}
	}
}
