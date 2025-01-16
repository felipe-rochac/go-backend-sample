package common

import (
	"testing"
)

func Test_IsValidEmail_ExpectSuccess(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name+tag+sorting@example.com", true},
		{"user.name@example.co.uk", true},
		{"invalid-email", false},
		{"@missingusername.com", false},
		{"missingatsign.com", false},
	}

	for _, test := range tests {
		result := IsValidEmail(test.email)
		if result != test.expected {
			t.Errorf("IsValidEmail(%v) = %v; want %v", test.email, result, test.expected)
		}
	}
}

func Test_IsValidUuid_ExpectSucces(t *testing.T) {
	tests := []struct {
		uuidStr  string
		expected bool
	}{
		{"123e4567-e89b-12d3-a456-426614174000", true},
		{"invalid-uuid", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsValidUuid(test.uuidStr)
		if result != test.expected {
			t.Errorf("IsValidUuid(%v) = %v; want %v", test.uuidStr, result, test.expected)
		}
	}
}

func Test_StringMinMaxLength_ExpectSuccess(t *testing.T) {
	tests := []struct {
		text     string
		min      int
		max      int
		expected bool
	}{
		{"hello", 1, 10, true},
		{"hello", 6, 10, false},
		{"hello", 1, 4, false},
		{"", 1, 10, false},
		{"hello world", 1, 10, false},
	}

	for _, test := range tests {
		result := StringMinMaxLength(test.text, test.min, test.max)
		if result != test.expected {
			t.Errorf("StringMinMaxLength(%v, %v, %v) = %v; want %v", test.text, test.min, test.max, result, test.expected)
		}
	}
}
