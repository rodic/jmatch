package jmatch

import "fmt"

type UnexpectedEndOfInputErr struct{}

func (e UnexpectedEndOfInputErr) Error() string {
	return "invalid JSON. Unexpected end of JSON input"
}

type UnexpectedTokenErr struct {
	token  string
	line   int
	column int
}

func (e UnexpectedTokenErr) Error() string {
	return fmt.Sprintf("invalid JSON. unexpected token %s at line %d column %d", e.token, e.line, e.column)
}
