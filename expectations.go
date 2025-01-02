package wtester

import (
	"bytes"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type ExpectFunc func(expected []byte) bool

func StringMatch(expected string, exact bool) ExpectFunc {
	return func(actual []byte) bool {
		if exact {
			return expected == string(actual)
		}

		return strings.Contains(string(actual), expected)
	}
}

func PrefixMatch(expected string) ExpectFunc {
	return func(actual []byte) bool {
		return strings.HasPrefix(string(actual), expected)
	}
}

func SuffixMatch(expected string) ExpectFunc {
	return func(actual []byte) bool {
		return strings.HasSuffix(string(actual), expected)
	}
}

func ValidUTF8() ExpectFunc {
	return func(actual []byte) bool {
		return utf8.Valid(actual)
	}
}

func RegexMatch(pattern string) ExpectFunc {
	re := regexp.MustCompile(pattern)
	return func(actual []byte) bool {
		return re.Match(actual)
	}
}

func BytesMatch(expected []byte) ExpectFunc {
	return func(actual []byte) bool {
		return bytes.Equal(expected, actual)
	}
}

func RunesMatch(expected []rune) ExpectFunc {
	return func(actual []byte) bool {
		return slices.Equal(bytes.Runes(actual), expected)
	}
}
