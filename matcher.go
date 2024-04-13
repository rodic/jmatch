package jmatch

import (
	"fmt"

	m "github.com/rodic/jmatch/matcher"
	p "github.com/rodic/jmatch/parser"
	t "github.com/rodic/jmatch/token"
	z "github.com/rodic/jmatch/tokenizer"
)

type Token = t.Token

func Match(json string, m m.Matcher) error {

	tokensChan := make(chan t.Token, 256)

	tokenizer := z.NewTokenizer(json, tokensChan)
	go tokenizer.Tokenize()

	tokens := make([]t.Token, 0, 32)

	for t := range tokensChan {
		fmt.Printf(">>>>>>>>>>> %v\n", t)
		tokens = append(tokens, t)
	}

	parser := p.NewParser(tokens, m)
	err := parser.Parse()

	if err != nil {
		return err
	}

	return nil
}
