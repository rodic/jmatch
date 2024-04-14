package parser

import (
	t "github.com/rodic/jmatch/token"
	z "github.com/rodic/jmatch/tokenizer"
)

type tokenList struct {
	tokensChan <-chan z.TokenResult
	current    t.Token
	next       t.Token
	hasNext    bool
}

func (t *tokenList) move() error {
	var nextResult z.TokenResult

	nextResult, t.hasNext = <-t.tokensChan

	if nextResult.Error != nil {
		return nextResult.Error
	}

	t.current = t.next
	t.next = nextResult.Token

	return nil
}

func NewTokens(tokensChan <-chan z.TokenResult) (*tokenList, error) {
	currentResult, hasNext := <-tokensChan

	if currentResult.Error != nil {
		return nil, currentResult.Error
	}

	var nextResult z.TokenResult

	if hasNext {
		nextResult, hasNext = <-tokensChan

		if nextResult.Error != nil {
			return nil, nextResult.Error
		}
	}

	tl := tokenList{
		tokensChan: tokensChan,
		current:    currentResult.Token,
		next:       nextResult.Token,
		hasNext:    hasNext,
	}

	return &tl, nil
}
