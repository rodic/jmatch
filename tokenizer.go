package jmatch

import (
	"fmt"
	"unicode"
)

type Tokenizer struct {
	input    []rune
	inputLen int
	position int
}

func NewTokenizer(jInput string) Tokenizer {
	runes := []rune(jInput)
	return Tokenizer{
		input:    runes,
		inputLen: len(runes),
		position: 0,
	}
}

func (t *Tokenizer) done() bool {
	return t.position >= t.inputLen
}

func (t *Tokenizer) current() rune {
	r := t.input[t.position]
	t.move()
	return r
}

func (t *Tokenizer) getString() string {
	res := []rune{}

	for {
		c := t.current()
		if c == '"' {
			break
		}
		res = append(res, c)
	}
	return string(res)
}

func (t *Tokenizer) getNumber() string {
	res := []rune{}
	dotCount := 0 // one dot in num is allowed
	isFirst := true

	t.rewind()

	for {
		c := t.current()

		if c == '.' {
			dotCount += 1
		}

		isMinus := c == '-' && isFirst
		isDigit := unicode.IsDigit(c)
		isDot := c == '.' && dotCount <= 1 && !isFirst

		isFirst = false

		if isMinus || isDigit || isDot {
			res = append(res, c)
		} else {
			t.rewind()
			break
		}
	}
	return string(res)
}

func (t *Tokenizer) getText() string {
	res := []rune{}
	t.rewind()

	for {
		c := t.current()

		if unicode.IsLetter(c) {
			res = append(res, c)
		} else {
			t.rewind()
			break
		}
	}
	return string(res)
}

func (t *Tokenizer) move() {
	t.position++
}

func (t *Tokenizer) rewind() {
	t.position--
}

func (t *Tokenizer) Tokenize() ([]Token, error) {

	res := make([]Token, 0, 8)

	for !t.done() {
		switch c := t.current(); c {
		case '{':
			res = append(res, Token{tokenType: LeftBrace, Value: "{"})
		case '}':
			res = append(res, Token{tokenType: RightBrace, Value: "}"})
		case '[':
			res = append(res, Token{tokenType: LeftBracket, Value: "["})
		case ']':
			res = append(res, Token{tokenType: RightBracket, Value: "]"})
		case ',':
			res = append(res, Token{tokenType: Comma, Value: ","})
		case '"':
			str := t.getString()
			res = append(res, Token{tokenType: String, Value: str})
		case ':':
			res = append(res, Token{tokenType: Colon, Value: ":"})
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			digit := t.getNumber()
			res = append(res, Token{tokenType: Number, Value: digit})
		case ' ':
			continue
		default:
			text := t.getText()

			if text == "true" || text == "false" {
				res = append(res, Token{tokenType: Boolean, Value: text})
			} else if text == "null" {
				res = append(res, Token{tokenType: Null, Value: text})
			} else if text != "" {
				errorPos := t.position + 1 - len(text)
				return nil, fmt.Errorf("unexpected token %s in JSON at position %d", text, errorPos)
			} else {
				return nil, fmt.Errorf("unexpected token %c in JSON at position %d", c, t.position+1)
			}
		}
	}

	return res, nil
}
