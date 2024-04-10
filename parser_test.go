package jmatch

import (
	"reflect"
	"testing"
)

type parsingResult = map[string]Token

type collectorMatcher struct {
	collection parsingResult
}

func (c *collectorMatcher) Match(path string, token Token) {
	c.collection[path] = token // just collect everything
}

func newCollectorMatcher() collectorMatcher {
	return collectorMatcher{
		collection: make(parsingResult),
	}
}

func TestSuccessParse(t *testing.T) {
	testCases := []struct {
		name     string
		tokens   []Token
		expected parsingResult
	}{
		// single value
		{name: "'1'",
			tokens: []Token{
				{tokenType: String, Value: "1"},
			},
			expected: parsingResult{
				".": Token{tokenType: String, Value: "1"},
			},
		},

		// objects
		{name: "{'a': '1'}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: String, Value: "1"},
				{tokenType: RightBrace, Value: "}"},
			},
			expected: parsingResult{
				".a": Token{tokenType: String, Value: "1"},
			},
		},
		{name: "{'a': '1', 'b': '2', 'c': 3}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: String, Value: "1"},
				{tokenType: Comma, Value: ","},
				{tokenType: String, Value: "b"},
				{tokenType: Colon, Value: ":"},
				{tokenType: String, Value: "2"},
				{tokenType: Comma, Value: ","},
				{tokenType: String, Value: "c"},
				{tokenType: Colon, Value: ":"},
				{tokenType: String, Value: "3"},
				{tokenType: RightBrace, Value: "}"},
			},
			expected: parsingResult{
				".a": Token{tokenType: String, Value: "1"},
				".b": Token{tokenType: String, Value: "2"},
				".c": Token{tokenType: String, Value: "3"},
			},
		},
		{name: "{'a': { 'b': { 'c': 3}}}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "b"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "c"},
				{tokenType: Colon, Value: ":"},
				{tokenType: String, Value: "3"},
				{tokenType: RightBrace, Value: "}"},
				{tokenType: RightBrace, Value: "}"},
				{tokenType: RightBrace, Value: "}"},
			},
			expected: parsingResult{
				".a.b.c": Token{tokenType: String, Value: "3"},
			},
		},

		// arrays
		{name: "['1', '2']",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "["},
				{tokenType: String, Value: "1"},
				{tokenType: Comma, Value: ","},
				{tokenType: String, Value: "2"},
				{tokenType: RightBracket, Value: "]"},
			},
			expected: parsingResult{
				".[0]": Token{tokenType: String, Value: "1"},
				".[1]": Token{tokenType: String, Value: "2"},
			},
		},

		{name: "['1', ['2']]",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "["},
				{tokenType: String, Value: "1"},
				{tokenType: Comma, Value: ","},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: String, Value: "2"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBracket, Value: "]"},
			},
			expected: parsingResult{
				".[0]":     Token{tokenType: String, Value: "1"},
				".[1].[0]": Token{tokenType: String, Value: "2"},
			},
		},
		{name: "{'a': [[[[[[1]]]]]]}",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "["},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: String, Value: "1"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBracket, Value: "]"},
			},
			expected: parsingResult{
				".[0].[0].[0].[0].[0].[0]": Token{tokenType: String, Value: "1"},
			},
		},

		// mixed
		{name: "{'a': ['1', '2']}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: String, Value: "1"},
				{tokenType: Comma, Value: ","},
				{tokenType: String, Value: "2"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBrace, Value: "}"},
			},
			expected: parsingResult{
				".a.[0]": Token{tokenType: String, Value: "1"},
				".a.[1]": Token{tokenType: String, Value: "2"},
			},
		},

		{name: "{'a': [[[[[[1]]]]]]}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: String, Value: "1"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBrace, Value: "}"},
			},
			expected: parsingResult{
				".a.[0].[0].[0].[0].[0].[0]": Token{tokenType: String, Value: "1"},
			},
		},

		{name: "{'a': ['1', ['2', '3']]}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: String, Value: "1"},
				{tokenType: Comma, Value: ","},
				{tokenType: LeftBracket, Value: "]"},
				{tokenType: String, Value: "2"},
				{tokenType: Comma, Value: ","},
				{tokenType: String, Value: "3"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBrace, Value: "}"},
			},
			expected: parsingResult{
				".a.[0]":     Token{tokenType: String, Value: "1"},
				".a.[1].[0]": Token{tokenType: String, Value: "2"},
				".a.[1].[1]": Token{tokenType: String, Value: "3"},
			},
		},

		{name: "{'a': ['1', {'b': ['2', {'c': 3}]}]}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: String, Value: "1"},
				{tokenType: Comma, Value: ","},
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "b"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: String, Value: "2"},
				{tokenType: Comma, Value: ","},
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "c"},
				{tokenType: Colon, Value: ":"},
				{tokenType: String, Value: "3"},
				{tokenType: RightBrace, Value: "}"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBrace, Value: "}"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBrace, Value: "}"},
			},
			expected: parsingResult{
				".a.[0]":         Token{tokenType: String, Value: "1"},
				".a.[1].b.[0]":   Token{tokenType: String, Value: "2"},
				".a.[1].b.[1].c": Token{tokenType: String, Value: "3"},
			},
		},

		{name: "{'a': ['1', '2'], 'b': '3', 'c': ['4']}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: String, Value: "1"},
				{tokenType: Comma, Value: ","},
				{tokenType: String, Value: "2"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: String, Value: "b"},
				{tokenType: Colon, Value: ":"},
				{tokenType: String, Value: "3"},
				{tokenType: Comma, Value: ","},
				{tokenType: String, Value: "c"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: String, Value: "4"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBrace, Value: "}"},
			},
			expected: parsingResult{
				".a.[0]": Token{tokenType: String, Value: "1"},
				".a.[1]": Token{tokenType: String, Value: "2"},
				".b":     Token{tokenType: String, Value: "3"},
				".c.[0]": Token{tokenType: String, Value: "4"},
			},
		},

		{name: "{'s': {'t': [[1], -2.0, '3', true, {'x': false}]}}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "s"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "t"},
				{tokenType: Colon, Value: ":"},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: LeftBracket, Value: "["},
				{tokenType: Number, Value: "1"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: Comma, Value: ","},
				{tokenType: Number, Value: "-2.0"},
				{tokenType: Comma, Value: ","},
				{tokenType: String, Value: "3"},
				{tokenType: Comma, Value: ","},
				{tokenType: Boolean, Value: "true"},
				{tokenType: Comma, Value: ","},
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "x"},
				{tokenType: Colon, Value: ":"},
				{tokenType: Boolean, Value: "false"},
				{tokenType: RightBrace, Value: "}"},
				{tokenType: RightBracket, Value: "]"},
				{tokenType: RightBrace, Value: "}"},
				{tokenType: RightBrace, Value: "}"},
			},
			expected: parsingResult{
				".s.t.[0].[0]": Token{tokenType: Number, Value: "1"},
				".s.t.[1]":     Token{tokenType: Number, Value: "-2.0"},
				".s.t.[2]":     Token{tokenType: String, Value: "3"},
				".s.t.[3]":     Token{tokenType: Boolean, Value: "true"},
				".s.t.[4].x":   Token{tokenType: Boolean, Value: "false"},
			},
		},
	}

	for _, tc := range testCases {
		//_ = tc.expected
		t.Run(tc.name, func(t *testing.T) {
			m := newCollectorMatcher()
			p := newParser(tc.tokens, &m)
			err := p.parse()
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(m.collection, tc.expected) {
				t.Errorf("Expected '%v', got '%v' instead\n", tc.expected, m.collection)
			}
		})
	}

}

func TestFailParse(t *testing.T) {
	testCases := []struct {
		name     string
		tokens   []Token
		expected string
	}{
		{name: "{",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
			},
			expected: "invalid JSON. unexpected token { found at line 1 column 1"},
		{name: "{{",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 2},
			},
			expected: "invalid JSON. unexpected token { found at line 1 column 2"},
		{name: "{{{",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 2},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 3},
			},
			expected: "invalid JSON. unexpected token { found at line 1 column 2"},
		{name: "[",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "[", line: 1, column: 1},
			},
			expected: "invalid JSON. unexpected token [ found at line 1 column 1"},
		{name: "[[",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "[", line: 1, column: 1},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 2},
			},
			expected: "invalid JSON. unexpected token [ found at line 1 column 2"},
		{name: "[[[",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "[", line: 1, column: 1},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 2},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 3},
			},
			expected: "invalid JSON. unexpected token [ found at line 1 column 3"},
		{name: "{]",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: RightBracket, Value: "]", line: 1, column: 2},
			},
			expected: "invalid JSON. unexpected token ] found at line 1 column 2"},
		{name: "[}",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "[", line: 1, column: 1},
				{tokenType: RightBrace, Value: "}", line: 1, column: 2},
			},
			expected: "invalid JSON. unexpected token } found at line 1 column 2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := newCollectorMatcher()
			p := newParser(tc.tokens, &m)
			err := p.parse()

			if err == nil {
				t.Errorf("Expected error but input was parsed %v", m.collection)
			} else if err.Error() != tc.expected {
				t.Errorf("Expected '%s', got '%s' instead\n", tc.expected, err)
			}
		})
	}
}
