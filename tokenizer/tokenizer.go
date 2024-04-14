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
	textPositionCounter textPositionCounter
	tokenStream         chan token.Token
}

func NewTokenizer(jInput string) tokenizer {
	runes := []rune(jInput)
	return tokenizer{
		input:               runes,
		inputLen:            len(runes),
		runePosition:        0,
		textPositionCounter: newTextPositionCounter(),
		tokenStream:         make(chan token.Token),
	}
}

func (t *tokenizer) done() bool {
	return t.runePosition >= t.inputLen
}

func (t *tokenizer) current() rune {
	r := t.input[t.runePosition]
	t.move()
	t.textPositionCounter.update(r)
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
	t.textPositionCounter.decreaseColumn()
}

func (t *tokenizer) GetTokenReadStream() <-chan token.Token {
	return t.tokenStream
}

func (t *tokenizer) Tokenize() error {

	defer close(t.tokenStream)

	for !t.done() {
		current := t.current()

		line := t.textPositionCounter.line
		column := t.textPositionCounter.column

		switch current {
		case '{':
			t.tokenStream <- token.NewLeftBraceToken(line, column)
		case '}':
			t.tokenStream <- token.NewRightBraceToken(line, column)
		case '[':
			t.tokenStream <- token.NewLeftBracketToken(line, column)
		case ']':
			t.tokenStream <- token.NewRightBracketToken(line, column)
		case ',':
			t.tokenStream <- token.NewCommaToken(line, column)
		case '"':
			str := t.getString()
			t.tokenStream <- token.NewStringToken(str, line, column)
		case ':':
			t.tokenStream <- token.NewColonToken(line, column)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			digit := t.getNumber()
			t.tokenStream <- token.NewNumberToken(digit, line, column)
		case ' ':
		case '\n':
			continue
		default:
			text := t.getText()

			if text == "true" || text == "false" {
				t.tokenStream <- token.NewBooleanToken(text, line, column)
			} else if text == "null" {
				t.tokenStream <- token.NewNullToken(line, column)
			} else if text != "" {
				return c.UnexpectedTokenErr{Token: text, Line: line, Column: column}
			} else {
				return c.UnexpectedTokenErr{Token: string(current), Line: line, Column: column}
			}
		}
	}

	return nil
}
