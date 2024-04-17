package jmatch

import (
	"io"

	p "github.com/rodic/jmatch/parser"
	t "github.com/rodic/jmatch/tokenizer"
)

type Token = t.Token

type Matcher func(path string, token t.Token)

// tokenizer -> parser -> matcher
func Match(reader io.Reader, matcher Matcher) error {

	tokenizer := t.NewTokenizer(reader)

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

		matcher(parsingResult.Path, parsingResult.Token)
	}

	return nil
}
