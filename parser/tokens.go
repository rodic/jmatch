package parser

import t "github.com/rodic/jmatch/token"

type tokenList struct {
	tokens      []t.Token
	tokensCount int
	currentId   int
	nextId      int
}

func (t *tokenList) current() t.Token {
	return t.tokens[t.currentId]
}

func (t *tokenList) next() t.Token {
	return t.tokens[t.currentId+1]
}

func (t *tokenList) move() {
	t.currentId++
}

func (t *tokenList) empty() bool {
	return t.tokensCount == 0
}

func (t *tokenList) hasNext() bool {
	return t.currentId < t.tokensCount-1
}

func NewTokens(tokens []t.Token) tokenList {
	return tokenList{
		tokens:      tokens,
		tokensCount: len(tokens),
		currentId:   0,
		nextId:      1,
	}
}
