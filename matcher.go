package jmatch

type Matcher interface {
	Match(path string, token Token)
}

func Match(json string, m Matcher) (Matcher, error) {
	tokenizer := newTokenizer(json)
	tokens, err := tokenizer.tokenize()

	if err != nil {
		return nil, err
	}

	parser := newParser(tokens, m)
	err = parser.parse()

	if err != nil {
		return nil, err
	}

	return m, nil
}
