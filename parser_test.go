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
		{name: "{}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: RightBrace, Value: "}"},
			},
			expected: parsingResult{},
		},
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
		{name: "[{}]",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "["},
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: RightBrace, Value: "}"},
				{tokenType: RightBracket, Value: "]"},
			},
			expected: parsingResult{},
		},
		{name: "[{'a': 1}]",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "["},
				{tokenType: LeftBrace, Value: "{"},
				{tokenType: String, Value: "a"},
				{tokenType: Colon, Value: ":"},
				{tokenType: Number, Value: "1"},
				{tokenType: RightBrace, Value: "}"},
				{tokenType: RightBracket, Value: "]"},
			},
			expected: parsingResult{
				".[0].a": Token{tokenType: Number, Value: "1"},
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
				{tokenType: Comma, Value: ","},
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
		// objects
		{name: "{",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{{",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 2},
			},
			expected: "invalid JSON. unexpected token { at line 1 column 2"},
		{name: "{{{",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 2},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 3},
			},
			expected: "invalid JSON. unexpected token { at line 1 column 2"},
		{name: "{{}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 2},
				{tokenType: RightBrace, Value: "}", line: 1, column: 3},
			},
			expected: "invalid JSON. unexpected token { at line 1 column 2"},
		{name: "{}}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: RightBrace, Value: "}", line: 1, column: 2},
				{tokenType: RightBrace, Value: "}", line: 1, column: 3},
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{1",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: Number, Value: "1", line: 1, column: 2},
			},
			expected: "invalid JSON. unexpected token 1 at line 1 column 2"},
		{name: "{true",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: Boolean, Value: "true", line: 1, column: 2},
			},
			expected: "invalid JSON. unexpected token true at line 1 column 2"},
		{name: "{null",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: Null, Value: "null", line: 1, column: 2},
			},
			expected: "invalid JSON. unexpected token null at line 1 column 2"},
		{name: "{:",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: Colon, Value: ":", line: 1, column: 2},
			},
			expected: "invalid JSON. unexpected token : at line 1 column 2"},
		{name: "{,",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: Comma, Value: ",", line: 1, column: 2},
			},
			expected: "invalid JSON. unexpected token , at line 1 column 2"},
		{name: "{'a',",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Comma, Value: ",", line: 1, column: 3},
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{'a',1}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Comma, Value: ",", line: 1, column: 3},
				{tokenType: Number, Value: "1", line: 1, column: 4},
				{tokenType: RightBrace, Value: "}", line: 1, column: 5},
			},
			expected: "invalid JSON. unexpected token , at line 1 column 3"},
		{name: "{'a','1'}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Comma, Value: ",", line: 1, column: 3},
				{tokenType: String, Value: "1", line: 1, column: 4},
				{tokenType: RightBrace, Value: "}", line: 1, column: 5},
			},
			expected: "invalid JSON. unexpected token , at line 1 column 3"},
		{name: "{'a': 1,}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 3},
				{tokenType: Number, Value: "1", line: 1, column: 4},
				{tokenType: Comma, Value: ",", line: 1, column: 5},
				{tokenType: RightBrace, Value: "}", line: 1, column: 6},
			},
			expected: "invalid JSON. unexpected token } at line 1 column 6"},
		{name: "{'a': 1, 2}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 3},
				{tokenType: Number, Value: "1", line: 1, column: 4},
				{tokenType: Comma, Value: ",", line: 1, column: 5},
				{tokenType: Number, Value: "2", line: 1, column: 6},
				{tokenType: RightBrace, Value: "}", line: 1, column: 7},
			},
			expected: "invalid JSON. unexpected token 2 at line 1 column 6"},
		{name: "{'a': 'b': 1}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 3},
				{tokenType: String, Value: "b", line: 1, column: 4},
				{tokenType: Colon, Value: ":", line: 1, column: 5},
				{tokenType: Number, Value: "1", line: 1, column: 6},
				{tokenType: RightBrace, Value: "}", line: 1, column: 7},
			},
			expected: "invalid JSON. unexpected token : at line 1 column 5"},
		{name: "{'a': {,}}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 3},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 4},
				{tokenType: Comma, Value: ",", line: 1, column: 5},
				{tokenType: RightBrace, Value: "}", line: 1, column: 6},
				{tokenType: RightBrace, Value: "}", line: 1, column: 7},
			},
			expected: "invalid JSON. unexpected token , at line 1 column 5"},
		{name: "{'a': {{}}}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 3},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 4},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 5},
				{tokenType: RightBrace, Value: "}", line: 1, column: 6},
				{tokenType: RightBrace, Value: "}", line: 1, column: 7},
				{tokenType: RightBrace, Value: "}", line: 1, column: 8},
			},
			expected: "invalid JSON. unexpected token { at line 1 column 5"},

		// arrays
		{name: "[",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "[", line: 1, column: 1},
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[[",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "[", line: 1, column: 1},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 2},
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[[[",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "[", line: 1, column: 1},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 2},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 3},
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[[]",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "[", line: 1, column: 1},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 2},
				{tokenType: RightBracket, Value: "]", line: 1, column: 3},
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[]]",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "[", line: 1, column: 1},
				{tokenType: RightBracket, Value: "]", line: 1, column: 2},
				{tokenType: RightBracket, Value: "]", line: 1, column: 3},
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[,",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "[", line: 1, column: 1},
				{tokenType: Comma, Value: ",", line: 1, column: 2},
			},
			expected: "invalid JSON. unexpected token , at line 1 column 2"},
		{name: "[1,]",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "[", line: 1, column: 1},
				{tokenType: Number, Value: "1", line: 1, column: 2},
				{tokenType: Comma, Value: ",", line: 1, column: 3},
				{tokenType: RightBracket, Value: "]", line: 1, column: 4},
			},
			expected: "invalid JSON. unexpected token ] at line 1 column 4"},

		// mixed
		{name: "{]",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: RightBracket, Value: "]", line: 1, column: 2},
			},
			expected: "invalid JSON. unexpected token ] at line 1 column 2"},
		{name: "[}",
			tokens: []Token{
				{tokenType: LeftBracket, Value: "[", line: 1, column: 1},
				{tokenType: RightBrace, Value: "}", line: 1, column: 2},
			},
			expected: "invalid JSON. unexpected token } at line 1 column 2"},
		{name: "{[}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 2},
				{tokenType: RightBrace, Value: "}", line: 1, column: 3},
			},
			expected: "invalid JSON. unexpected token [ at line 1 column 2"},
		{name: "{'a': [}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 3},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 4},
				{tokenType: RightBrace, Value: "}", line: 1, column: 5},
			},
			expected: "invalid JSON. unexpected token } at line 1 column 5"},
		{name: "{'a': [,",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 3},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 4},
				{tokenType: Comma, Value: ",", line: 1, column: 5},
			},
			expected: "invalid JSON. unexpected token , at line 1 column 5"},
		{name: "{'a': [,",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 3},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 4},
				{tokenType: Comma, Value: ",", line: 1, column: 5},
			},
			expected: "invalid JSON. unexpected token , at line 1 column 5"},
		{name: "{'a': [{},",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "a", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 3},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 4},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 5},
				{tokenType: RightBrace, Value: "}", line: 1, column: 6},
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{'s': {'t': [[1], -2.0, '3', true, {'x': false,}]}}",
			tokens: []Token{
				{tokenType: LeftBrace, Value: "{", line: 1, column: 1},
				{tokenType: String, Value: "s", line: 1, column: 2},
				{tokenType: Colon, Value: ":", line: 1, column: 3},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 4},
				{tokenType: String, Value: "t", line: 1, column: 5},
				{tokenType: Colon, Value: ":", line: 1, column: 6},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 7},
				{tokenType: LeftBracket, Value: "[", line: 1, column: 8},
				{tokenType: Number, Value: "1", line: 1, column: 9},
				{tokenType: RightBracket, Value: "]", line: 1, column: 10},
				{tokenType: Comma, Value: ",", line: 1, column: 11},
				{tokenType: Number, Value: "-2.0", line: 1, column: 12},
				{tokenType: Comma, Value: ",", line: 1, column: 13},
				{tokenType: String, Value: "3", line: 1, column: 14},
				{tokenType: Comma, Value: ",", line: 1, column: 15},
				{tokenType: Boolean, Value: "true", line: 1, column: 16},
				{tokenType: Comma, Value: ",", line: 1, column: 17},
				{tokenType: LeftBrace, Value: "{", line: 1, column: 18},
				{tokenType: String, Value: "x", line: 1, column: 19},
				{tokenType: Colon, Value: ":", line: 1, column: 20},
				{tokenType: Boolean, Value: "false", line: 1, column: 21},
				{tokenType: Comma, Value: ",", line: 1, column: 22},
				{tokenType: RightBrace, Value: "}", line: 1, column: 23},
				{tokenType: RightBracket, Value: "]", line: 1, column: 24},
				{tokenType: RightBrace, Value: "}", line: 1, column: 25},
				{tokenType: RightBrace, Value: "}", line: 1, column: 26},
			},
			expected: "invalid JSON. unexpected token } at line 1 column 23",
		},
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
