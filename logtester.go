package wtester

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

// WTester is a wrapper around an [io.Writer] that allows
// for expectations to be set on the output of the writer.
// It can be used to test loggers, writers, or any other
// [io.Writer] implementation.
type WTester struct {
	w       io.Writer
	muW     sync.Mutex // guards w
	expects map[string]*Expect
	errors  map[string]*ExpectError
	muErr   sync.Mutex
}

func NewWTester(w io.Writer) *WTester {
	return &WTester{
		w:       w,
		expects: make(map[string]*Expect),
		errors:  make(map[string]*ExpectError),
	}
}

// AppendWriter appends an [io.Writer] to the WTester's
// underlying writer. This allows for multiple writers to
// be written to simultaneously.
func (l *WTester) AppendWriter(w io.Writer) {
	l.muW.Lock()
	defer l.muW.Unlock()

	if l.w == nil {
		l.w = w
		return
	}

	l.w = io.MultiWriter(l.w, w)
}

// Write writes the provided byte slice to the underlying
// [io.Writer] and checks if the byte slice matches any of
// the expectations set on the WTester.
func (l *WTester) Write(p []byte) (n int, err error) {
	// Only unmarshal JSON once. And only if there are JSON expectations.
	var m map[string]any

	for _, e := range l.expects {
		var ok bool
		switch exp := e.exp.(type) {
		case JSONExpecter:
			if m == nil {
				if err := json.Unmarshal(p, &m); err != nil {
					l.appendError(e.title, ErrorRecord{
						Bytes: p,
						Err:   fmt.Errorf("failed to unmarshal JSON: %s", err.Error()),
					})
					continue
				}
			}

			ok = exp.ExpectJSON(m)
		default:
			ok = exp.Expect(p)
		}

		if ok {
			e.matched()
			continue
		}

		if !ok && e.every {
			l.appendError(e.title, ErrorRecord{
				Bytes: p,
			})
		}
	}

	l.muW.Lock()
	defer l.muW.Unlock()

	return l.w.Write(p)
}

// Close closes the underlying io.Writer if it implements
// the [io.Closer] interface.
func (l *WTester) Close() error {
	if c, ok := l.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// Expect sets an expectation on the WTester. The title
// is used to identify the expectation and the f parameter
// is a function that takes a byte slice and returns a boolean
// indicating if the expectation is met.
func (l *WTester) Expect(title string, exp Expecter) *Expect {
	e := NewExpect(title, exp)
	l.expects[title] = e

	return e
}

// ExpectFunc sets a function expectation on the WTester.
// The title is used to identify the expectation and the f
// parameter is a function that takes a byte slice and returns
// a boolean indicating if the expectation is met.
func (l *WTester) ExpectFunc(title string, f func(actual []byte) bool) *Expect {
	return l.Expect(title, ExpectFunc(f))
}

// Reset resets the WTester by clearing all expectations
// and errors.
func (l *WTester) Reset() {
	l.muW.Lock()
	defer l.muW.Unlock()

	l.muErr.Lock()
	defer l.muErr.Unlock()

	l.expects = make(map[string]*Expect)
	l.errors = make(map[string]*ExpectError)
}

// Validate validates the expectations set on the WTester
// and returns an error if any of the expectations are not met.
// If there are no validation errors, nil is returned.
// You must cast the err ve, ok := err.(*ValidationErrors) to access
// the underlying validation errors.
func (l *WTester) Validate() error {
	for _, e := range l.expects {
		switch {
		case e.min > 0 && e.matches < e.min:
			l.appendError(e.title, ErrorRecord{
				Err: fmt.Errorf("expected at least %d matches, got %d", e.min, e.matches),
			})
		case (e.max > 0 || e.noMatch) && e.matches > e.max:
			l.appendError(e.title, ErrorRecord{
				Err: fmt.Errorf("expected at most %d matches, got %d", e.max, e.matches),
			})
		}
	}

	ve := &ValidationErrors{}
	for _, e := range l.errors {
		ve.Errs = append(ve.Errs, *e)
	}

	if ve.IsEmpty() {
		return nil
	}

	return ve
}

// appendError appends an error to the WTester's error map.
func (l *WTester) appendError(title string, e ErrorRecord) {
	l.muErr.Lock()
	defer l.muErr.Unlock()

	if l.errors[title] == nil {
		l.errors[title] = &ExpectError{
			Title: title,
		}
	}

	l.errors[title].Errors = append(l.errors[title].Errors, e)
}
