package wtester

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

const (
	regexDate = "[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}"
)

func TestWTester_Write(t *testing.T) {
	buf := new(bytes.Buffer)
	wt := NewWTester(buf)

	wt.expects = map[string]*Expect{
		"test": {
			title: "test",
			f: func(p []byte) bool {
				return bytes.Contains(p, []byte("hello"))
			},
		},
	}

	n, err := wt.Write([]byte("hello world"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n != 11 {
		t.Fatalf("expected 11 bytes written, got %d", n)
	}

	if len(wt.errors) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(wt.errors))
	}

	// Assert that the buffer contains the written bytes
	if !bytes.Contains(buf.Bytes(), []byte("hello world")) {
		t.Fatalf("expected buffer to contain 'hello world'")
	}
}

func TestWTester_Close(t *testing.T) {
	// Close writer without implementing io.Closer
	var buf bytes.Buffer
	wt := NewWTester(&buf)

	err := wt.Close()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Close writer implementing io.Closer
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	wt = NewWTester(f)
	err = wt.Close()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = os.Remove(f.Name())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWTester_ExpectWithLogger(t *testing.T) {
	wt := NewWTester(io.Discard)

	wt.Expect("Match hello world", RegexMatch(fmt.Sprintf("^%s hello world\n$", regexDate))).WithMax(1).WithMin(1)
	wt.Expect("error in server", RegexMatch(fmt.Sprintf("^%s error in server\n$", regexDate))).WithMax(1).WithMin(1)
	wt.Expect("everything utf8", ValidUTF8()).Every()

	// Set the logger output to the WTester
	log.SetOutput(wt)

	log.Printf("hello world")
	log.Printf("error in server")

	err := wt.Validate()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWTester_ExpectMax(t *testing.T) {
	buf := new(bytes.Buffer)
	wt := NewWTester(buf)

	wt.Expect("Match str", StringMatch("hi there", false)) // Allow many matches
	wt.Expect("Match bytes", BytesMatch([]byte("julian776"))).WithMax(3)
	wt.Expect("everything utf8", ValidUTF8()).Every()

	wt.Write([]byte("hi there"))
	wt.Write([]byte("hi there"))
	wt.Write([]byte("hi there"))
	wt.Write([]byte("hi there"))
	wt.Write([]byte("hi there"))

	wt.Write([]byte("julian776"))
	wt.Write([]byte("julian776"))
	wt.Write([]byte("julian776"))
	wt.Write([]byte("julian776"))

	err := wt.Validate()
	if err != nil {
		ve, ok := err.(ValidationErrors)
		if !ok {
			t.Fatalf("expected ValidationErrors, got %T", err)
		}

		valErr := ve.Errs[0]
		if valErr.Title != "Match bytes" {
			t.Fatalf("expected title 'Match bytes', got %q", valErr.Title)
		}

		err = valErr.Errors[0]
		if err.Error() != "expected at most 3 matches, got 4\n" {
			t.Fatalf("expected error 'expected at most 3 matches, got 4\n', got %q", err.Error())
		}
	} else {
		t.Fatalf("expected error, got nil")
	}
}
