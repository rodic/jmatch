package tokenizer

import (
	"reflect"
	"strings"
	"testing"
)

func TestTokenizeValidInputs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []Token
	}{
		// Simple cases, {} and one key value pair
		{name: "empty",
			input: "{}",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewRightBraceToken(1, 2)}},
		{name: "emptyWithSpaces",
			input: "{     }",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewRightBraceToken(1, 7)}},
		{name: "simplePair",
			input: "{\"a\":\"1\"}",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewStringToken("a", 1, 2),
				NewColonToken(1, 5),
				NewStringToken("1", 1, 6),
				NewRightBraceToken(1, 9)}},
		{name: "simplePairWithSpaces",
			input: "{  \"a\"  :   \"1\"  }",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewStringToken("a", 1, 4),
				NewColonToken(1, 9),
				NewStringToken("1", 1, 13),
				NewRightBraceToken(1, 18)}},
		{name: "simplePairWithInt",
			input: "{\"a\":1}",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewStringToken("a", 1, 2),
				NewColonToken(1, 5),
				NewNumberToken("1", 1, 6),
				NewRightBraceToken(1, 7)}},
		{name: "simplePairWithNegativeInt",
			input: "{\"a\":-1}",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewStringToken("a", 1, 2),
				NewColonToken(1, 5),
				NewNumberToken("-1", 1, 6),
				NewRightBraceToken(1, 8)}},
		{name: "simplePairWithIntAndSpaces",
			input: "{  \"a\"  : 1  }",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewStringToken("a", 1, 4),
				NewColonToken(1, 9),
				NewNumberToken("1", 1, 11),
				NewRightBraceToken(1, 14)}},
		{name: "simplePairWithDecimal",
			input: "{\"a\":1.01}",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewStringToken("a", 1, 2),
				NewColonToken(1, 5),
				NewNumberToken("1.01", 1, 6),
				NewRightBraceToken(1, 10)}},
		{name: "simplePairWithTrue",
			input: "{\"a\":true}",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewStringToken("a", 1, 2),
				NewColonToken(1, 5),
				NewBooleanToken("true", 1, 6),
				NewRightBraceToken(1, 10)}},
		{name: "simplePairWithFalse",
			input: "{\"a\":false}",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewStringToken("a", 1, 2),
				NewColonToken(1, 5),
				NewBooleanToken("false", 1, 6),
				NewRightBraceToken(1, 11)}},
		{name: "simplePairWithNull",
			input: "{\"a\":null}",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewStringToken("a", 1, 2),
				NewColonToken(1, 5),
				NewNullToken(1, 6),
				NewRightBraceToken(1, 10)}},
		{name: "simplePairWithDotInKey",
			input: "{\"a.b\":\"1\"}",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewStringToken("a.b", 1, 2),
				NewColonToken(1, 7),
				NewStringToken("1", 1, 8),
				NewRightBraceToken(1, 11)}},
		{name: "simplePairWithSpaceInKey",
			input: "{\"a b\":\"1\"}",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewStringToken("a b", 1, 2),
				NewColonToken(1, 7),
				NewStringToken("1", 1, 8),
				NewRightBraceToken(1, 11)}},

		// Array
		{name: "simplePairWithArray",
			input: "{\"a\":[1, \"2\", true, null]}",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewStringToken("a", 1, 2),
				NewColonToken(1, 5),
				NewLeftBracketToken(1, 6),
				NewNumberToken("1", 1, 7),
				NewCommaToken(1, 8),
				NewStringToken("2", 1, 10),
				NewCommaToken(1, 13),
				NewBooleanToken("true", 1, 15),
				NewCommaToken(1, 19),
				NewNullToken(1, 21),
				NewRightBracketToken(1, 25),
				NewRightBraceToken(1, 26)}},

		// Nested
		{name: "nested",
			input: "{\"a\": \"ünicode\", \"b\"  : { \"list\": [true, [false, null, {  }]]}}",
			expected: []Token{
				NewLeftBraceToken(1, 1),
				NewStringToken("a", 1, 2),
				NewColonToken(1, 5),
				NewStringToken("ünicode", 1, 7),

				NewCommaToken(1, 16),
				NewStringToken("b", 1, 18),
				NewColonToken(1, 23),
				NewLeftBraceToken(1, 25),
				NewStringToken("list", 1, 27),

				NewColonToken(1, 33),
				NewLeftBracketToken(1, 35),
				NewBooleanToken("true", 1, 36),
				NewCommaToken(1, 40),
				NewLeftBracketToken(1, 42),
				NewBooleanToken("false", 1, 43),
				NewCommaToken(1, 48),
				NewNullToken(1, 50),
				NewCommaToken(1, 54),
				NewLeftBraceToken(1, 56),
				NewRightBraceToken(1, 59),
				NewRightBracketToken(1, 60),
				NewRightBracketToken(1, 61),
				NewRightBraceToken(1, 62),
				NewRightBraceToken(1, 63),
			}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := make([]Token, 0, 10)
			tokenizer := NewTokenizer(strings.NewReader(tc.input))
			go tokenizer.Tokenize()

			for tokenResult := range tokenizer.GetTokenReadStream() {
				if tokenResult.Error != nil {
					t.Error(tokenResult.Error)
				}
				result = append(result, tokenResult.Token)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected '%v', got '%v' instead\n", tc.expected, result)
			}
		})
	}
}

func TestTokenizeInvalidInputs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "invalidMinus",
			input:    "{\"a\":-}",
			expected: "invalid JSON. unexpected token } at line 1 column 7"},
		{name: "invalidDot",
			input:    "{\"a\":.}",
			expected: "invalid JSON. unexpected token . at line 1 column 6"},
		{name: "invalidNumWithDot",
			input:    "{\"a\":123.}",
			expected: "invalid JSON. unexpected token } at line 1 column 10"},
		{name: "invalidNumber",
			input:    "{\"a\":1.2.3}",
			expected: "invalid JSON. unexpected token . at line 1 column 9"},
		{name: "invalidText",
			input:    "{\"a\":    truef}",
			expected: "invalid JSON. unexpected token truef at line 1 column 10"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenizer := NewTokenizer(strings.NewReader(tc.input))
			go tokenizer.Tokenize()

			var lastTokenResult TokenResult

			for tokenResult := range tokenizer.GetTokenReadStream() {
				lastTokenResult = tokenResult
			}

			if lastTokenResult.Error == nil {
				t.Errorf("Expected error %s but got %v", tc.expected, lastTokenResult.Token)
			} else if lastTokenResult.Error.Error() != tc.expected {
				t.Errorf("Expected error %s got %s", tc.expected, lastTokenResult.Error)
			}
		})
	}
}
