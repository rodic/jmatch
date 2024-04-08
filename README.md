# jmatch

The `jmatch` library is designed to parse and process JSON data in Go.

To parse JSON using `jmatch`, you would typically follow these steps:

1. Import the `jmatch` package in your Go code:
    ```go
    import "github.com/arodic/jmatch"
    ```

2. Create a matcher object that implements `jmatch.Matcher` interface:
    ```go
    type Matcher interface {
	    match(path string, token Token)
    }
    ```

    `jmatch` parser will produce `path` -> `token` mappings, where `path` is a string in from `.key.[arrayIndex]` and `token` is a struct with value (stored as string) and type (string, number, bool, null)
    ```go
    type Token struct {
    	tokenType TokenType
    	value     string
    }
    ```
    the matcher will be called with both of them.

    ```go
    type FixedTokenValueMatch struct {
    	matchingString string
	    matches        []string
    }

    func (fm *FixedTokenValueMatch) match(path string, token Token) {
	    if token.value == fm.matchingString {
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

    jmatch.Match(tc.json, &fm)

    fmt.Printf("%v", fm.matches) // {'.a.b.[1]'}

   ```


## todo

- handle json with more complex keys {"a.b": 1}
- improve error reporting
- cmd tool
- fixed match
- excluding bools from matched text
- regular expressions
- numbers
- streaming