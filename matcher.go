package jmatch

import (
	m "github.com/rodic/jmatch/matcher"
	p "github.com/rodic/jmatch/parser"
	t "github.com/rodic/jmatch/tokenizer"
)

func Match(json string, m m.Matcher) (m.Matcher, error) {
	tokenizer := t.NewTokenizer(json)
	tokens, err := tokenizer.Tokenize()

	if err != nil {
		return nil, err
	}

	parser := p.NewParser(tokens, m)
	err = parser.Parse()

	if err != nil {
		return nil, err
	}

	return m, nil
}
