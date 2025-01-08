# Error Iterator

[![Go Reference](https://pkg.go.dev/badge/vawter.tech/eiter.svg)](https://pkg.go.dev/vawter.tech/eiter)

```shell
go get vawter.tech/eiter
```

This package contains iterator types that support error reporting. They interoperates with the stdlib `iter` package.

```go
package main

import (
	"errors"
	"fmt"

	"vawter.tech/eiter"
)

// Counter emits the requested number of values and then an error.
func Counter(ct int) eiter.Seq[int] {
	return eiter.Of(func(yield func(int) bool) error {
		for i := range ct {
			if !yield(i) {
				return nil
			}
		}
		return errors.New("Error World!")
	})
}

func main() {
	var err error
	for i := range Counter(3).Unwrap(&err) {
		fmt.Println(i)
	}
	// This is similar to checking sql.Rows.Err()
	if err != nil {
		fmt.Println(err)
	}
}
```