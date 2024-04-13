package parser

import t "github.com/rodic/jmatch/token"

type tokenList struct {
	tokensChan <-chan t.Token
	current    t.Token
	next       t.Token
	hasNext    bool
}

func (t *tokenList) move() {
	t.current = t.next
	t.next, t.hasNext = <-t.tokensChan
}

func NewTokens(tokensChan <-chan t.Token) tokenList {
	current, hasNext := <-tokensChan

	var next t.Token

	if hasNext {
		next, hasNext = <-tokensChan
	}

	return tokenList{
		tokensChan: tokensChan,
		current:    current,
		next:       next,
		hasNext:    hasNext,
	}
}
