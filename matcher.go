package jmatch

type Matcher interface {
	match(path string, token Token)
}

func Match(json string, m Matcher) (Matcher, error) {
	tokenizer := NewTokenizer(json)
	tokens, err := tokenizer.Tokenize()

	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	result, err := parser.parse()

	if err != nil {
		return nil, err
	}

	for path, token := range result {
		m.match(path, token)
	}

	return m, nil
}
