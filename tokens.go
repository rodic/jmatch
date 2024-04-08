package jmatch

type TokenType int

type Token struct {
	tokenType TokenType
	Value     string
}

const (
	LeftBrace TokenType = iota
	RightBrace
	LeftBracket
	RightBracket
	Comma
	String
	Number
	Boolean
	Null
	Colon
)

func (t Token) IsLeftBrace() bool {
	return t.tokenType == LeftBrace
}

func (t Token) IsRightBrace() bool {
	return t.tokenType == RightBrace
}

func (t Token) IsLeftBracket() bool {
	return t.tokenType == LeftBracket
}

func (t Token) IsRightBracket() bool {
	return t.tokenType == RightBracket
}

func (t Token) IsComma() bool {
	return t.tokenType == Comma
}

func (t Token) IsString() bool {
	return t.tokenType == String
}

func (t Token) IsNumber() bool {
	return t.tokenType == Number
}

func (t Token) IsBoolean() bool {
	return t.tokenType == Boolean
}

func (t Token) IsNull() bool {
	return t.tokenType == Null
}

func (t Token) IsColon() bool {
	return t.tokenType == Colon
}

type Tokens struct {
	tokens      []Token
	tokensCount int
	currentId   int
	nextId      int
}

func (t *Tokens) current() Token {
	return t.tokens[t.currentId]
}

func (t *Tokens) next() Token {
	return t.tokens[t.currentId+1]
}

func (t *Tokens) hasNext() bool {
	return t.currentId < t.tokensCount-2
}

func (t *Tokens) move() {
	t.currentId++
}

func NewTokens(tokens []Token) Tokens {
	return Tokens{
		tokens:      tokens,
		tokensCount: len(tokens),
		currentId:   0,
		nextId:      1,
	}
}