package client

import (
	"testing"
)

func TestValidateRGBHex(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"#FF0000", true},
		{"#ff0000", true},
		{"#f00", false},
		{"FF0000", false},
		{"#GG0000", false},
		{"#1234567", false}, // Too long
		{"#12345", false},   // Too short
		{"123456", false},   // Missing #
	}

	for _, tc := range testCases {
		t.Run(
			tc.input, func(t *testing.T) {
				actual := validateRGBHex(tc.input)
				if actual != tc.expected {
					t.Errorf("isValidRGBHex(%q) = %v; want %v", tc.input, actual, tc.expected)
				}
			},
		)
	}
}
