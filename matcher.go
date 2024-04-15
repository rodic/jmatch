package jmatch

import (
	m "github.com/rodic/jmatch/matcher"
	p "github.com/rodic/jmatch/parser"
	t "github.com/rodic/jmatch/tokenizer"
)

type Token = t.Token

// tokenizer -> parser -> matcher
func Match(json string, matcher m.Matcher) error {

	tokenizer := t.NewTokenizer(json)

	go tokenizer.Tokenize()

	parser, err := p.NewParser(tokenizer.GetTokenReadStream())

	if err != nil {
		return err
	}

	go parser.Parse()

	for parsingResult := range parser.GetResultReadStream() {

		if parsingResult.Error != nil {
			return parsingResult.Error
		}

		matcher.Match(parsingResult.Path, parsingResult.Token)
	}

	return nil
}
