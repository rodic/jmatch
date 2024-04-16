package matcher

import t "github.com/rodic/jmatch/tokenizer"

type Matcher interface {
	Match(path string, token t.Token)
}
