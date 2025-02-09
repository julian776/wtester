package wtester

import (
	"fmt"
	"testing"
)

func ExampleObfuscatedMatch() {
	// Create a match function that checks if the "credit_card" field is obfuscated.
	// Percentage obfuscated is set to 0.4, meaning at least 40%
	// The percentage could be calculated by counting the number of obfuscateChar
	// in the string divided by the total number of characters in the string.
	// For example, "1234-****-****-1234" has 8 obfuscated characters out of
	// 19 total characters (The '-' character is not obfuscated).
	// (8 / 19) = 0.42105263
	matchFunc := ObfuscatedMatch("*", 0.4, "credit_card")

	// Test the match function against a JSON object.
	input := []byte(`{"credit_card": "1234-****-****-1234"}`)
	if matchFunc(input) {
		fmt.Println("The credit card field is obfuscated.")
	} else {
		fmt.Println("The credit card field is not obfuscated.")
	}

	// Output: The credit card field is obfuscated.
}

func TestObfuscatedMatch(t *testing.T) {
	tests := map[string]struct {
		obfuscateChar        string
		percentageObfuscated float64
		fields               []string
		input                []byte
		expected             bool
		panic                bool
	}{
		"Field is sufficiently obfuscated": {
			obfuscateChar:        "*",
			percentageObfuscated: 0.5,
			fields:               []string{"password"},
			input:                []byte(`{"password": "p******d"}`),
			expected:             true,
		},
		"Field is obfuscated credit card": {
			obfuscateChar:        "*",
			percentageObfuscated: 0.5,
			fields:               []string{"credit_card"},
			input:                []byte(`{"credit_card": "**34-****-****-12**"}`),
			expected:             true,
		},
		"Field is not sufficiently obfuscated": {
			obfuscateChar:        "*",
			percentageObfuscated: 1,
			fields:               []string{"password"},
			input:                []byte(`{"password": "*******d"}`),
			expected:             false,
		},
		"Field is not a string": {
			obfuscateChar:        "*",
			percentageObfuscated: 0.5,
			fields:               []string{"password"},
			input:                []byte(`{"password": 123456}`),
			expected:             false,
		},
		"Field is missing": {
			obfuscateChar:        "*",
			percentageObfuscated: 0.5,
			fields:               []string{"password"},
			input:                []byte(`{"username": "user1"}`),
			expected:             false,
		},
		"Multiple fields, one sufficiently obfuscated": {
			obfuscateChar:        "*",
			percentageObfuscated: 0.5,
			fields:               []string{"password", "token"},
			input:                []byte(`{"password": "p******d", "token": "t****n"}`),
			expected:             true,
		},
		"Multiple fields, none sufficiently obfuscated": {
			obfuscateChar:        "*",
			percentageObfuscated: 0.8,
			fields:               []string{"password", "token"},
			input:                []byte(`{"password": "p******d", "token": "t****n"}`),
			expected:             false,
		},
		"Empty fields slice": {
			obfuscateChar:        "*",
			percentageObfuscated: 0.5,
			fields:               []string{},
			input:                []byte(`{"password": "p******d"}`),
			panic:                true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.panic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("ObfuscatedMatch() did not panic")
					}
				}()
			}

			matchFunc := ObfuscatedMatch(tt.obfuscateChar, tt.percentageObfuscated, tt.fields...)
			if got := matchFunc(tt.input); got != tt.expected {
				t.Errorf("ObfuscatedMatch() = %v, want %v", got, tt.expected)
			}
		})
	}
}
