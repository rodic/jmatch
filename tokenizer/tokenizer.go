package tokenizer

import (
	"fmt"
	"io"
	"unicode"

	c "github.com/rodic/jmatch/common"
)

type TokenResult struct {
	Token Token
	Error error
}

type tokenizer struct {
	runes       RuneReader
	tokenStream chan TokenResult
}

func NewTokenizer(r io.Reader) tokenizer {
	return tokenizer{
		runes:       NewRuneReader(r),
		tokenStream: make(chan TokenResult),
	}
}

func (t *tokenizer) GetTokenReadStream() <-chan TokenResult {
	return t.tokenStream
}

func (t *tokenizer) getString() (string, error) {
	res := []rune{}

	if err := t.runes.move(); err != nil {
		return "", err
	}

	for {
		if t.runes.current == '"' {
			break
		}
		res = append(res, t.runes.current)

		if err := t.runes.move(); err != nil {
			return "", err
		}
	}
	return string(res), nil
}

func (t *tokenizer) getNumber() (string, error) {
	res := []rune{t.runes.current}

	dotCount := 0
	isDigitSet := t.runes.current != '-'

	for {
		if err := t.runes.move(); err != nil {
			return "", err
		}

		if unicode.IsDigit(t.runes.current) {
			isDigitSet = true
			res = append(res, t.runes.current)
		} else if t.runes.current == '.' && dotCount == 0 && isDigitSet {
			dotCount++
			res = append(res, t.runes.current)
		} else {
			t.runes.rewind()
			break
		}
	}

	// if unexpected end of number
	if !unicode.IsDigit(t.runes.current) {
		t.runes.move()
		return string(t.runes.current), fmt.Errorf("invalid number")
	}

	return string(res), nil
}

func (t *tokenizer) getText() (string, error) {
	res := []rune{t.runes.current}

	for {
		if err := t.runes.move(); err != nil {
			return "", err
		}

		if unicode.IsLetter(t.runes.current) {
			res = append(res, t.runes.current)
		} else {
			t.runes.rewind()
			break
		}
	}
	return string(res), nil
}

func (t *tokenizer) writeTokenResult(token Token) {
	t.tokenStream <- TokenResult{Token: token, Error: nil}
}

func (t *tokenizer) writeError(err error) {
	t.tokenStream <- TokenResult{Error: err}
}

func (t *tokenizer) Tokenize() {

	defer close(t.tokenStream)

	for {
		t.runes.move()

		if t.runes.done {
			break
		}

		current := t.runes.current

		line := t.runes.line
		column := t.runes.column

		switch current {
		case ' ':
		case '\n':
			continue
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
		case ':':
			t.writeTokenResult(NewColonToken(line, column))
		case '"':
			str, err := t.getString()
			if err == nil {
				t.writeTokenResult(NewStringToken(str, line, column))
			} else {
				t.writeError(err)
				return
			}
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			digit, err := t.getNumber()
			if err == nil {
				t.writeTokenResult(NewNumberToken(digit, line, column))
			} else {
				t.writeError(c.UnexpectedTokenErr{Token: digit, Line: t.runes.line, Column: t.runes.column})
				return
			}
		default:
			text, err := t.getText()

			if err != nil {
				t.writeError(err)
				return
			}

			if text == "true" || text == "false" {
				t.writeTokenResult(NewBooleanToken(text, line, column))
			} else if text == "null" {
				t.writeTokenResult(NewNullToken(line, column))
			} else {
				t.writeError(c.UnexpectedTokenErr{Token: text, Line: line, Column: column})
				return
			}
		}
	}
}
