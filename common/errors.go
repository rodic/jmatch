package common

import "fmt"

type UnexpectedEndOfInputErr struct{}

func (e UnexpectedEndOfInputErr) Error() string {
	return "invalid JSON. Unexpected end of JSON input"
}

type UnexpectedTokenErr struct {
	Token  string
	Line   int
	Column int
}

func (e UnexpectedTokenErr) Error() string {
	return fmt.Sprintf("invalid JSON. unexpected token %s at line %d column %d", e.Token, e.Line, e.Column)
}
