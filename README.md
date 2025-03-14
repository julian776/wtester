# WTester | Logtester | Log Tester

`wtester` is a Go package designed for testing log outputs or any other byte stream. It allows you to define expectations on the output and validate whether those expectations are met. This package focuses on governance and compliance with the structure of the logs, ensuring that they adhere to predefined formats and standards.

It is not intended for testing business logic or application logic. It is designed to be used in tests, not in production code. The validations are an overhead that is not needed in production.

## Installation

To install the package, run:

```sh
go get github.com/julian776/wtester
```

## Usage

Check the GoDoc for detailed usage instructions: [GoDoc](https://pkg.go.dev/github.com/julian776/wtester)

## Customization

The package is designed to be flexible and customizable. You can define your own expectations to suit your needs.
Check the [expectations.go](expectations.go) file for some examples.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
