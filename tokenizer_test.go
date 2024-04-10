package jmatch

import (
	"reflect"
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
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: RightBrace, Value: "}", line: 1, column: 2}}},
		{name: "emptyWithSpaces",
			input: "{     }",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: RightBrace, Value: "}", line: 1, column: 7}}},
		{name: "simplePair",
			input: "{\"a\":\"1\"}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 5},
				{tokenType: String, Value: "1", line: 1, column: 6},
				{tokenType: RightBrace, Value: "}", line: 1, column: 9}}},
		{name: "simplePairWithSpaces",
			input: "{  \"a\"  :   \"1\"  }",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 4},
				{tokenType: Colon, Value: ":", line: 1, column: 9},
				{tokenType: String, Value: "1", line: 1, column: 13},
				{tokenType: RightBrace, Value: "}", line: 1, column: 18}}},
		{name: "simplePairWithInt",
			input: "{\"a\":1}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 5},
				{tokenType: Number, Value: "1", line: 1, column: 6},
				{tokenType: RightBrace, Value: "}", line: 1, column: 7}}},
		{name: "simplePairWithNegativeInt",
			input: "{\"a\":-1}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 5},
				{tokenType: Number, Value: "-1", line: 1, column: 6},
				{tokenType: RightBrace, Value: "}", line: 1, column: 8}}},
		{name: "simplePairWithIntAndSpaces",
			input: "{  \"a\"  : 1  }",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 4},
				{tokenType: Colon, Value: ":", line: 1, column: 9},
				{tokenType: Number, Value: "1", line: 1, column: 11},
				{tokenType: RightBrace, Value: "}", line: 1, column: 14}}},
		{name: "simplePairWithDecimal",
			input: "{\"a\":1.01}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 5},
				{tokenType: Number, Value: "1.01", line: 1, column: 6},
				{tokenType: RightBrace, Value: "}", line: 1, column: 10}}},
		{name: "simplePairWithTrue",
			input: "{\"a\":true}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 5},
				{tokenType: Boolean, Value: "true", line: 1, column: 6},
				{tokenType: RightBrace, Value: "}", line: 1, column: 10}}},
		{name: "simplePairWithFalse",
			input: "{\"a\":false}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 5},
				{tokenType: Boolean, Value: "false", line: 1, column: 6},
				{tokenType: RightBrace, Value: "}", line: 1, column: 11}}},
		{name: "simplePairWithNull",
			input: "{\"a\":null}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 5},
				{tokenType: Null, Value: "null", line: 1, column: 6},
				{tokenType: RightBrace, Value: "}", line: 1, column: 10}}},
		{name: "simplePairWithDotInKey",
			input: "{\"a.b\":\"1\"}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a.b", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 7},
				{tokenType: String, Value: "1", line: 1, column: 8},
				{tokenType: RightBrace, Value: "}", line: 1, column: 11}}},
		{name: "simplePairWithSpaceInKey",
			input: "{\"a b\":\"1\"}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a b", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 7},
				{tokenType: String, Value: "1", line: 1, column: 8},
				{tokenType: RightBrace, Value: "}", line: 1, column: 11}}},

		// Array
		{name: "simplePairWithArray",
			input: "{\"a\":[1, \"2\", true, null]}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 5},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 6},
				{tokenType: Number, Value: "1", line: 1, column: 7},
				{tokenType: Comma, Value: ",", line: 1, column: 8},
				{tokenType: String, Value: "2", line: 1, column: 10},
				{tokenType: Comma, Value: ",", line: 1, column: 13},
				{tokenType: Boolean, Value: "true", line: 1, column: 15},
				{tokenType: Comma, Value: ",", line: 1, column: 19},
				{tokenType: Null, Value: "null", line: 1, column: 21},
				{tokenType: RightBracket, Value: "]", line: 1, column: 25},
				{tokenType: RightBrace, Value: "}", line: 1, column: 26}}},

		// Nested
		{name: "nested",
			input: "{\"a\": \"ünicode\", \"b\"  : { \"list\": [true, [false, null, {  }]]}}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 5},
				{tokenType: String, Value: "ünicode", line: 1, column: 7},

				{tokenType: Comma, Value: ",", line: 1, column: 16},
				{tokenType: String, Value: "b", line: 1, column: 18},
				{tokenType: Colon, Value: ":", line: 1, column: 23},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 25},
				{tokenType: String, Value: "list", line: 1, column: 27},

				{tokenType: Colon, Value: ":", line: 1, column: 33},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 35},
				{tokenType: Boolean, Value: "true", line: 1, column: 36},
				{tokenType: Comma, Value: ",", line: 1, column: 40},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 42},
				{tokenType: Boolean, Value: "false", line: 1, column: 43},
				{tokenType: Comma, Value: ",", line: 1, column: 48},
				{tokenType: Null, Value: "null", line: 1, column: 50},
				{tokenType: Comma, Value: ",", line: 1, column: 54},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 56},
				{tokenType: RightBrace, Value: "}", line: 1, column: 59},
				{tokenType: RightBracket, Value: "]", line: 1, column: 60},
				{tokenType: RightBracket, Value: "]", line: 1, column: 61},
				{tokenType: RightBrace, Value: "}", line: 1, column: 62},
				{tokenType: RightBrace, Value: "}", line: 1, column: 63},
			}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenizer := newTokenizer(tc.input)
			result, err := tokenizer.tokenize()
			if err != nil {
				t.Error(err)
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
		{name: "invalidNumber",
			input:    "{\"a\":1.2.3}",
			expected: "invalid JSON. unexpected token . at line 1 column 9"},
		{name: "invalidText",
			input:    "{\"a\":    truef}",
			expected: "invalid JSON. unexpected token truef at line 1 column 10"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenizer := newTokenizer(tc.input)
			res, err := tokenizer.tokenize()
			if err == nil {
				t.Errorf("Expected error %s got result %v", tc.expected, res)
			}
			if err.Error() != tc.expected {
				t.Errorf("Expected error %s got %s", tc.expected, err)
			}

		})
	}
}
