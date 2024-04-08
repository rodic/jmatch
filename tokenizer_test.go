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
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "emptyWithSpaces",
			input: "{     }",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "simplePair",
			input: "{\"a\":\"1\"}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: String, Value: "1"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "simplePairWithSpaces",
			input: "{  \"a\"  :   \"1\"  }",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: String, Value: "1"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "simplePairWithInt",
			input: "{\"a\":1}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: Number, Value: "1"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "simplePairWithNegativeInt",
			input: "{\"a\":-1}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: Number, Value: "-1"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "simplePairWithIntAndSpaces",
			input: "{  \"a\"  : 1  }",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: Number, Value: "1"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "simplePairWithDecimal",
			input: "{\"a\":1.01}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: Number, Value: "1.01"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "simplePairWithDecimalAndSpaces",
			input: "{  \"a\"  :  1.01  }",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: Number, Value: "1.01"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "simplePairWithTrue",
			input: "{\"a\":true}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: Boolean, Value: "true"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "simplePairWithTrueAndSpaces",
			input: "{  \"a\"  :  true  }",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: Boolean, Value: "true"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "simplePairWithFalse",
			input: "{\"a\":false}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: Boolean, Value: "false"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "simplePairWithFalseAndSpaces",
			input: "{  \"a\"  :  false  }",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: Boolean, Value: "false"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "simplePairWithNull",
			input: "{\"a\":null}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: Null, Value: "null"},
				{tokenType: RightBrace, Value: "}"}}},
		{name: "simplePairWithNullAndSpaces",
			input: "{  \"a\"  :  null  }",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: Null, Value: "null"},
				{tokenType: RightBrace, Value: "}"}}},

		// Array
		{name: "simplePairWithArray",
			input: "{\"a\":[1, \"2\", true, null]}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: Number, Value: "1"},
				{tokenType: Comma, Value: ","},
				{tokenType: String, Value: "2"},
				{tokenType: Comma, Value: ","},
				{tokenType: Boolean, Value: "true"},
				{tokenType: Comma, Value: ","},
				{tokenType: Null, Value: "null"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBrace, Value: "}"}}},

		// Nested
		{name: "nested",
			input: "{\"a\": \"ünicode\", \"b\"  : { \"list\": [true, [false, null, {  }]]}}",
			expected: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: String, Value: "ünicode"},
				{tokenType: Comma, Value: ","},
				{tokenType: String, Value: "b"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "list"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: Boolean, Value: "true"},
				{tokenType: Comma, Value: ","},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: Boolean, Value: "false"},
				{tokenType: Comma, Value: ","},
				{tokenType: Null, Value: "null"},
				{tokenType: Comma, Value: ","},
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: RightBrace, Value: "}"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBrace, Value: "}"},
				{tokenType: RightBrace, Value: "}"},
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
			expected: "unexpected token . in JSON at position 9"},
		{name: "invalidText",
			input:    "{\"a\":    truef}",
			expected: "unexpected token truef in JSON at position 10"},
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
