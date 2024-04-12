package tokenizer

import (
	"unicode"

	c "github.com/rodic/jmatch/common"
	token "github.com/rodic/jmatch/token"
)

type tokenizer struct {
	input               []rune
	inputLen            int
	runePosition        int
	runePositionCounter positionCounter
}

func NewTokenizer(jInput string) tokenizer {
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

func (t *tokenizer) Tokenize() ([]token.Token, error) {

	res := make([]token.Token, 0, 8)

	for !t.done() {
		current := t.current()

		line := t.runePositionCounter.line
		column := t.runePositionCounter.column

		switch current {
		case '{':
			res = append(res, token.NewLeftBraceToken(line, column))
		case '}':
			res = append(res, token.NewRightBraceToken(line, column))
		case '[':
			res = append(res, token.NewLeftBracketToken(line, column))
		case ']':
			res = append(res, token.NewRightBracketToken(line, column))
		case ',':
			res = append(res, token.NewCommaToken(line, column))
		case '"':
			str := t.getString()
			res = append(res, token.NewStringToken(str, line, column))
		case ':':
			res = append(res, token.NewColonToken(line, column))
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			digit := t.getNumber()
			res = append(res, token.NewNumberToken(digit, line, column))
		case ' ':
		case '\n':
			continue
		default:
			text := t.getText()

			if text == "true" || text == "false" {
				res = append(res, token.NewBooleanToken(text, line, column))
			} else if text == "null" {
				res = append(res, token.NewNullToken(line, column))
			} else if text != "" {
				return nil, c.UnexpectedTokenErr{Token: text, Line: line, Column: column}
			} else {
				return nil, c.UnexpectedTokenErr{Token: string(current), Line: line, Column: column}
			}
		}
	}

	return res, nil
}
