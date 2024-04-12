package token

type TokenList struct {
	tokens      []Token
	tokensCount int
	currentId   int
	nextId      int
}

func (t *TokenList) Current() Token {
	return t.tokens[t.currentId]
}

func (t *TokenList) Next() Token {
	return t.tokens[t.currentId+1]
}

func (t *TokenList) Move() {
	t.currentId++
}

func (t *TokenList) Empty() bool {
	return t.tokensCount == 0
}

func (t *TokenList) HasNext() bool {
	return t.currentId < t.tokensCount-1
}

func NewTokens(tokens []Token) TokenList {
	return TokenList{
		tokens:      tokens,
		tokensCount: len(tokens),
		currentId:   0,
		nextId:      1,
	}
}
