# WTester | Logtester | Log Tester

`wtester` is a Go package designed for testing log outputs or any other byte stream. It allows you to define expectations on the output and validate whether those expectations are met. This is particularly useful for testing logs and ensuring that the expected log messages are produced.

## Installation

To install the package, run:

```sh
go get github.com/julian776/wtester
```

## Usage

```go
package main

import (
 "fmt"
 "io"
 "log"

 "github.com/julian776/wtester"
)

func main() {
 wt := wtester.NewWTester(io.Discard)

 wt.Expect("Match hello world", wtester.RegexMatch(`hello world`)).WithMax(1).WithMin(1)
 wt.Expect("Valid UTF-8", wtester.ValidUTF8()).Every()

 log.SetOutput(wt)

 log.Printf("hello world")

 ve := wt.Validate()
 fmt.Println(ve.Errors())
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
