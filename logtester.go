package wtester

import (
	"fmt"
	"io"
	"sync"
)

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

func (l *WTester) Write(p []byte) (n int, err error) {
	for _, e := range l.expects {
		ok := e.f(p)
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

func (l *WTester) Close() error {
	if c, ok := l.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func (l *WTester) Expect(title string, f ExpectFunc) *Expect {
	e := NewExpect(title, f)
	l.expects[title] = e

	return e
}

func (l *WTester) Reset() {
	l.muW.Lock()
	defer l.muW.Unlock()

	l.muErr.Lock()
	defer l.muErr.Unlock()

	l.expects = make(map[string]*Expect)
	l.errors = make(map[string]*ExpectError)
}

func (l *WTester) Validate() error {
	for _, e := range l.expects {
		if e.min > 0 && e.matches < e.min {
			l.appendError(e.title, ErrorRecord{
				Err: fmt.Errorf("expected at least %d matches, got %d", e.min, e.matches),
			})
		}

		if e.max > 0 && e.matches > e.max {
			l.appendError(e.title, ErrorRecord{
				Err: fmt.Errorf("expected at most %d matches, got %d", e.max, e.matches),
			})
		}
	}

	ve := ValidationErrors{}
	for _, e := range l.errors {
		ve.Errs = append(ve.Errs, *e)
	}

	if ve.IsEmpty() {
		return nil
	}

	return ve
}

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
