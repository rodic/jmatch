package jmatch

import (
	m "github.com/rodic/jmatch/matcher"
	p "github.com/rodic/jmatch/parser"
	t "github.com/rodic/jmatch/token"
	z "github.com/rodic/jmatch/tokenizer"
)

type Token = t.Token

func Match(json string, matcher m.Matcher) error {

	tokenizer := z.NewTokenizer(json)
	go tokenizer.Tokenize()

	parser := p.NewParser(tokenizer.GetTokenReadStream(), matcher)
	err := parser.Parse()

	if err != nil {
		return err
	}

	return nil
}
