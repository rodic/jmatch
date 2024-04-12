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

	err := jmatch.Match(json, &fm)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%v\n", fm.matches) // ['.a.b.[1]']
}

```

## TODO

- handle json with more complex keys {"a.b": 1}
- cmd tool
- streaming

## License

This project is licensed under the MIT license. Please see the [LICENSE](LICENSE) file for more details.
