# jmatch

The `jmatch` library is designed to parse and process JSON data in Go. 

It provides a `Match` function that can be invoked with an instance of `io.Reader` containing JSON data and a match function that implements the `Matcher` interface.

```go
type Matcher interface {
	Match(path string, token t.Token)
}
```

params for match functions are:
- `path` is a string that is compatible with `jq` syntax.
- `token` is a struct in Go that contains the value and type information.

```go
type Token struct {
	_type  tokenType
	Value  string
}
```

`type` of the token is a `json` type and can be queried using the following methods:
- `token.IsString()`
- `token.IsNumber()`
- `token.IsBoolean()`
- `token.IsNull()`

## Usage example:

```go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rodic/jmatch"
)

type FixedTokenValueMatch struct {
	matchingString string
	matches        []string
}

func (fm *FixedTokenValueMatch) Match(path string, token jmatch.Token) {
	if token.IsNumber() && token.Value == fm.matchingString {
		fm.matches = append(fm.matches, path)
	}
}

func main() {

	fm := FixedTokenValueMatch{
		matchingString: "2",
		matches:        make([]string, 0, 8),
	}

	jsonReader := strings.NewReader("{\"a\": {\"b.c\": [\"2\", 2]}}")

	err := jmatch.Match(jsonReader, &fm)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// [.a."b.c"[1]] only snd elem matches since the first is string
	fmt.Printf("%v\n", fm.matches)
}
```

## TODO

- improve integration tests with invalid inputs
- cmd tool

## License

This project is licensed under the MIT license. Please see the [LICENSE](LICENSE) file for more details.
