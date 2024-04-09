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
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: RightBrace, Value: "}", line: 1, position: 2}}},
		{name: "emptyWithSpaces",
			input: "{     }",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: RightBrace, Value: "}", line: 1, position: 7}}},
		{name: "simplePair",
			input: "{\"a\":\"1\"}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: String, Value: "a", line: 1, position: 2},
				{tokenType: Colon, Value: ":", line: 1, position: 5},
				{tokenType: String, Value: "1", line: 1, position: 6},
				{tokenType: RightBrace, Value: "}", line: 1, position: 9}}},
		{name: "simplePairWithSpaces",
			input: "{  \"a\"  :   \"1\"  }",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: String, Value: "a", line: 1, position: 4},
				{tokenType: Colon, Value: ":", line: 1, position: 9},
				{tokenType: String, Value: "1", line: 1, position: 13},
				{tokenType: RightBrace, Value: "}", line: 1, position: 18}}},
		{name: "simplePairWithInt",
			input: "{\"a\":1}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: String, Value: "a", line: 1, position: 2},
				{tokenType: Colon, Value: ":", line: 1, position: 5},
				{tokenType: Number, Value: "1", line: 1, position: 6},
				{tokenType: RightBrace, Value: "}", line: 1, position: 7}}},
		{name: "simplePairWithNegativeInt",
			input: "{\"a\":-1}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: String, Value: "a", line: 1, position: 2},
				{tokenType: Colon, Value: ":", line: 1, position: 5},
				{tokenType: Number, Value: "-1", line: 1, position: 6},
				{tokenType: RightBrace, Value: "}", line: 1, position: 8}}},
		{name: "simplePairWithIntAndSpaces",
			input: "{  \"a\"  : 1  }",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: String, Value: "a", line: 1, position: 4},
				{tokenType: Colon, Value: ":", line: 1, position: 9},
				{tokenType: Number, Value: "1", line: 1, position: 11},
				{tokenType: RightBrace, Value: "}", line: 1, position: 14}}},
		{name: "simplePairWithDecimal",
			input: "{\"a\":1.01}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: String, Value: "a", line: 1, position: 2},
				{tokenType: Colon, Value: ":", line: 1, position: 5},
				{tokenType: Number, Value: "1.01", line: 1, position: 6},
				{tokenType: RightBrace, Value: "}", line: 1, position: 10}}},
		{name: "simplePairWithTrue",
			input: "{\"a\":true}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: String, Value: "a", line: 1, position: 2},
				{tokenType: Colon, Value: ":", line: 1, position: 5},
				{tokenType: Boolean, Value: "true", line: 1, position: 6},
				{tokenType: RightBrace, Value: "}", line: 1, position: 10}}},
		{name: "simplePairWithFalse",
			input: "{\"a\":false}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: String, Value: "a", line: 1, position: 2},
				{tokenType: Colon, Value: ":", line: 1, position: 5},
				{tokenType: Boolean, Value: "false", line: 1, position: 6},
				{tokenType: RightBrace, Value: "}", line: 1, position: 11}}},
		{name: "simplePairWithNull",
			input: "{\"a\":null}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: String, Value: "a", line: 1, position: 2},
				{tokenType: Colon, Value: ":", line: 1, position: 5},
				{tokenType: Null, Value: "null", line: 1, position: 6},
				{tokenType: RightBrace, Value: "}", line: 1, position: 10}}},
		{name: "simplePairWithDotInKey",
			input: "{\"a.b\":\"1\"}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: String, Value: "a.b", line: 1, position: 2},
				{tokenType: Colon, Value: ":", line: 1, position: 7},
				{tokenType: String, Value: "1", line: 1, position: 8},
				{tokenType: RightBrace, Value: "}", line: 1, position: 11}}},
		{name: "simplePairWithSpaceInKey",
			input: "{\"a b\":\"1\"}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: String, Value: "a b", line: 1, position: 2},
				{tokenType: Colon, Value: ":", line: 1, position: 7},
				{tokenType: String, Value: "1", line: 1, position: 8},
				{tokenType: RightBrace, Value: "}", line: 1, position: 11}}},

		// Array
		{name: "simplePairWithArray",
			input: "{\"a\":[1, \"2\", true, null]}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: String, Value: "a", line: 1, position: 2},
				{tokenType: Colon, Value: ":", line: 1, position: 5},
				{tokenType: LeftBracket, Value: "[", line: 1, position: 6},
				{tokenType: Number, Value: "1", line: 1, position: 7},
				{tokenType: Comma, Value: ",", line: 1, position: 8},
				{tokenType: String, Value: "2", line: 1, position: 10},
				{tokenType: Comma, Value: ",", line: 1, position: 13},
				{tokenType: Boolean, Value: "true", line: 1, position: 15},
				{tokenType: Comma, Value: ",", line: 1, position: 19},
				{tokenType: Null, Value: "null", line: 1, position: 21},
				{tokenType: RightBracket, Value: "]", line: 1, position: 25},
				{tokenType: RightBrace, Value: "}", line: 1, position: 26}}},

		// Nested
		{name: "nested",
			input: "{\"a\": \"ünicode\", \"b\"  : { \"list\": [true, [false, null, {  }]]}}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, position: 1},
				{tokenType: String, Value: "a", line: 1, position: 2},
				{tokenType: Colon, Value: ":", line: 1, position: 5},
				{tokenType: String, Value: "ünicode", line: 1, position: 7},

				{tokenType: Comma, Value: ",", line: 1, position: 16},
				{tokenType: String, Value: "b", line: 1, position: 18},
				{tokenType: Colon, Value: ":", line: 1, position: 23},
				{tokenType: LeftBrace, Value: "{", line: 1, position: 25},
				{tokenType: String, Value: "list", line: 1, position: 27},

				{tokenType: Colon, Value: ":", line: 1, position: 33},
				{tokenType: LeftBracket, Value: "[", line: 1, position: 35},
				{tokenType: Boolean, Value: "true", line: 1, position: 36},
				{tokenType: Comma, Value: ",", line: 1, position: 40},
				{tokenType: LeftBracket, Value: "[", line: 1, position: 42},
				{tokenType: Boolean, Value: "false", line: 1, position: 43},
				{tokenType: Comma, Value: ",", line: 1, position: 48},
				{tokenType: Null, Value: "null", line: 1, position: 50},
				{tokenType: Comma, Value: ",", line: 1, position: 54},
				{tokenType: LeftBrace, Value: "{", line: 1, position: 56},
				{tokenType: RightBrace, Value: "}", line: 1, position: 59},
				{tokenType: RightBracket, Value: "]", line: 1, position: 60},
				{tokenType: RightBracket, Value: "]", line: 1, position: 61},
				{tokenType: RightBrace, Value: "}", line: 1, position: 62},
				{tokenType: RightBrace, Value: "}", line: 1, position: 63},
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
			expected: "invalid JSON. unexpected token . at line 1 position 9"},
		{name: "invalidText",
			input:    "{\"a\":    truef}",
			expected: "invalid JSON. unexpected token truef at line 1 position 10"},
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
