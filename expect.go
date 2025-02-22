package wtester

import "sync"

type Expect struct {
	title   string
	exp     Expecter
	every   bool
	min     int
	max     int
	noMatch bool
	matches int
	mu      sync.Mutex // guards matches
}

func NewExpect(title string, exp Expecter) *Expect {
	return &Expect{
		title: title,
		exp:   exp,
		min:   1,
	}
}

// WithMin sets the minimum number of times the expectation should match
// the default is 1.
func (e *Expect) WithMin(min int) *Expect {
	// If noMatch is set, we don't want to set a min
	if e.noMatch {
		return e
	}

	e.min = min
	return e
}

// WithMax sets the maximum number of times the expectation should match.
// Defaults to 0, which means no maximum. But if the max is set explicitly
// to 0, the expect will assert '0' matches. Also overrides min to 0.
func (e *Expect) WithMax(max int) *Expect {
	if max == 0 {
		e.noMatch = true
		e.min = 0
	}

	e.max = max
	return e
}

// Every sets the expectation that should match every time
// the Write method is called. If the expectation does not
// match, a validation error is recorded.
func (e *Expect) Every() *Expect {
	e.every = true
	return e
}

func (e *Expect) matched() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.matches++
}
