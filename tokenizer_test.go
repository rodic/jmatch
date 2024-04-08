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
				{tokenType: LeftBrace, value: "{"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "emptyWithSpaces",
			input: "{     }",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "simplePair",
			input: "{\"a\":\"1\"}",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: String, value: "1"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "simplePairWithSpaces",
			input: "{  \"a\"  :   \"1\"  }",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: String, value: "1"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "simplePairWithInt",
			input: "{\"a\":1}",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: Number, value: "1"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "simplePairWithNegativeInt",
			input: "{\"a\":-1}",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: Number, value: "-1"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "simplePairWithIntAndSpaces",
			input: "{  \"a\"  : 1  }",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: Number, value: "1"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "simplePairWithDecimal",
			input: "{\"a\":1.01}",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: Number, value: "1.01"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "simplePairWithDecimalAndSpaces",
			input: "{  \"a\"  :  1.01  }",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: Number, value: "1.01"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "simplePairWithTrue",
			input: "{\"a\":true}",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: Boolean, value: "true"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "simplePairWithTrueAndSpaces",
			input: "{  \"a\"  :  true  }",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: Boolean, value: "true"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "simplePairWithFalse",
			input: "{\"a\":false}",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: Boolean, value: "false"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "simplePairWithFalseAndSpaces",
			input: "{  \"a\"  :  false  }",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: Boolean, value: "false"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "simplePairWithNull",
			input: "{\"a\":null}",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: Null, value: "null"},
				{tokenType: RightBrace, value: "}"}}},
		{name: "simplePairWithNullAndSpaces",
			input: "{  \"a\"  :  null  }",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: Null, value: "null"},
				{tokenType: RightBrace, value: "}"}}},

		// Array
		{name: "simplePairWithArray",
			input: "{\"a\":[1, \"2\", true, null]}",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: LeftBracket, value: "["},
				{tokenType: Number, value: "1"},
				{tokenType: Comma, value: ","},
				{tokenType: String, value: "2"},
				{tokenType: Comma, value: ","},
				{tokenType: Boolean, value: "true"},
				{tokenType: Comma, value: ","},
				{tokenType: Null, value: "null"},
				{tokenType: RightBracket, value: "]"},
				{tokenType: RightBrace, value: "}"}}},

		// Nested
		{name: "nested",
			input: "{\"a\": \"ünicode\", \"b\"  : { \"list\": [true, [false, null, {  }]]}}",
			expected: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: String, value: "ünicode"},
				{tokenType: Comma, value: ","},
				{tokenType: String, value: "b"},
				{tokenType: Colon, value: ":"},
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "list"},
				{tokenType: Colon, value: ":"},
				{tokenType: LeftBracket, value: "["},
				{tokenType: Boolean, value: "true"},
				{tokenType: Comma, value: ","},
				{tokenType: LeftBracket, value: "["},
				{tokenType: Boolean, value: "false"},
				{tokenType: Comma, value: ","},
				{tokenType: Null, value: "null"},
				{tokenType: Comma, value: ","},
				{tokenType: LeftBrace, value: "{"},
				{tokenType: RightBrace, value: "}"},
				{tokenType: RightBracket, value: "]"},
				{tokenType: RightBracket, value: "]"},
				{tokenType: RightBrace, value: "}"},
				{tokenType: RightBrace, value: "}"},
			}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenizer := NewTokenizer(tc.input)
			result, err := tokenizer.Tokenize()
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
			expected: "unexpected token . in JSON at position 9"},
		{name: "invalidText",
			input:    "{\"a\":    truef}",
			expected: "unexpected token truef in JSON at position 10"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenizer := NewTokenizer(tc.input)
			res, err := tokenizer.Tokenize()
			if err == nil {
				t.Errorf("Expected error %s got result %v", tc.expected, res)
			}
			if err.Error() != tc.expected {
				t.Errorf("Expected error %s got %s", tc.expected, err)
			}

		})
	}
}
