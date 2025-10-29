package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFormatPhoneNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"+1234567890", "+1234567890"},
		{"1234567890", "+1234567890"},
		{"+91987654321", "+91987654321"},
		{"123456789012", "+123456789012"},
	}

	for _, test := range tests {
		result := formatPhoneNumber(test.input)
		if result != test.expected {
			t.Errorf("formatPhoneNumber(%s) = %s; want %s",
				test.input, result, test.expected)
		}
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"+1234567890", true},
		{"+123", false},
		{"abcdefghijk", false},
		{"+91987654321", true},
	}

	for _, test := range tests {
		result := validatePhoneNumber(test.input)
		if result != test.expected {
			t.Errorf("validatePhoneNumber(%s) = %v; want %v",
				test.input, result, test.expected)
		}
	}
}

func TestSearchHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/search?phone=+1234567890", nil)
	w := httptest.NewRecorder()

	searchHandler(w, req)

	if w.Code != http.StatusNotFound && w.Code != http.StatusOK {
		t.Errorf("Expected status code 404 or 200, got %v", w.Code)
	}
}

func TestAddNumberHandler(t *testing.T) {
	formData := strings.NewReader("phone_number=+1234567890&reason=test&store_location=test&check_amount=100")
	req := httptest.NewRequest("POST", "/add", formData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	addNumberHandler(w, req)

	if w.Code != http.StatusSeeOther && w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 303 or 400, got %v", w.Code)
	}
}
