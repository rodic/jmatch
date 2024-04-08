# jmatch

The `jmatch` library is designed to parse and process JSON data in Go.

To parse JSON using `jmatch`, you would typically follow these steps:


1. Import the `jmatch` package in your Go code:
    ```go
    import "github.com/rodic/jmatch"
    ```

2. Create a matcher object that implements `jmatch.Matcher` interface:
    ```go
    package jmatch

    type Matcher interface {
        Match(path string, token Token)
    }
    ```

    `jmatch` parser will produce (`path`, `token`) pairs, where `path` is a string in from `.key.[arrayIndex]` and `token` is a struct with value (stored as string) and type (string, number, bool, null)
    ```go
    package jmatch

    type Token struct {
        tokenType TokenType
        Value     string
    }
    ```
    the matcher will be called with `path` and `token`.

    ```go
    type FixedTokenValueMatch struct {
        matchingString string
        matches        []string
    }

    func (fm *FixedTokenValueMatch) match(path string, token Token) {
	    if token.Value == fm.matchingString {
            fm.matches = append(fm.matches, path)
        }
    }

    fm := FixedTokenValueMatch{
        matchingString: "2",
        matches: make([]string, 0, 8)
    }
    ```

3. Call Match fn against input JSON and matcher:
    ```go
    json := "{\"a\": {\"b\": [1, 2]}}"

    jmatch.Match(json, &fm)

    fmt.Printf("%v\n", fm.matches) // {'.a.b.[1]'}
    ```

Complete script
```go
package main

import (
	"fmt"

	"github.com/rodic/jmatch"
)

type FixedTokenValueMatch struct {
	matchingString string
	matches        []string
}

func (fm *FixedTokenValueMatch) Match(path string, token jmatch.Token) {
	if token.Value == fm.matchingString {
		fm.matches = append(fm.matches, path)
	}
}

func main() {
	fm := FixedTokenValueMatch{
		matchingString: "2",
		matches:        make([]string, 0, 8),
	}

	json := "{\"a\": {\"b\": [1, 2]}}"

	jmatch.Match(json, &fm)

	fmt.Printf("%v\n", fm.matches) // [.a.b.[1]]
}
```

## TODO

- improve text match to allow other characters beside letters
- handle json with more complex keys {"a.b": 1}
- improve error reporting
- cmd tool
- fixed match
- excluding bools from matched text
- regular expressions
- numbers
- streaming

## License

This project is licensed under the MIT license. Please see the [LICENSE](LICENSE) file for more details.
