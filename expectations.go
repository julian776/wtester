package wtester

import (
	"bytes"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type ExpectFunc func(expected []byte) bool

// StringMatch returns an ExpectFunc that checks if the actual byte slice matches
// the expected string. If exact is true, it checks for an exact match. Otherwise,
// it checks if the expected string is contained within the actual byte slice.
func StringMatch(expected string, exact bool) ExpectFunc {
	return func(actual []byte) bool {
		if exact {
			return expected == string(actual)
		}

		return strings.Contains(string(actual), expected)
	}
}

// PrefixMatch returns an ExpectFunc that checks if the actual byte slice
// starts with the expected string.
func PrefixMatch(expected string) ExpectFunc {
	return func(actual []byte) bool {
		return strings.HasPrefix(string(actual), expected)
	}
}

// SuffixMatch returns an ExpectFunc that checks if the actual byte slice
// ends with the expected string.
func SuffixMatch(expected string) ExpectFunc {
	return func(actual []byte) bool {
		return strings.HasSuffix(string(actual), expected)
	}
}

// ValidUTF8 returns an ExpectFunc that checks if the actual byte slice
// is valid UTF-8.
func ValidUTF8() ExpectFunc {
	return func(actual []byte) bool {
		return utf8.Valid(actual)
	}
}

// RegexMatch returns an ExpectFunc that checks for matches against the
// provided regular expression pattern.
func RegexMatch(pattern string) ExpectFunc {
	re := regexp.MustCompile(pattern)
	return func(actual []byte) bool {
		return re.Match(actual)
	}
}

// BytesMatch returns an ExpectFunc that checks if the actual byte slice
// matches the expected byte slice.
func BytesMatch(expected []byte) ExpectFunc {
	return func(actual []byte) bool {
		return bytes.Equal(expected, actual)
	}
}

// RunesMatch returns an ExpectFunc that checks if the actual byte slice
// matches the expected rune slice.
func RunesMatch(expected []rune) ExpectFunc {
	return func(actual []byte) bool {
		return slices.Equal(bytes.Runes(actual), expected)
	}
}
