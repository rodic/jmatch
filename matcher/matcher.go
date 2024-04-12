package matcher

import "github.com/rodic/jmatch/token"

type Matcher interface {
	Match(path string, token token.Token)
}
