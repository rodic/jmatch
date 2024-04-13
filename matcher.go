package jmatch

import (
	m "github.com/rodic/jmatch/matcher"
	p "github.com/rodic/jmatch/parser"
	t "github.com/rodic/jmatch/token"
	z "github.com/rodic/jmatch/tokenizer"
)

type Token = t.Token

func Match(json string, m m.Matcher) error {

	tokenStream := make(chan t.Token)

	tokenizer := z.NewTokenizer(json, tokenStream)
	go tokenizer.Tokenize()

	parser := p.NewParser(tokenStream, m)
	err := parser.Parse()

	if err != nil {
		return err
	}

	return nil
}
