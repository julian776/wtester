package wtester

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"testing"
)

const (
	regexDate = "[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}"
)

func Example() {
	wt := NewWTester(io.Discard)

	wt.Expect("Match hello world", RegexMatch(`hello world`)).WithMax(1).WithMin(1)
	wt.Expect("Valid UTF-8", ValidUTF8()).Every()

	log.SetOutput(wt)

	log.Printf("hello world")

	err := wt.Validate()
	if err != nil {
		// No errors should be reported
		fmt.Println("Wt 1:", err)
	}

	wt.Reset()

	wt.Expect("Match server started", StringMatch("server started\n", true)).WithMax(1).WithMin(1)
	wt.Expect("Valid UTF-8", ValidUTF8()).Every()

	log.SetOutput(wt)

	log.Printf("hello world")

	err = wt.Validate()
	if err != nil {
		// Demonstrating type assertion
		ve, ok := err.(*ValidationErrors)
		if !ok {
			fmt.Printf("Error is not of type ValidationError: %T\n", err)
			return
		}

		// One error should be reported
		fmt.Println("Wt 2:", ve.Error())
	}

	// With Slog
	wt.Reset()

	reqId := "56de9ab2-bacd-4e24-a3c7-8221860d8edc"

	wt.Expect("All logs contain req_id", RegexMatch(`"req_id":"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"`)).Every()
	wt.Expect("Must not contain errors", StringMatch("error", false)).WithMax(0)

	logger := slog.New(slog.NewJSONHandler(wt, nil))

	logger.Info("handling parse document", "title", "test", "req_id", reqId)
	logger.Error("failed to parse document", "error", "unexpected EOF", "req_id", reqId)

	err = wt.Validate()
	if err != nil {
		// One error should be reported
		fmt.Println("Wt 3:", err)
	}

	// Output:
	// Wt 2: validation "Match server started"
	// Fails On:
	// expected at least 1 matches, got 0
	// Wt 3: validation "Must not contain errors"
	// Fails On:
	// expected at most 0 matches, got 1
}

func TestWTester_WritesToUnderlyingWriterAreValid(t *testing.T) {
	t.Parallel()

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

func TestWTester_Close_ClosesUnderlyingWriter(t *testing.T) {
	t.Parallel()

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

func TestWTester_ExpectNoErrorsWithLoggerAndValidLogs(t *testing.T) {
	t.Parallel()

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

func TestWTester_ValidateMinMaxExpectationsAreReported(t *testing.T) {
	t.Parallel()

	buf := new(bytes.Buffer)
	wt := NewWTester(buf)

	wt.Expect("Match str", StringMatch("hi there", false)) // Allow many matches
	wt.Expect("Match bytes", BytesMatch([]byte("julian776"))).WithMax(3)
	wt.Expect("Must not match", StringMatch("hello", true)).WithMax(0)
	wt.Expect("Min 4 matches", StringMatch("super dog", false)).WithMin(4)
	wt.Expect("everything utf8", ValidUTF8()).Every()

	wt.Write([]byte("hi there"))
	wt.Write([]byte("hi there"))
	wt.Write([]byte("hi there"))
	wt.Write([]byte("hi there"))
	wt.Write([]byte("hi there"))

	wt.Write([]byte("super dog"))
	wt.Write([]byte("super dog"))
	wt.Write([]byte("super dog"))
	wt.Write([]byte("super dog"))

	wt.Write([]byte("julian776"))
	wt.Write([]byte("julian776"))
	wt.Write([]byte("julian776"))
	wt.Write([]byte("julian776"))

	wt.Write([]byte("hello"))

	err := wt.Validate()
	if err != nil {
		ve, ok := err.(*ValidationErrors)
		if !ok {
			t.Fatalf("expected ValidationErrors, got %T", err)
		}

		matchBytes := false
		mustNotMatch := false
		for _, e := range ve.Errs {
			switch e.Title {
			case "Match bytes":
				matchBytes = true

				err = e.Errors[0]
				if err.Error() != "expected at most 3 matches, got 4\n" {
					t.Fatalf("expected error 'expected at most 3 matches, got 4\n', got %q", err.Error())
				}
			case "Must not match":
				mustNotMatch = true

				err = e.Errors[0]
				if err.Error() != "expected at most 0 matches, got 1\n" {
					t.Fatalf("expected error 'expected at most 0 matches, got 1\n', got %q", err.Error())
				}
			}
		}

		if !matchBytes {
			t.Fatalf("expected error for 'Match bytes'")
		}

		if !mustNotMatch {
			t.Fatalf("expected error for 'Must not match'")
		}
	} else {
		t.Fatalf("expected error, got nil")
	}
}
