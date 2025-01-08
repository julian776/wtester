// wtester is designed for testing log outputs
// or any other byte stream.
// It allows you to define expectations on
// the output and validate whether those expectations are met.
// This is particularly useful for testing
// logs and ensuring that the expected log messages are produced.
//
// As an important note, do not get mad trying
// to test exact log messages, but rather test
// the structure of the log messages.
// For example, test that a log message contains
// a certain sub string, has the required fields,
// or has the expected structure.
//
// The package is designed to be used in tests
// not in production code. The validations are
// an overhead that is not needed in production.
package wtester
