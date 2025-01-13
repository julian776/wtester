package wtester

import (
	"fmt"
	"strings"
)

// ValidationErrors is a struct that holds a list of
// [ExpectError] structs. This is used to build a
// custom error message.
type ValidationErrors struct {
	Errs []ExpectError
}

func (v *ValidationErrors) Error() string {
	var s string
	for _, e := range v.Errs {
		s += e.Error() + "\n"
	}

	return strings.Trim(s, "\n")
}

// IsEmpty returns true if there are no validation errors.
// Shortcut for len(ValidationErrors.Errs) == 0.
func (v *ValidationErrors) IsEmpty() bool {
	return len(v.Errs) == 0
}

// ExpectError is a struct that holds the title of
// the validation and a list of [ErrorRecords] that you
// can use to build a custom error message.
type ExpectError struct {
	Title  string
	Errors []ErrorRecord
}

func (v ExpectError) Error() string {
	errs := ""
	for _, e := range v.Errors {
		if len(e.Bytes) != 0 {
			errs += string(e.Bytes) + "\n"
		}

		if e.Err != nil {
			errs += e.Err.Error() + "\n"
		}
	}
	return fmt.Sprintf("validation \"%s\"\nFails On:\n%s", v.Title, errs)
}

// ErrorRecord is a struct that holds the bytes that
// failed validation and the error that was returned.
type ErrorRecord struct {
	Bytes []byte
	Err   error
}

func (e ErrorRecord) Error() string {
	errs := ""
	if e.Err != nil {
		errs += e.Err.Error() + "\n"
	}

	if len(e.Bytes) != 0 {
		errs += string(e.Bytes) + "\n"
	}

	return errs
}
