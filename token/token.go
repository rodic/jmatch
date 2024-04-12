package token

import c "github.com/rodic/jmatch/common"

type TokenType int

type Token struct {
	tokenType TokenType
	Value     string
	line      int
	column    int
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

func (t Token) AsUnexpectedTokenErr() c.UnexpectedTokenErr {
	return c.UnexpectedTokenErr{Token: t.Value, Line: t.line, Column: t.column}
}

func NewToken(tokenType TokenType, value string, line int, column int) Token {
	return Token{
		tokenType: tokenType,
		Value:     value,
		line:      line,
		column:    column,
	}
}

func NewStringToken(value string, line int, column int) Token {
	return Token{tokenType: String, Value: value, line: line, column: column}
}

func NewNumberToken(value string, line int, column int) Token {
	return NewToken(Number, value, line, column)
}

func NewBooleanToken(value string, line int, column int) Token {
	return NewToken(Boolean, value, line, column)
}

func NewNullToken(line int, column int) Token {
	return NewToken(Null, "null", line, column)
}

func NewLeftBraceToken(line int, column int) Token {
	return Token{tokenType: LeftBrace, Value: "{", line: line, column: column}
}

func NewRightBraceToken(line int, column int) Token {
	return Token{tokenType: RightBrace, Value: "}", line: line, column: column}
}

func NewLeftBracketToken(line int, column int) Token {
	return Token{tokenType: LeftBracket, Value: "[", line: line, column: column}
}

func NewRightBracketToken(line int, column int) Token {
	return Token{tokenType: RightBracket, Value: "]", line: line, column: column}
}

func NewColonToken(line int, column int) Token {
	return Token{tokenType: Colon, Value: ":", line: line, column: column}
}

func NewCommaToken(line int, column int) Token {
	return Token{tokenType: Comma, Value: ",", line: line, column: column}
}
