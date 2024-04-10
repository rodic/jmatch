package jmatch

import (
	"fmt"
	"unicode"
)

type tokenizer struct {
	input               []rune
	inputLen            int
	runePosition        int
	runePositionCounter positionCounter
}

func newTokenizer(jInput string) tokenizer {
	runes := []rune(jInput)
	return tokenizer{
		input:               runes,
		inputLen:            len(runes),
		runePosition:        0,
		runePositionCounter: newRunePositionCounter(),
	}
}

func (t *tokenizer) done() bool {
	return t.runePosition >= t.inputLen
}

func (t *tokenizer) current() rune {
	r := t.input[t.runePosition]
	t.move()
	t.runePositionCounter.update(r)
	return r
}

func (t *tokenizer) getString() string {
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

func (t *tokenizer) getNumber() string {
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

func (t *tokenizer) getText() string {
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

func (t *tokenizer) move() {
	t.runePosition++
}

func (t *tokenizer) rewind() {
	t.runePosition--
	t.runePositionCounter.decreaseColumn()
}

func (t *tokenizer) tokenize() ([]Token, error) {

	res := make([]Token, 0, 8)

	for !t.done() {
		c := t.current()

		line := t.runePositionCounter.line
		column := t.runePositionCounter.column

		switch c {
		case '{':
			res = append(res, newToken(LeftBrace, "{", line, column))
		case '}':
			res = append(res, newToken(RightBrace, "}", line, column))
		case '[':
			res = append(res, newToken(LeftBracket, "[", line, column))
		case ']':
			res = append(res, newToken(RightBracket, "]", line, column))
		case ',':
			res = append(res, newToken(Comma, ",", line, column))
		case '"':
			str := t.getString()
			res = append(res, newToken(String, str, line, column))
		case ':':
			res = append(res, newToken(Colon, ":", line, column))
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			digit := t.getNumber()
			res = append(res, newToken(Number, digit, line, column))
		case ' ':
			continue
		default:
			text := t.getText()

			if text == "true" || text == "false" {
				res = append(res, newToken(Boolean, text, line, column))
			} else if text == "null" {
				res = append(res, newToken(Null, text, line, column))
			} else if text != "" {
				return nil, fmt.Errorf(
					"invalid JSON. unexpected token %s at line %d column %d", text, line, column,
				)
			} else {
				return nil, fmt.Errorf(
					"invalid JSON. unexpected token %c at line %d column %d", c, line, column,
				)
			}
		}
	}

	return res, nil
}
