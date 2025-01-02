package wtester

import "fmt"

type ValidationErrors struct {
	errs []ValidationError
}

func (v *ValidationErrors) Errors() string {
	var s string
	for _, e := range v.errs {
		s += e.Error() + "\n"
	}

	return s
}

type ValidationError struct {
	Title  string
	Errors []ErrorRecord
}

func (v *ValidationError) Error() string {
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
