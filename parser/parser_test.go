package parser

import (
	"reflect"
	"testing"

	token "github.com/rodic/jmatch/token"
)

type parsingResult = map[string]token.Token

type collectorMatcher struct {
	collection parsingResult
}

func (c *collectorMatcher) Match(path string, token token.Token) {
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
		tokens   []token.Token
		expected parsingResult
	}{
		// single value
		{name: "'1'",
			tokens: []token.Token{
				token.NewStringToken("1", 1, 1),
			},
			expected: parsingResult{
				".": token.NewStringToken("1", 1, 1),
			},
		},

		// objects
		{name: "{}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewRightBraceToken(1, 1),
			},
			expected: parsingResult{},
		},
		{name: "{'a': '1'}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewStringToken("1", 1, 4),
				token.NewRightBraceToken(1, 5),
			},
			expected: parsingResult{
				".a": token.NewStringToken("1", 1, 4),
			},
		},
		{name: "{'a.b': '1'}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a.b", 1, 2),
				token.NewColonToken(1, 3),
				token.NewStringToken("1", 1, 4),
				token.NewRightBraceToken(1, 5),
			},
			expected: parsingResult{
				".\"a.b\"": token.NewStringToken("1", 1, 4),
			},
		},
		{name: "{'a b': '1'}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a b", 1, 2),
				token.NewColonToken(1, 3),
				token.NewStringToken("1", 1, 4),
				token.NewRightBraceToken(1, 5),
			},
			expected: parsingResult{
				".\"a b\"": token.NewStringToken("1", 1, 4),
			},
		},
		{name: "{'a': '1', 'b': '2', 'c': 3}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewStringToken("1", 1, 4),
				token.NewCommaToken(1, 5),
				token.NewStringToken("b", 1, 6),
				token.NewColonToken(1, 7),
				token.NewStringToken("2", 1, 8),
				token.NewCommaToken(1, 9),
				token.NewStringToken("c", 1, 10),
				token.NewColonToken(1, 111),
				token.NewStringToken("3", 1, 12),
				token.NewRightBraceToken(1, 5),
			},
			expected: parsingResult{
				".a": token.NewStringToken("1", 1, 4),
				".b": token.NewStringToken("2", 1, 8),
				".c": token.NewStringToken("3", 1, 12),
			},
		},
		{name: "{'a': { 'b': { 'c': 3}}}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBraceToken(1, 4),
				token.NewStringToken("b", 1, 5),
				token.NewColonToken(1, 6),
				token.NewLeftBraceToken(1, 7),
				token.NewStringToken("c", 1, 8),
				token.NewColonToken(1, 9),
				token.NewStringToken("3", 1, 10),
				token.NewRightBraceToken(1, 11),
				token.NewRightBraceToken(1, 12),
				token.NewRightBraceToken(1, 13),
			},
			expected: parsingResult{
				".a.b.c": token.NewStringToken("3", 1, 10),
			},
		},

		// arrays
		{name: "['1', '2']",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewStringToken("1", 1, 2),
				token.NewCommaToken(1, 3),
				token.NewStringToken("2", 1, 4),
				token.NewRightBracketToken(1, 5),
			},
			expected: parsingResult{
				".[0]": token.NewStringToken("1", 1, 2),
				".[1]": token.NewStringToken("2", 1, 4),
			},
		},
		{name: "['1', ['2']]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewStringToken("1", 1, 2),
				token.NewCommaToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewStringToken("2", 1, 5),
				token.NewRightBracketToken(1, 6),
				token.NewRightBracketToken(1, 7),
			},
			expected: parsingResult{
				".[0]":    token.NewStringToken("1", 1, 2),
				".[1][0]": token.NewStringToken("2", 1, 5),
			},
		},
		{name: "[{}]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewLeftBraceToken(1, 3),
				token.NewRightBraceToken(1, 5),
				token.NewRightBracketToken(1, 1),
			},
			expected: parsingResult{},
		},
		{name: "[{'a': 1}]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewLeftBraceToken(1, 2),
				token.NewStringToken("a", 1, 3),
				token.NewColonToken(1, 4),
				token.NewNumberToken("1", 1, 5),
				token.NewRightBraceToken(1, 6),
				token.NewRightBracketToken(1, 7),
			},
			expected: parsingResult{
				".[0].a": token.NewNumberToken("1", 1, 5),
			},
		},
		{name: "[[[[[[1]]]]]]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewLeftBracketToken(1, 2),
				token.NewLeftBracketToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewLeftBracketToken(1, 5),
				token.NewLeftBracketToken(1, 6),
				token.NewStringToken("1", 1, 7),
				token.NewRightBracketToken(1, 8),
				token.NewRightBracketToken(1, 9),
				token.NewRightBracketToken(1, 10),
				token.NewRightBracketToken(1, 11),
				token.NewRightBracketToken(1, 12),
				token.NewRightBracketToken(1, 13),
			},
			expected: parsingResult{
				".[0][0][0][0][0][0]": token.NewStringToken("1", 1, 7),
			},
		},

		// mixed
		{name: "{'a': ['1', '2']}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewStringToken("1", 1, 5),
				token.NewCommaToken(1, 6),
				token.NewStringToken("2", 1, 7),
				token.NewRightBracketToken(1, 1),
				token.NewRightBraceToken(1, 5),
			},
			expected: parsingResult{
				".a[0]": token.NewStringToken("1", 1, 5),
				".a[1]": token.NewStringToken("2", 1, 7),
			},
		},
		{name: "{'a': [[[[[[1]]]]]]}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewLeftBracketToken(1, 5),
				token.NewLeftBracketToken(1, 6),
				token.NewLeftBracketToken(1, 7),
				token.NewLeftBracketToken(1, 8),
				token.NewLeftBracketToken(1, 9),
				token.NewStringToken("1", 1, 10),
				token.NewRightBracketToken(1, 11),
				token.NewRightBracketToken(1, 12),
				token.NewRightBracketToken(1, 13),
				token.NewRightBracketToken(1, 14),
				token.NewRightBracketToken(1, 15),
				token.NewRightBracketToken(1, 16),
				token.NewRightBraceToken(1, 17),
			},
			expected: parsingResult{
				".a[0][0][0][0][0][0]": token.NewStringToken("1", 1, 10),
			},
		},
		{name: "{'a': ['1', ['2', '3']]}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewStringToken("1", 1, 5),
				token.NewCommaToken(1, 6),
				token.NewLeftBracketToken(1, 7),
				token.NewStringToken("2", 1, 8),
				token.NewCommaToken(1, 9),
				token.NewStringToken("3", 1, 10),
				token.NewRightBracketToken(1, 11),
				token.NewRightBracketToken(1, 12),
				token.NewRightBraceToken(1, 13),
			},
			expected: parsingResult{
				".a[0]":    token.NewStringToken("1", 1, 5),
				".a[1][0]": token.NewStringToken("2", 1, 8),
				".a[1][1]": token.NewStringToken("3", 1, 10),
			},
		},
		{name: "{'a': ['1', {'b': ['2', {'c': 3}]}]}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewStringToken("1", 1, 5),
				token.NewCommaToken(1, 6),
				token.NewLeftBraceToken(1, 7),
				token.NewStringToken("b", 1, 8),
				token.NewColonToken(1, 9),
				token.NewLeftBracketToken(1, 10),
				token.NewStringToken("2", 1, 11),
				token.NewCommaToken(1, 12),
				token.NewLeftBraceToken(1, 13),
				token.NewStringToken("c", 1, 14),
				token.NewColonToken(1, 15),
				token.NewStringToken("3", 1, 16),
				token.NewRightBraceToken(1, 17),
				token.NewRightBracketToken(1, 18),
				token.NewRightBraceToken(1, 19),
				token.NewRightBracketToken(1, 20),
				token.NewRightBraceToken(1, 21),
			},
			expected: parsingResult{
				".a[0]":        token.NewStringToken("1", 1, 5),
				".a[1].b[0]":   token.NewStringToken("2", 1, 11),
				".a[1].b[1].c": token.NewStringToken("3", 1, 16),
			},
		},
		{name: "{'a': ['1', '2'], 'b': '3', 'c': ['4']}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewStringToken("1", 1, 5),
				token.NewCommaToken(1, 6),
				token.NewStringToken("2", 1, 7),
				token.NewRightBracketToken(1, 8),
				token.NewCommaToken(1, 9),
				token.NewStringToken("b", 1, 10),
				token.NewColonToken(1, 11),
				token.NewStringToken("3", 1, 12),
				token.NewCommaToken(1, 13),
				token.NewStringToken("c", 1, 14),
				token.NewColonToken(1, 15),
				token.NewLeftBracketToken(1, 16),
				token.NewStringToken("4", 1, 17),
				token.NewRightBracketToken(1, 18),
				token.NewRightBraceToken(1, 19),
			},
			expected: parsingResult{
				".a[0]": token.NewStringToken("1", 1, 5),
				".a[1]": token.NewStringToken("2", 1, 7),
				".b":    token.NewStringToken("3", 1, 12),
				".c[0]": token.NewStringToken("4", 1, 17),
			},
		},
		{name: "{'s': {'t': [[1], -2.0, '3', true, {'x': false}]}}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("s", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBraceToken(1, 4),
				token.NewStringToken("t", 1, 5),
				token.NewColonToken(1, 6),
				token.NewLeftBracketToken(1, 7),
				token.NewLeftBracketToken(1, 8),
				token.NewNumberToken("1", 1, 9),
				token.NewRightBracketToken(1, 10),
				token.NewCommaToken(1, 11),
				token.NewNumberToken("-2.0", 1, 12),
				token.NewCommaToken(1, 13),
				token.NewStringToken("3", 1, 14),
				token.NewCommaToken(1, 15),
				token.NewBooleanToken("true", 1, 16),
				token.NewCommaToken(1, 17),
				token.NewLeftBraceToken(1, 18),
				token.NewStringToken("x", 1, 19),
				token.NewColonToken(1, 20),
				token.NewBooleanToken("false", 1, 21),
				token.NewRightBraceToken(1, 22),
				token.NewRightBracketToken(1, 23),
				token.NewRightBraceToken(1, 24),
				token.NewRightBraceToken(1, 25),
			},
			expected: parsingResult{
				".s.t[0][0]": token.NewNumberToken("1", 1, 9),
				".s.t[1]":    token.NewNumberToken("-2.0", 1, 12),
				".s.t[2]":    token.NewStringToken("3", 1, 14),
				".s.t[3]":    token.NewBooleanToken("true", 1, 16),
				".s.t[4].x":  token.NewBooleanToken("false", 1, 21),
			},
		},
	}

	for _, tc := range testCases {
		//_ = tc.expected
		t.Run(tc.name, func(t *testing.T) {
			m := newCollectorMatcher()
			p := NewParser(tc.tokens, &m)
			err := p.Parse()
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
		tokens   []token.Token
		expected string
	}{
		// objects
		{name: "{",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{{",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewLeftBraceToken(1, 2),
			},
			expected: "invalid JSON. unexpected token { at line 1 column 2"},
		{name: "{{{",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewLeftBraceToken(1, 2),
				token.NewLeftBraceToken(1, 3),
			},
			expected: "invalid JSON. unexpected token { at line 1 column 2"},
		{name: "{{}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewLeftBraceToken(1, 2),
				token.NewRightBraceToken(1, 3),
			},
			expected: "invalid JSON. unexpected token { at line 1 column 2"},
		{name: "{}}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewRightBraceToken(1, 2),
				token.NewRightBraceToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{1",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewNumberToken("1", 1, 2),
			},
			expected: "invalid JSON. unexpected token 1 at line 1 column 2"},
		{name: "{true",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewBooleanToken("true", 1, 2),
			},
			expected: "invalid JSON. unexpected token true at line 1 column 2"},
		{name: "{null",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewNullToken(1, 2),
			},
			expected: "invalid JSON. unexpected token null at line 1 column 2"},
		{name: "{:",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewColonToken(1, 2),
			},
			expected: "invalid JSON. unexpected token : at line 1 column 2"},
		{name: "{,",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewCommaToken(1, 2),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 2"},
		{name: "{'a',",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewCommaToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{'a',1}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewCommaToken(1, 3),
				token.NewNumberToken("1", 1, 4),
				token.NewRightBraceToken(1, 5),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 3"},
		{name: "{'a','1'}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewCommaToken(1, 3),
				token.NewStringToken("1", 1, 4),
				token.NewRightBraceToken(1, 5),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 3"},
		{name: "{'a': 1,}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewNumberToken("1", 1, 4),
				token.NewCommaToken(1, 5),
				token.NewRightBraceToken(1, 6),
			},
			expected: "invalid JSON. unexpected token } at line 1 column 6"},
		{name: "{'a': 1, 2}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewNumberToken("1", 1, 4),
				token.NewCommaToken(1, 5),
				token.NewNumberToken("2", 1, 6),
				token.NewRightBraceToken(1, 7),
			},
			expected: "invalid JSON. unexpected token 2 at line 1 column 6"},
		{name: "{'a': 'b': 1}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewStringToken("b", 1, 4),
				token.NewColonToken(1, 5),
				token.NewNumberToken("1", 1, 6),
				token.NewRightBraceToken(1, 7),
			},
			expected: "invalid JSON. unexpected token : at line 1 column 5"},
		{name: "{'a': {,}}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBraceToken(1, 4),
				token.NewCommaToken(1, 5),
				token.NewRightBraceToken(1, 6),
				token.NewRightBraceToken(1, 7),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 5"},
		{name: "{'a': {{}}}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBraceToken(1, 4),
				token.NewLeftBraceToken(1, 5),
				token.NewRightBraceToken(1, 6),
				token.NewRightBraceToken(1, 7),
				token.NewRightBraceToken(1, 8),
			},
			expected: "invalid JSON. unexpected token { at line 1 column 5"},

		// arrays
		{name: "[",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[[",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewLeftBracketToken(1, 2),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[[[",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewLeftBracketToken(1, 2),
				token.NewLeftBracketToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[[]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewLeftBracketToken(1, 2),
				token.NewRightBracketToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[]]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewRightBracketToken(1, 2),
				token.NewRightBracketToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[,",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewCommaToken(1, 2),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 2"},
		{name: "[1,]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewNumberToken("1", 1, 2),
				token.NewCommaToken(1, 3),
				token.NewRightBracketToken(1, 4),
			},
			expected: "invalid JSON. unexpected token ] at line 1 column 4"},

		// mixed
		{name: "{]",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewRightBracketToken(1, 2),
			},
			expected: "invalid JSON. unexpected token ] at line 1 column 2"},
		{name: "[}",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewRightBraceToken(1, 2),
			},
			expected: "invalid JSON. unexpected token } at line 1 column 2"},
		{name: "{[}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewLeftBracketToken(1, 2),
				token.NewRightBraceToken(1, 3),
			},
			expected: "invalid JSON. unexpected token [ at line 1 column 2"},
		{name: "{'a': [}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewRightBraceToken(1, 5),
			},
			expected: "invalid JSON. unexpected token } at line 1 column 5"},
		{name: "{'a': [,",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewCommaToken(1, 5),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 5"},
		{name: "{'a': [,",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewCommaToken(1, 5),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 5"},
		{name: "{'a': [{},",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewLeftBraceToken(1, 5),
				token.NewRightBraceToken(1, 6),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{'s': {'t': [[1], -2.0, '3', true, {'x': false,}]}}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("s", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBraceToken(1, 4),
				token.NewStringToken("t", 1, 5),
				token.NewColonToken(1, 6),
				token.NewLeftBracketToken(1, 7),
				token.NewLeftBracketToken(1, 8),
				token.NewNumberToken("1", 1, 9),
				token.NewRightBracketToken(1, 10),
				token.NewCommaToken(1, 11),
				token.NewNumberToken("-2.0", 1, 12),
				token.NewCommaToken(1, 13),
				token.NewStringToken("3", 1, 14),
				token.NewCommaToken(1, 15),
				token.NewBooleanToken("true", 1, 16),
				token.NewCommaToken(1, 17),
				token.NewLeftBraceToken(1, 18),
				token.NewStringToken("x", 1, 19),
				token.NewColonToken(1, 20),
				token.NewBooleanToken("false", 1, 21),
				token.NewCommaToken(1, 22),
				token.NewRightBraceToken(1, 23),
				token.NewRightBracketToken(1, 24),
				token.NewRightBraceToken(1, 25),
				token.NewRightBraceToken(1, 26),
			},
			expected: "invalid JSON. unexpected token } at line 1 column 23",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := newCollectorMatcher()
			p := NewParser(tc.tokens, &m)

			err := p.Parse()

			if err == nil {
				t.Errorf("Expected error but input was parsed %v", m.collection)
			} else if err.Error() != tc.expected {
				t.Errorf("Expected '%s', got '%s' instead\n", tc.expected, err)
			}
		})
	}
}
