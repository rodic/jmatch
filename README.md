# jmatch

The `jmatch` library is designed to parse and process JSON data in Go.

```go
package main

import (
	"fmt"
	"os"

	"github.com/rodic/jmatch"
)

type FixedTokenValueMatch struct {
	matchingString string
	matches        []string
}

// implementing Match interface, path is jq compatible string
func (fm *FixedTokenValueMatch) Match(path string, token jmatch.Token) {
	if token.IsNumber() && token.Value == fm.matchingString { // check the token type and value
		fm.matches = append(fm.matches, path)
	}
}

func main() {

	fm := FixedTokenValueMatch{
		matchingString: "2",
		matches:        make([]string, 0, 8),
	}

	json := "{\"a\": {\"b.c\": [\"2\", 2]}}"

	err := jmatch.Match(json, &fm)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%v\n", fm.matches) // [.a."b.c"[1]] only snd elem matches since the first is string
}
```

## TODO

- initialize match with a reader instead of json string
- improve integration tests with invalid inputs
- cmd tool

## License

This project is licensed under the MIT license. Please see the [LICENSE](LICENSE) file for more details.
