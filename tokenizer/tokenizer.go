package tokenizer

import (
	"unicode"

	c "github.com/rodic/jmatch/common"
)

type TokenResult struct {
	Token Token
	Error error
}

type tokenizer struct {
	input               []rune
	inputLen            int
	runePosition        int
	textPositionCounter textPositionCounter
	tokenStream         chan TokenResult
}

func NewTokenizer(jInput string) tokenizer {
	runes := []rune(jInput)
	return tokenizer{
		input:               runes,
		inputLen:            len(runes),
		runePosition:        0,
		textPositionCounter: newTextPositionCounter(),
		tokenStream:         make(chan TokenResult),
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

func (t *tokenizer) GetTokenReadStream() <-chan TokenResult {
	return t.tokenStream
}

func (t *tokenizer) writeTokenResult(token Token) {
	t.tokenStream <- TokenResult{Token: token, Error: nil}
}

func (t *tokenizer) writeError(err error) {
	t.tokenStream <- TokenResult{Error: err}
}

func (t *tokenizer) Tokenize() {

	defer close(t.tokenStream)

	for !t.done() {
		current := t.current()

		line := t.textPositionCounter.line
		column := t.textPositionCounter.column

		switch current {
		case '{':
			t.writeTokenResult(NewLeftBraceToken(line, column))
		case '}':
			t.writeTokenResult(NewRightBraceToken(line, column))
		case '[':
			t.writeTokenResult(NewLeftBracketToken(line, column))
		case ']':
			t.writeTokenResult(NewRightBracketToken(line, column))
		case ',':
			t.writeTokenResult(NewCommaToken(line, column))
		case '"':
			str := t.getString()
			t.writeTokenResult(NewStringToken(str, line, column))
		case ':':
			t.writeTokenResult(NewColonToken(line, column))
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			digit := t.getNumber()
			t.writeTokenResult(NewNumberToken(digit, line, column))
		case ' ':
		case '\n':
			continue
		default:
			text := t.getText()

			if text == "true" || text == "false" {
				t.writeTokenResult(NewBooleanToken(text, line, column))
			} else if text == "null" {
				t.writeTokenResult(NewNullToken(line, column))
			} else if text != "" {
				t.writeError(c.UnexpectedTokenErr{Token: text, Line: line, Column: column})
				return
			} else {
				t.writeError(c.UnexpectedTokenErr{Token: string(current), Line: line, Column: column})
				return
			}
		}
	}
}
