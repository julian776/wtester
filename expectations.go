package wtester

import (
	"bytes"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type Expecter interface {
	Expect(actual []byte) bool
}

// JSONExpecter is an interface for JSON expectations.
// If implemented, the expectation will be unmarshaled into a map
// before being passed to ExpectJSON instead of calling the Expect method.
// Improves performance by unmarshaling the JSON data only once per Write call.
type JSONExpecter interface {
	ExpectJSON(actual map[string]any) bool
}

type ExpectFunc func(actual []byte) bool

func (f ExpectFunc) Expect(actual []byte) bool {
	return f(actual)
}

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
	return utf8.Valid
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

type obfuscatedMatch struct {
	obfuscateChar        string
	percentageObfuscated float64
	fields               []string
}

// ObfuscatedMatch returns an Expecter that checks if the provided fields
// in a JSON are obfuscated with the provided character and percentage.
//
// Any value that corresponds to a field in the fields slice is expected to
// be a string. If it is not a string, the function returns false.
// An empty value for any field will return false.
//
// Panics if the percentageObfuscated is not between 0 and 1 or if
// the fields slice is empty.
func ObfuscatedMatch(
	obfuscateChar string,
	percentageObfuscated float64,
	fields ...string,
) Expecter {
	if percentageObfuscated < 0 || percentageObfuscated > 1 {
		panic("percentageObfuscated must be between 0 and 1")
	}

	if len(fields) == 0 {
		panic("fields cannot be empty")
	}

	return &obfuscatedMatch{
		obfuscateChar:        obfuscateChar,
		percentageObfuscated: percentageObfuscated,
		fields:               fields,
	}
}

// Only for satisfy the Expecter interface.
func (om *obfuscatedMatch) Expect(actual []byte) bool {
	return false
}

// ExpectJSON checks if the provided fields in a JSON are obfuscated with the
// provided character and percentage.
func (om *obfuscatedMatch) ExpectJSON(m map[string]any) bool {
	for k, v := range m {
		if !slices.Contains(om.fields, k) {
			continue
		}

		str, ok := v.(string)
		if !ok {
			return false
		}

		obfuscated := strings.Count(str, om.obfuscateChar)
		total := len(str)
		percent := float64(obfuscated) / float64(total)

		if percent >= om.percentageObfuscated {
			return true
		}
	}

	return false
}

// Not returns an ExpectFunc that negates the result of the provided ExpectFunc.
func Not(expectFunc ExpectFunc) ExpectFunc {
	return func(actual []byte) bool {
		return !expectFunc(actual)
	}
}

// AndMatch returns an ExpectFunc that checks if all provided ExpectFuncs return true.
func AndMatch(expectFuncs ...ExpectFunc) ExpectFunc {
	return func(actual []byte) bool {
		for _, expectFunc := range expectFuncs {
			if !expectFunc(actual) {
				return false
			}
		}

		return true
	}
}

// OrMatch returns an ExpectFunc that checks if any of the provided ExpectFuncs return true.
func OrMatch(expectFuncs ...ExpectFunc) ExpectFunc {
	return func(actual []byte) bool {
		for _, expectFunc := range expectFuncs {
			if expectFunc(actual) {
				return true
			}
		}

		return false
	}
}
