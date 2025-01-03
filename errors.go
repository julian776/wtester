package wtester

import "fmt"

// ValidationErrors is a struct that holds a list of
// ValidationError structs. This is used to build a
// custom error message.
type ValidationErrors struct {
	Errs []ValidationError
}

func (v ValidationErrors) Error() string {
	var s string
	for _, e := range v.Errs {
		s += e.Error() + "\n"
	}

	return s
}

// IsEmpty returns true if there are no validation errors.
// Shortcut for len(ValidationErrors.Errs) == 0.
func (v *ValidationErrors) IsEmpty() bool {
	return len(v.Errs) == 0
}

// ValidationError is a struct that holds the title of
// the validation and a list of ErrorRecords that you
// can use to build a custom error message.
type ValidationError struct {
	Title  string
	Errors []ErrorRecord
}

func (v ValidationError) Error() string {
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
