package wtester

import "sync"

type Expect struct {
	title   string
	f       ExpectFunc
	every   bool
	min     int
	max     int
	matches int
	mu      sync.Mutex // guards matches
}

func NewExpect(title string, f ExpectFunc) *Expect {
	return &Expect{
		title: title,
		f:     f,
		min:   1,
	}
}

// WithMin sets the minimum number of times the expectation should match
// the default is 1.
func (e *Expect) WithMin(min int) *Expect {
	e.min = min
	return e
}

// WithMax sets the maximum number of times the expectation should match
func (e *Expect) WithMax(max int) *Expect {
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
