package tokenizer

import c "github.com/rodic/jmatch/common"

type tokenType int

const (
	leftBrace tokenType = iota
	rightBrace
	leftBracket
	rightBracket
	comma
	str
	number
	boolean
	null
	colon
)

type Token struct {
	_type  tokenType
	Value  string
	Line   int
	Column int
}

func new(t tokenType, value string, line int, column int) Token {
	return Token{
		_type:  t,
		Value:  value,
		Line:   line,
		Column: column,
	}
}

func NewStringToken(value string, line int, column int) Token {
	return Token{_type: str, Value: value, Line: line, Column: column}
}

func NewNumberToken(value string, line int, column int) Token {
	return new(number, value, line, column)
}

func NewBooleanToken(value string, line int, column int) Token {
	return new(boolean, value, line, column)
}

func NewNullToken(line int, column int) Token {
	return new(null, "null", line, column)
}

func NewLeftBraceToken(line int, column int) Token {
	return Token{_type: leftBrace, Value: "{", Line: line, Column: column}
}

func NewRightBraceToken(line int, column int) Token {
	return Token{_type: rightBrace, Value: "}", Line: line, Column: column}
}

func NewLeftBracketToken(line int, column int) Token {
	return Token{_type: leftBracket, Value: "[", Line: line, Column: column}
}

func NewRightBracketToken(line int, column int) Token {
	return Token{_type: rightBracket, Value: "]", Line: line, Column: column}
}

func NewColonToken(line int, column int) Token {
	return Token{_type: colon, Value: ":", Line: line, Column: column}
}

func NewCommaToken(line int, column int) Token {
	return Token{_type: comma, Value: ",", Line: line, Column: column}
}

func (t Token) IsLeftBrace() bool {
	return t._type == leftBrace
}

func (t Token) IsRightBrace() bool {
	return t._type == rightBrace
}

func (t Token) IsLeftBracket() bool {
	return t._type == leftBracket
}

func (t Token) IsRightBracket() bool {
	return t._type == rightBracket
}

func (t Token) IsComma() bool {
	return t._type == comma
}

func (t Token) IsString() bool {
	return t._type == str
}

func (t Token) IsNumber() bool {
	return t._type == number
}

func (t Token) IsBoolean() bool {
	return t._type == boolean
}

func (t Token) IsNull() bool {
	return t._type == null
}

func (t Token) IsColon() bool {
	return t._type == colon
}

func (t Token) AsUnexpectedTokenErr() c.UnexpectedTokenErr {
	return c.UnexpectedTokenErr{Token: t.Value, Line: t.Line, Column: t.Column}
}
