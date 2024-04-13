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
	runePositionCounter textPositionCounter
	tokensChan          chan<- token.Token
}

func NewTokenizer(jInput string, tokensChan chan<- token.Token) tokenizer {
	runes := []rune(jInput)
	return tokenizer{
		input:               runes,
		inputLen:            len(runes),
		runePosition:        0,
		runePositionCounter: newTextPositionCounter(),
		tokensChan:          tokensChan,
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

func (t *tokenizer) Tokenize() error {

	for !t.done() {
		current := t.current()

		line := t.runePositionCounter.line
		column := t.runePositionCounter.column

		switch current {
		case '{':
			t.tokensChan <- token.NewLeftBraceToken(line, column)
		case '}':
			t.tokensChan <- token.NewRightBraceToken(line, column)
		case '[':
			t.tokensChan <- token.NewLeftBracketToken(line, column)
		case ']':
			t.tokensChan <- token.NewRightBracketToken(line, column)
		case ',':
			t.tokensChan <- token.NewCommaToken(line, column)
		case '"':
			str := t.getString()
			t.tokensChan <- token.NewStringToken(str, line, column)
		case ':':
			t.tokensChan <- token.NewColonToken(line, column)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			digit := t.getNumber()
			t.tokensChan <- token.NewNumberToken(digit, line, column)
		case ' ':
		case '\n':
			continue
		default:
			text := t.getText()

			if text == "true" || text == "false" {
				t.tokensChan <- token.NewBooleanToken(text, line, column)
			} else if text == "null" {
				t.tokensChan <- token.NewNullToken(line, column)
			} else if text != "" {
				return c.UnexpectedTokenErr{Token: text, Line: line, Column: column}
			} else {
				return c.UnexpectedTokenErr{Token: string(current), Line: line, Column: column}
			}
		}
	}

	close(t.tokensChan)

	return nil
}
